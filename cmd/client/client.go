package main

import (
	"context"
	"fmt"
	"os"
	"pow-client-server/internal/client"
)

func main() {
	fmt.Println("Starting tpc client")

	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	ctx := context.Background()

	serverUrl := serverHost + ":" + serverPort

	cli := client.NewClient(serverUrl, ctx)
	err := cli.Execute()
	if err != nil {
		panic(err)
	}
}
