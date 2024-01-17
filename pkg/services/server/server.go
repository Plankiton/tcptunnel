package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/plankiton/tcptunnel/pkg/models"
)

var (
	mutex = sync.Mutex{}

	clients  = make(map[string]models.Client)
	messages = make([]models.Message, 0)

	messageStream = make(chan models.Message, 5)
	errorsStream  = make(chan error)

	responseStream       = make(chan models.Message, 5)
	responseErrorsStream = make(chan error)
)

func isRelativeFromLastMessageId(messageId, messageIndex int) bool {
	return getMessageListIndex(messageId) == messageIndex
}

func getMessageListIndex(messageId int) int {
	return (len(messages) - 1) - messageId
}

func handleClient(client models.Client) {
	defer client.Conn.Close()

	for {
		var (
			message models.Message
		)

		decoder := json.NewDecoder(client.Conn)
		err := decoder.Decode(&message)

		if err != nil {
			messageJSON, _ := json.Marshal(message)
			errorMessage := fmt.Sprintf("error during client message encoding: %v", err)

			fmt.Printf("{\"message\": %s, \"error\": \"%s\"}\n", messageJSON, errorMessage)
			messageStream <- message
			errorsStream <- models.Message{
				Message: errorMessage,
				Extra: map[string]interface{}{
					"message": message,
				},
			}
		}

		message.CreatedAt = message.SendedAt
		message.SenderID = client.ID

		rawJsonMessage, _ := json.Marshal(message)
		fmt.Println(string(rawJsonMessage))

		outputMsg := models.Message{
			SendedAt: time.Now(),
			SenderID: client.ID,
		}

		mutex.Lock()
		switch message.Action.Name {
		case models.GetActionFlag:
			selectedMsg := selectMessage(message)

			if selectedMsg.Message == "" {
				selectedMsg.Message = "message not found, please select each other message id"
			}

			outputMsg.Message = selectedMsg.Message

		case models.UpdateActionFlag:
			messageIndex := getMessageListIndex(getMessageID(message))
			message.CreatedAt = time.Now()

			messages[messageIndex] = message
		case models.DeleteActionFlag:
			messageIndex := getMessageListIndex(getMessageID(message))

			swapMessages := make([]models.Message, 0)
			for i, message := range messages {
				if i != messageIndex {
					swapMessages = append(swapMessages, message)
				}
			}

			messages = swapMessages
		case models.CreateActionFlag:
			message.CreatedAt = time.Now()

			outputMsg.Message = message.Message
			messages = append(messages, message)
		}

		mutex.Unlock()

		sendMessage(client, outputMsg)
		messageStream <- outputMsg
		errorsStream <- nil
	}
}

func streamMessages() {
	for outputMsg := range messageStream {
		err := <-errorsStream
		errorMessage := fmt.Sprintf("error receiving message: %v", err)
		messageJSON, _ := json.Marshal(outputMsg)

		fmt.Println("{\"message\": " + string(messageJSON) + ", \"error\": \"erro ao decodificar mensagem: " + err.Error() + "\"}")
		if err != nil {
			responseStream <- outputMsg
			responseErrorsStream <- models.Message{
				Message: errorMessage,
				Extra: map[string]interface{}{
					"message": outputMsg,
				},
			}
		}

		for _, currClient := range clients {
			if currClient.ID == outputMsg.SenderID {
				continue
			}

			if err = sendMessage(currClient, outputMsg); err != nil {
				if err != nil && outputMsg.RetryCount < 5 {
					messageStream <- outputMsg
					errorsStream <- nil
				}

				break
			}
		}

		outputMsg.RetryCount++
	}

}

func formatClientID(id int) string {
	return fmt.Sprintf("%x", id)
}

func generateRandomHexCode() (string, error) {
	// Defina o número de bytes que você precisa para o código hexadecimal (neste caso, 3 bytes para 5 caracteres hexadecimais)
	numBytes := 3

	// Crie um slice de bytes para armazenar os bytes aleatórios
	randomBytes := make([]byte, numBytes)

	// Leia bytes aleatórios da fonte de entropia criptográfica
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Converta os bytes para uma string hexadecimal
	hexCode := hex.EncodeToString(randomBytes)

	// Pegue os primeiros 5 caracteres do código hexadecimal
	hexCode = hexCode[:5]

	return hexCode, nil
}

func getMessageID(message models.Message) int {
	if messageID, err := strconv.Atoi(message.Action.MessageID); err != nil {
		return messageID
	}

	return 0
}

func selectMessage(message models.Message) models.Message {
	messageID := getMessageID(message)
	fmt.Println("\n\nMessages:", messages, "message id:", messageID)

	var found models.Message
	for m, message := range messages {
		fmt.Println(" -> sel message index:", m, " query message index:", getMessageListIndex(messageID), "message:", message.Message)
		if isRelativeFromLastMessageId(messageID, m) {
			fmt.Print("     FOUND!!\n\n\n")
			found = message
			break
		}

		fmt.Print("			NOT FOUND!!\n\n\n")
	}

	return found
}

func sendMessage(client models.Client, message models.Message) error {
	encoder := json.NewEncoder(client.Conn)
	err := encoder.Encode(message)
	if err != nil {
		fmt.Println("Erro ao enviar mensagem para o cliente", client.ID, ":", err)
		delete(clients, client.ID)
		return err
	}

	return nil
}

func StartServer(serverPort string) {
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
