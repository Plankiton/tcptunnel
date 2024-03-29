package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/plankiton/tcptunnel/pkg/models"
	"github.com/plankiton/tcptunnel/pkg/services/client"
)

var (
	responseStream chan models.Message
)

func main() {
	serverURL := getServerURLFlag()

	action, messageId := getCLIActionFlag()

	conn, err := net.Dial("tcp", serverURL)
	if err != nil {
		fmt.Printf("{\"error\": \"error connecting with local server: %v\"}", err)
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

	response, err := client.SendMessageToServer(msg, conn, responseStream)
	if err != nil {
		fmt.Printf("{\"error\": \"error during client trigger startup: %v\"}", err)
		os.Exit(1)
	}

	fmt.Println(response.Message)
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

func inputIsNeeded(msg models.Message) bool {
	return msg.Action.Name == models.CreateActionFlag || msg.Action.Name == models.UpdateActionFlag
}
