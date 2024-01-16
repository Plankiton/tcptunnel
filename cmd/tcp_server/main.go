package main

import (
	"os"

	"github.com/plankiton/tcptunnel/pkg/services/server"
)

func main() {
	serverPort := os.Getenv("TCPTUNNEL_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	server.StartServer(serverPort)
}
