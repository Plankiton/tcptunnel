package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/plankiton/tcptunnel/pkg/models"
)

func main() {
	serverURL := getServerURLFlag()

	action, messageId := getCLIActionFlag()

	conn, err := net.Dial("tcp", serverURL)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		os.Exit(1)
	}
	defer conn.Close()

	msg := models.Message{
		SendedAt: time.Now(),
		Action: models.Action{
			Name:      action,
			MessageID: messageId,
		},
	}

	if inputIsNeeded(msg) {
		fmt.Scan(&msg.Message)
	}

	encoder := json.NewEncoder(conn)
	err = encoder.Encode(msg)
	if err != nil {
		fmt.Println("Erro ao enviar mensagem:", err)
		return
	}

	errorChan := make(chan error, 5)
	go handleIncomingMessages(conn, errorChan)
	go func() {
		time.Sleep(30 * time.Second)
		errorChan <- fmt.Errorf("Time out")
	}()

	err = <-errorChan
	if err != nil {
		fmt.Println("Erro ao receber a resposta do servidor:", err)
	}
}

func inputIsNeeded(msg models.Message) bool {
	return msg.Action.Name == models.CreateActionFlag || msg.Action.Name == models.UpdateActionFlag
}
