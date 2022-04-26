package main

import (
	"fmt"
	"os"
	"pow-client-server/internal/client"
)

func main() {
	fmt.Println("Starting tcp client")

	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")

	serverUrl := serverHost + ":" + serverPort

	cli := client.NewClient(serverUrl)
	err := cli.Execute()
	if err != nil {
		panic(err)
	}
}
