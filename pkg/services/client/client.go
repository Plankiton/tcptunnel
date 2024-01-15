package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/plankiton/tcptunnel/pkg/models"
)

func SendMessageToServer(msg models.Message, conn net.Conn, responseStream chan models.Message) (response models.Message, err error) {
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(msg)
	if err != nil {
		return response, fmt.Errorf("error decoding message: %v", err)
	}

	errorChan := make(chan error, 5)
	go handleIncomingMessages(conn, responseStream, errorChan)
	go func() {
		time.Sleep(30 * time.Second)
		errorChan <- fmt.Errorf("internal time out when receiving server incoming messages")
	}()

	response = <-responseStream
	err = <-errorChan
	return response, err
}

func receiveServerMessage(conn net.Conn) (models.Message, error) {
	var message models.Message
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&message)
	return message, err
}

func handleIncomingMessages(serverConn net.Conn, responseStream chan models.Message, signalChannel chan error) {
	for {
		msg, err := receiveServerMessage(serverConn)
		if err != nil {
			signalChannel <- fmt.Errorf("Erro ao decodificar mensagem: %v", err)
			return
		}

		responseStream <- msg
		signalChannel <- nil
	}
}

func rawMessageFlagIsEnabled() bool {
	for _, arg := range os.Args {
		if "-r" == arg {
			return true
		}
	}

	return false
}
