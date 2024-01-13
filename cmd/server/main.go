package main

import (
	"fmt"
	"net"
	"os"

	"github.com/plankiton/tcptunnel/pkg/models"
)

func main() {
	serverPort := os.Getenv("TCPTUNNEL_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	listener, err := net.Listen("tcp", fmt.Sprint(":", serverPort))
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Servidor escutando na porta %s...\n", serverPort)

	go streamMessages()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		clientID, _ := generateRandomHexCode()
		fmt.Println("Client accepted:", clientID)

		client := models.Client{
			ID:       clientID,
			Conn:     conn,
			Messages: make(chan models.Message),
		}

		mutex.Lock()
		clients[client.ID] = client
		mutex.Unlock()

		go handleClient(client)
	}
}
