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
	clients          = make(map[string]models.Client)
	messages         = make([]models.Message, 0)
	mutex            = sync.Mutex{}
	messageStreaming = make(chan models.Message, 5)
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
		var msg models.Message
		decoder := json.NewDecoder(client.Conn)
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Erro ao decodificar mensagem:", err)
			return
		}

		msg.CreatedAt = msg.SendedAt
		msg.SenderID = client.ID

		rawJsonMessage, _ := json.Marshal(msg)
		fmt.Println(string(rawJsonMessage))

		outputMsg := models.Message{
			SendedAt: time.Now(),
			SenderID: client.ID,
		}

		mutex.Lock()
		switch msg.Action.Name {
		case models.GetActionFlag:
			selectedMsg := selectMessage(msg)

			if selectedMsg.Message == "" {
				selectedMsg.Message = "message not found, please select each other message id"
			}

			outputMsg.Message = selectedMsg.Message

		case models.UpdateActionFlag:
			messageIndex := getMessageListIndex(getMessageID(msg))
			msg.CreatedAt = time.Now()

			messages[messageIndex] = msg
		case models.DeleteActionFlag:
			messageIndex := getMessageListIndex(getMessageID(msg))

			swapMessages := make([]models.Message, 0)
			for i, msg := range messages {
				if i != messageIndex {
					swapMessages = append(swapMessages, msg)
				}
			}

			messages = swapMessages
		case models.CreateActionFlag:
			msg.CreatedAt = time.Now()

			outputMsg.Message = msg.Message
			messages = append(messages, msg)
		}

		mutex.Unlock()

		fmt.Println()
		sendMessage(client, outputMsg)
		messageStreaming <- outputMsg
	}
}

func streamMessages() {
	for outputMsg := range messageStreaming {
		hasError := false

		for _, currClient := range clients {
			if currClient.ID == outputMsg.SenderID {
				continue
			}

			if err := sendMessage(currClient, outputMsg); err != nil {
				hasError = true
				break
			}
		}

		if hasError && outputMsg.RetryCount < 5 {
			messageStreaming <- outputMsg
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

func getMessageID(msg models.Message) int {
	if messageID, err := strconv.Atoi(msg.Action.MessageID); err != nil {
		return messageID
	}

	return 0
}

func selectMessage(msg models.Message) models.Message {
	messageID := getMessageID(msg)
	fmt.Println("\n\nMessages:", messages, "message id:", messageID)

	var found models.Message
	for m, msg := range messages {
		fmt.Println(" -> sel message index:", m, " query msg index:", getMessageListIndex(messageID), "msg:", msg.Message)
		if isRelativeFromLastMessageId(messageID, m) {
			fmt.Print("     FOUND!!\n\n\n")
			found = msg
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
