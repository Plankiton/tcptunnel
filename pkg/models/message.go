package models

import "time"

// Message struct para armazenar os metadados da mensagem
type Message struct {
	SendedAt   time.Time `json:"sendedAt"`
	CreatedAt  time.Time `json:"createdAt"`
	Message    string    `json:"message"`
	SenderID   string    `json:"senderID"`
	Action     Action    `json:"action"`
	RetryCount int       `json:"retryCount"`
}

type Action struct {
	Name      string `json:"name"`
	MessageID string `json:"messageID"`
}

const (
	UpdateActionFlag = "update"
	CreateActionFlag = "create"
	DeleteActionFlag = "delete"
	GetActionFlag    = "get"
)
