package server

import (
	"bufio"
	"fmt"
	"net"
	"pow-client-server/internal/pkg/processor"
)

// server - structure that represents the server side
type server struct {
	listener *net.Listener
}

// NewServer - generates new instance of server by the provided url
func NewServer(serverUrl string) *server {
	listener, err := net.Listen("tcp", serverUrl)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully created the server - message listener on url: %s\n", serverUrl)
	return &server{
		listener: &listener,
	}
}

// Listen - starts listening to clients and processing their messages
func (srv *server) Listen() error {
	// defer the connection close
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}(*srv.listener)

	fmt.Printf("Started listening to address{%s}\n", (*srv.listener).Addr())
	for {
		conn, err := (*srv.listener).Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}
		go handleClientConn(conn)
	}
}

// handleClientConn - processes each connection separately
// Algorithm of handling client connection
// 1. read the incoming message from the client
// 2. try to process it with special request processor
// 3. send the response for the client request (if processing was successful)
func handleClientConn(conn net.Conn) {
	fmt.Printf("Established new client: %s\n", conn.RemoteAddr())
	connectionReader := bufio.NewReader(conn)
	for {
		// step 1
		clientRequest, err := connectionReader.ReadString('\n')
		if err != nil {
			fmt.Printf("failed to read the request from the connection - %s\n", err)
			return
		}
		// step 2
		message, err := processor.Process(clientRequest, conn.RemoteAddr().String())
		if err != nil || message == nil {
			fmt.Printf("failed to process incoming message: %s\n", err)
			return
		}
		// step 3
		serialized, err := message.Serialize()
		if err != nil {
			fmt.Printf("failed to serialize message: %s\n", err)
			return
		}
		_, err = conn.Write(append(serialized, byte('\n')))
		if err != nil {
			fmt.Printf("failed to send message back to client: %s\n", err)
		}
	}

}
