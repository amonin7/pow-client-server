package main

import (
	"fmt"
	"os"
	"pow-client-server/internal/server"
)

func main() {
	fmt.Println("Starting tcp server")

	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")

	serverUrl := serverHost + ":" + serverPort

	srv := server.NewServer(serverUrl)
	err := srv.Listen()
	if err != nil {
		panic(err)
	}
}
