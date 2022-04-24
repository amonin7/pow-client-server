package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"pow-client-server/internal/pkg/model"
	"pow-client-server/internal/pkg/pow/isrm"
	"time"
)

type client struct {
	conn *net.Conn
	ctx  *context.Context
}

func NewClient(serverUrl string, ctx context.Context) *client {
	conn, err := net.Dial("tcp", serverUrl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Established connection to server{%s}\n", serverUrl)
	return &client{
		conn: &conn,
		ctx:  &ctx,
	}
}

func (cli *client) Execute() error {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(*cli.conn)

	for {
		message, err := SinglePoWAlgoExecution(*cli.conn, *cli.conn)
		if err != nil {
			return err
		}
		fmt.Println("quote result:", message)
		time.Sleep(5 * time.Second)
	}
}

func SinglePoWAlgoExecution(readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)

	message := &model.Message{
		Header: model.ChallengeReq,
	}
	msgStr := fmt.Sprintf("%s\n", message.Serialize())
	_, err := writerConn.Write([]byte(msgStr))
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	// reading and parsing response
	msgStr, err = reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	message, err = model.DeserializeMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	var integerSquareRootModulo isrm.IntegerSquareRootModulo
	err = json.Unmarshal([]byte(message.Payload), &integerSquareRootModulo)
	if err != nil {
		return "", fmt.Errorf("err parse hashcash: %w", err)
	}
	fmt.Println("got hashcash:", integerSquareRootModulo.Serialize())

	// 2. got challenge, compute hashcash
	err = integerSquareRootModulo.FindSolution()
	println(integerSquareRootModulo.Proof)
	if err != nil {
		return "", fmt.Errorf("err compute hashcash: %w", err)
	}
	fmt.Println("proof found: ", integerSquareRootModulo.Proof)
	// marshal solution to json
	byteData, err := json.Marshal(integerSquareRootModulo)
	if err != nil {
		return "", fmt.Errorf("err marshal hashcash: %w", err)
	}

	// 3. send challenge solution back to server
	message = &model.Message{
		Header:  model.ResourceReq,
		Payload: string(byteData),
	}
	msgStr = fmt.Sprintf("%s\n", message.Serialize())
	_, err = writerConn.Write([]byte(msgStr))
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}
	fmt.Println("challenge sent to server")

	// 4. get result quote from server
	msgStr, err = reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	message, err = model.DeserializeMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	return message.Payload, nil
}
