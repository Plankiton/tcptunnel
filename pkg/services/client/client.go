package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/plankiton/tcptunnel/pkg/env"
	"github.com/plankiton/tcptunnel/pkg/models"
)

func SendMessage(message models.Message, conn net.Conn, responseStream chan models.Message, isConnClosed chan struct{}) (response models.Message, err error) {
	errorStream := make(chan error, env.MESSAGE_BATCH_SIZE)
	go StreamMessage(message, conn, responseStream, errorStream, isConnClosed)

	response = <-responseStream
	err = <-errorStream
	return response, err
}

func closeConnection(conn net.Conn, isConnClosed chan struct{}) {
	conn.Close()
	isConnClosed <- struct{}{}
}

func StreamMessage(message models.Message, conn net.Conn, responseStream chan models.Message, errorStream chan error, isConnClosed chan struct{}) {
	defer closeConnection(conn, isConnClosed)

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	err := encoder.Encode(message)
	if err != nil {
		responseStream <- models.Message{
			Message: fmt.Sprintf("error during message streaming: %v", err),
			Extra: map[string]interface{}{
				"message": message,
			},
		}
		errorStream <- fmt.Errorf("error decoding message: %v", err)
		return
	}

	var response models.Message
	err = decoder.Decode(&response)
	if err != nil {
		responseStream <- response
		errorStream <- fmt.Errorf("error decoding message: %v", err)
		return
	}

	responseStream <- response
	errorStream <- nil
}

func ListenMessages(conn net.Conn, messageStream, responseStream chan models.Message, errorStream chan error, isConnClosed chan struct{}) {
	go processResponseMessages(conn, messageStream, responseStream, errorStream, isConnClosed)
	go func() {
		time.Sleep(30 * time.Second)
		errorStream <- fmt.Errorf("internal time out when receiving server incoming messages")
	}()
}

func ListenMessageMessages(conn net.Conn, responseStream chan models.Message, errorStream chan error) {
	go handleIncomingMessages(conn, responseStream, errorStream)
	go func() {
		time.Sleep(30 * time.Second)
		errorStream <- fmt.Errorf("internal time out when receiving server incoming messages")
	}()
}

func receiveServerMessage(conn net.Conn) (models.Message, error) {
	var message models.Message
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&message)
	return message, err
}

func processResponseMessages(serverConn net.Conn, messageStream, responseStream chan models.Message, errorsStream chan error, isConnClosed chan struct{}) {
	for message := range messageStream {
		go StreamMessage(message, serverConn, responseStream, errorsStream, isConnClosed)
		go handleIncomingMessages(serverConn, responseStream, errorsStream)
	}
}

func handleIncomingMessages(serverConn net.Conn, responseStream chan models.Message, errorsStream chan error) {
	for {
		message, err := receiveServerMessage(serverConn)
		if err != nil {
			errorsStream <- fmt.Errorf("err receiving server message: %v", err)
			return
		}

		responseStream <- message
		errorsStream <- nil
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
