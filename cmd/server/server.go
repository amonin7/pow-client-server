package main

import (
	"context"
	"fmt"
	"os"
	"pow-client-server/internal/server"
)

func main() {
	fmt.Println("Starting tpc client")

	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	ctx := context.Background()

	serverUrl := serverHost + ":" + serverPort

	srv := server.NewServer(serverUrl, ctx)
	err := srv.Listen()
	if err != nil {
		panic(err)
	}
}
