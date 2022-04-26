package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"pow-client-server/internal/pkg/model"
	"pow-client-server/internal/pkg/pow/isrm"
	"time"
)

// client - structure that represents client
type client struct {
	conn *net.Conn
}

// NewClient - creates new client from the provided server url
func NewClient(serverUrl string) *client {
	conn, err := net.Dial("tcp", serverUrl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Established connection to server{%s}\n", serverUrl)
	return &client{
		conn: &conn,
	}
}

// Execute - start executing requests to server
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
		fmt.Printf("Word of wisdom string received from server: %s\n", message)
		// some time threshold to pause the algorithm
		time.Sleep(10 * time.Second)
	}
}

// SinglePoWAlgoExecution - function to execute one cycle of Proof of Work algorithm from the client side
// Algorithm:
//	1. send request for challenge
//	2. receive the challenge
//	3. calculate the proof for the challenge
//	4. send it back to server with the request for resource
//	5. receive the resource from server
func SinglePoWAlgoExecution(rConnection io.Reader, wConnection io.Writer) (string, error) {
	reader := bufio.NewReader(rConnection)

	// step 1 send request for challenge
	message := &model.Message{
		Header: model.ChallengeReq,
	}
	serialized, err := message.Serialize()
	if err != nil {
		return "", err
	}
	_, err = wConnection.Write(append(serialized, byte('\n')))
	if err != nil {
		return "", fmt.Errorf("failed to send request for challenge to the server: %w", err)
	}

	// step 2 receive the challenge
	serializedMessage, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read message from connection: %w", err)
	}
	message, err = model.DeserializeMessage(serializedMessage)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize message: %w", err)
	}
	integerSquareRootModulo, err := isrm.DeserializeIsrm(message.Payload)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize challenge: %w", err)
	}
	fmt.Printf("received challenge from the server: %s\n", integerSquareRootModulo.ToShortString())

	// step 3 calculate the proof for the challenge
	err = integerSquareRootModulo.FindSolution()
	if err != nil {
		return "", fmt.Errorf("err compute hashcash: %w", err)
	}
	fmt.Println("proof found: ", integerSquareRootModulo.Proof)
	byteData, err := json.Marshal(integerSquareRootModulo)
	if err != nil {
		return "", fmt.Errorf("err marshal hashcash: %w", err)
	}

	// step 4 send challenge back to server with the request for resource
	message = &model.Message{
		Header:  model.ResourceReq,
		Payload: string(byteData),
	}
	serialized, err = message.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize message, which should be sent to the server: %w", err)
	}
	_, err = wConnection.Write(append(serialized, byte('\n')))
	if err != nil {
		return "", fmt.Errorf("failed to send the request for resource: %w", err)
	}
	fmt.Println("challenge was successfully sent to server")

	//	step 5 receive the resource from server
	serializedMessage, err = reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	message, err = model.DeserializeMessage(serializedMessage)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	return message.Payload, nil
}
