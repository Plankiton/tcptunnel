package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/plankiton/tcptunnel/pkg/models"
)

func handleIncomingMessages(conn net.Conn) {
	for {
		var msg models.Message
		decoder := json.NewDecoder(conn)
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Erro ao decodificar mensagem:", err)
			return
		}

		fmt.Printf("[%s] Cliente %d: %s\n", msg.CreatedAt.Format("2006-01-02 15:04:05"), msg.SenderID, msg.Message)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Conectado ao servidor. ID do cliente:", os.Getpid())

	go handleIncomingMessages(conn)

	for {
		var text string
		fmt.Print("Digite uma mensagem: ")
		fmt.Scanln(&text)

		msg := models.Message{
			Message: text,
		}

		encoder := json.NewEncoder(conn)
		err := encoder.Encode(msg)
		if err != nil {
			fmt.Println("Erro ao enviar mensagem:", err)
			return
		}
	}
}
