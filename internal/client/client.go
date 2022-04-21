package client

import (
	"context"
	"fmt"
	"io"
	"net"
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
		message, err := SinglePoWAlgoExecution(*cli.ctx, *cli.conn, *cli.conn)
		if err != nil {
			return err
		}
		fmt.Println("quote result:", message)
		time.Sleep(5 * time.Second)
	}
}

func SinglePoWAlgoExecution(ctx context.Context, readerConn io.Reader, writerConn io.Writer) (string, error) {
	return "", nil
}
