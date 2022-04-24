package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"pow-client-server/internal/pkg/processor"
)

type server struct {
	listener *net.Listener
	ctx      *context.Context
}

func NewServer(serverUrl string, ctx context.Context) *server {
	listener, err := net.Listen("tcp", serverUrl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Established connection to server{%s}\n", serverUrl)
	return &server{
		listener: &listener,
		ctx:      &ctx,
	}
}

func (srv *server) Listen() error {
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}(*srv.listener)

	fmt.Printf("Started listening to address{%s}\n", (*srv.listener).Addr())
	for {
		// Listen for an incoming connection.
		conn, err := (*srv.listener).Accept()
		if err != nil {
			return fmt.Errorf("connection accept error: %w", err)
		}
		// Handle connections in a new goroutine.
		go handleClientConn(*srv.ctx, conn)
	}
}

func handleClientConn(ctx context.Context, conn net.Conn) {
	fmt.Println("Established new client:", conn.RemoteAddr())
	defer conn.Close()

	connectionReader := bufio.NewReader(conn)
	for {
		clientRequest, err := connectionReader.ReadString('\n')
		if err != nil {
			fmt.Printf("handled error while reading connection buffer - %s\n", err)
			return
		}
		msg, err := processor.Process(clientRequest, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			msgStr := fmt.Sprintf("%s\n", msg.Serialize())
			_, err := conn.Write([]byte(msgStr))
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}

}
