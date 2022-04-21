package server

import (
	"context"
	"fmt"
	"net"
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
		go handleConn(*srv.ctx, conn)
	}
}

func handleConn(ctx context.Context, conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()
}
