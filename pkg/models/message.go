package models

import "time"

// Message struct para armazenar os metadados da mensagem
type Message struct {
	SendedAt  time.Time `json:"sendedAt"`
	CreatedAt time.Time `json:"createdAt"`
	Message   string    `json:"message"`
	SenderID  int       `json:"senderID"`
}
