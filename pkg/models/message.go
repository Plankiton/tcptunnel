package models

import "time"

// Message struct para armazenar os metadados da mensagem
type Message struct {
	SendedAt   time.Time `json:"sendedAt,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	Message    string    `json:"message,omitempty"`
	SenderID   string    `json:"senderID,omitempty"`
	Action     Action    `json:"action,omitempty"`
	RetryCount int       `json:"retryCount,omitempty"`
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
