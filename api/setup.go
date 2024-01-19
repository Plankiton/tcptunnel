package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/plankiton/tcptunnel/pkg/env"
	"github.com/plankiton/tcptunnel/pkg/models"
	"github.com/plankiton/tcptunnel/pkg/services/client"
	"github.com/plankiton/tcptunnel/pkg/services/server"
)

var (
	wsHandler websocket.Upgrader

	wsMessagesStream chan models.Message // Response messages to websocket client
	wsErrorsStream   chan error          // Errors during websocket response processment

	messagesStream chan models.Message // Messages sent to tcptunnel
	responseStream chan models.Message // Response messages from tcptunnel
	errorStream    chan error          // Errors from tcptunnel

	wsConnStream   chan net.Conn              // Stream of websocket connections
	wsConnIsClosed map[net.Conn]chan struct{} // Stream to manage websocket connections closing

	isConnClosed chan struct{}
	connMutex    sync.Mutex
	conn         net.Conn

	serverURL string
	err       error
)

func init() {
	serverURL = "localhost:8080"
	go server.StartServer("8080")

	messagesStream = make(chan models.Message, env.MESSAGE_BATCH_SIZE)
	responseStream = make(chan models.Message, env.MESSAGE_BATCH_SIZE)
	errorStream = make(chan error, env.MESSAGE_BATCH_SIZE)

	isConnClosed = make(chan struct{}, 1)
	conn, err = net.Dial("tcp", serverURL)
	if err != nil {
		fmt.Printf("{\"error\": \"error connecting with local server: %v\"}\n", err)
	}

	go streamMessages()
	go listenMessages()
	go processResponseMessages()
	go handleWebsocketConnections()
}

func Healphcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello kitty!"))
}

func closeConnections() {
	connMutex.Lock()
	conn.Close()
	isConnClosed <- struct{}{}
	connMutex.Unlock()
}

func streamMessages() {
	defer closeConnections()

	for {
		processMessages(conn, messagesStream, responseStream, errorStream, isConnClosed)
	}
}

func listenMessages() {
	defer conn.Close()

}

func streamWebsocketMessages(conn net.Conn, isConnClosed chan struct{}) {
	for {
		for message := range wsMessagesStream {
			client.StreamMessage(message, conn, wsMessagesStream, wsErrorsStream, isConnClosed)
		}
	}
}

func waitWebsocketConnectionFinish(conn net.Conn) {

}

func handleWebsocketConnection(conn net.Conn) {
	wsConnStream <- conn
}

func handleWebsocketConnections() {
	for conn := range wsConnStream {
		go streamWebsocketMessages(conn, wsConnIsClosed[conn])
		go processResponseMessages()
	}
}

func processMessages(conn net.Conn, messagesStream, responseStream chan models.Message, errorStream chan error, isConnClosed chan struct{}) {
	select {
	case msg := <-messagesStream:
		messageJSON, _ := json.Marshal(msg)
		fmt.Printf("{\"info\": \"sending message to server\", \"message\": %s, \"serverURL\": %s}\n",
			string(messageJSON), serverURL)

		client.StreamMessage(msg, conn, responseStream, errorStream, isConnClosed)
		fmt.Printf("{\"info\": \"successful sent message\", \"message\": %s, \"serverURL\": %s}",
			string(messageJSON), serverURL)

	case _, isClosed := <-isConnClosed:
		if isClosed {
			return
		}

		close(isConnClosed)
		return
	}
}

func processResponseMessages() {
	for response := range responseStream {
		if err := <-errorStream; err != nil {
			wsErrorsStream <- models.Message{
				Message: fmt.Sprintf("error during tcptunnel processment: %v", err),
				Extra: map[string]interface{}{
					"response": response,
				},
			}

			continue
		}

		wsMessagesStream <- response
	}
}

func sendMessage(message models.Message, conn net.Conn, isConnClosed chan struct{}) (response models.Message, err error) {
	resMessage, err := client.SendMessage(message, conn, responseStream, isConnClosed)
	if err != nil {
		responseJson, _ := json.Marshal(responseError{
			Message: message,
			Error:   fmt.Sprint("error tunneling message request:", err),
			Extra: map[string]interface{}{
				"response": resMessage,
			},
		})

		fmt.Println(responseJson)
		return resMessage, err
	}

	return resMessage, nil
}

func NewHandler() http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/api/message", Message)
	router.HandleFunc("/api/stream", Stream)

	return http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}
}

type response struct {
	Info    string         `json:"info,omitempty"`
	Message models.Message `json:"message,omitempty"`
}

type responseError struct {
	Message models.Message `json:"message,omitempty"`
	Error   string         `json:"error,omitempty"`
	Extra   interface{}    `json:"extra,omitempty"`
}
