package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/plankiton/tcptunnel/pkg/models"
)

// Client representa um cliente conectado
type Client struct {
	ID       int
	Conn     net.Conn
	Messages chan models.Message
}

var (
	clients  = make(map[int]Client)
	messages = make([]models.Message, 0)
	mutex    = sync.Mutex{}
)

func handleClient(client Client) {
	defer client.Conn.Close()

	for {
		var msg models.Message
		decoder := json.NewDecoder(client.Conn)
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Erro ao decodificar mensagem:", err)
			return
		}

		msg.SendedAt = time.Now()
		msg.CreatedAt = msg.SendedAt
		msg.SenderID = client.ID

		mutex.Lock()
		messages = append(messages, msg)
		mutex.Unlock()

		for _, otherClient := range clients {
			if otherClient.ID != client.ID {
				encoder := json.NewEncoder(otherClient.Conn)
				err := encoder.Encode(msg)
				if err != nil {
					fmt.Println("Erro ao enviar mensagem para o cliente", otherClient.ID, ":", err)
				}
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor escutando na porta 8080...")

	clientID := 1

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		client := Client{
			ID:       clientID,
			Conn:     conn,
			Messages: make(chan models.Message),
		}

		mutex.Lock()
		clients[clientID] = client
		mutex.Unlock()

		clientID++

		go handleClient(client)
	}
}
