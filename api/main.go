package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/plankiton/tcptunnel/pkg/models"
	"github.com/plankiton/tcptunnel/pkg/services/client"
	"github.com/plankiton/tcptunnel/pkg/services/server"
)

var (
	messagesStream chan models.Message
	responseStream chan models.Message
	conn           net.Conn
	err            error
)

func init() {
	go server.StartServer("8080")
	go streamMessagesToServer("localhost:8080")
}

func streamMessagesToServer(serverURL string) {
	conn, err = net.Dial("tcp", serverURL)
	if err != nil {
		fmt.Printf("{\"error\": \"error connecting with local server: %v\"}\n", err)
	}
	defer conn.Close()

	for msg := range messagesStream {
		messageJSON, _ := json.Marshal(msg)
		fmt.Printf("{\"info\": \"sending message to server\", \"message\": %s, \"serverURL\": %s}\n",
			string(messageJSON), serverURL)

		response, err := client.SendMessageToServer(msg, conn, responseStream)
		if err != nil {
			responseJSON, _ := json.Marshal(response)
			fmt.Printf("{\"error\": \"error tunneling message to tcp server: %v\", \"message\": %s, \"serverURL\": %s, \"serverResponse\": %s}\n",
				err, string(messageJSON), serverURL, string(responseJSON))
			continue
		}

		fmt.Printf("{\"info\": \"successful sent message\", \"message\": %s, \"serverURL\": %s}",
			string(messageJSON), serverURL)
	}
}

func Api(w http.ResponseWriter, r *http.Request) {
	var message models.Message

	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		encoder.Encode(responseError{
			Message: message,
			Error:   fmt.Sprint("error decoding request body:", err),
		})

		w.WriteHeader(400)
		return
	}

	resMessage, err := client.SendMessageToServer(message, conn, responseStream)
	if err != nil {
		encoder.Encode(responseError{
			Message: message,
			Error:   fmt.Sprint("error tunneling message request:", err),
		})

		w.WriteHeader(500)
		return
	}

	encoder.Encode(response{
		Info:    "message successful sended",
		Message: resMessage,
	})
}

type response struct {
	Info    string         `json:"info"`
	Message models.Message `json:"message"`
}

type responseError struct {
	Message models.Message `json:"message"`
	Error   string         `json:"error"`
}
