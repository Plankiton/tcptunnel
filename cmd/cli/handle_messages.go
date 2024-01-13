package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/plankiton/tcptunnel/pkg/models"
)

func handleIncomingMessages(serverConn net.Conn, signalChannel chan error) {
	var msg models.Message
	decoder := json.NewDecoder(serverConn)
	err := decoder.Decode(&msg)
	if err != nil {
		signalChannel <- fmt.Errorf("Erro ao decodificar mensagem: %v", err)
		return
	}

	if !rawMessageFlagIsEnabled() {
		fmt.Print(msg.Message)
	} else {
		jsonMsg, _ := json.Marshal(msg)
		fmt.Print(string(jsonMsg))
	}

	signalChannel <- nil
}

func rawMessageFlagIsEnabled() bool {
	for _, arg := range os.Args {
		if "-r" == arg {
			return true
		}
	}

	return false
}

func getCLIActionFlag() (action string, messageId string) {
	messageId = ""
	action = models.GetActionFlag
	for a, arg := range os.Args {
		if arg == "-C" {
			action = models.CreateActionFlag
		}

		if arg == "-G" {
			if a+1 < len(os.Args) {
				action = models.GetActionFlag
				messageId = os.Args[a+1]
			} else {
				messageId = "0"
			}
		}

		if arg == "-U" {
			if a+1 < len(os.Args) {
				action = models.UpdateActionFlag
				messageId = os.Args[a+1]
			} else {
				messageId = "0"
			}
		}

		if arg == "-D" {
			if a+1 < len(os.Args) {
				action = models.DeleteActionFlag
				messageId = os.Args[a+1]
			} else {
				messageId = "0"
			}
		}

	}

	return action, messageId
}

func getServerURLFlag() string {
	serverURL := os.Getenv("TCPTUNNEL_URL")
	if serverURL != "" {
		return serverURL
	}

	for a, arg := range os.Args {
		if arg == "-s" {
			if a+1 < len(os.Args) {
				serverURL = os.Args[a+1]
				break
			}
		}
	}

	if serverURL != "" {
		return serverURL
	}
	return "localhost:8080"
}

func readUntilEOF() (string, error) {
	var message string
	fmt.Scan(&message)
	return message, nil
}
