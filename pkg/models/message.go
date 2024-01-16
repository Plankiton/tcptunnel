package models

import "time"

// Message struct para armazenar os metadados da mensagem
type Message struct {
	ID string `json:"id,omitempty"`

	RetryCount int `json:"retryCount,omitempty"`

	Action Action `json:"action,omitempty"`

	Message  string `json:"message,omitempty"`
	SenderID string `json:"senderID,omitempty"`

	Extra interface{} `json:"extra,omitempty"`

	SendedAt  time.Time `json:"sendedAt,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type Action struct {
	Name      string `json:"name"`
	MessageID string `json:"messageID"`
}

func (m Message) Error() string {
	return m.Message
}

const (
	UpdateActionFlag = "update"
	CreateActionFlag = "create"
	DeleteActionFlag = "delete"
	GetActionFlag    = "get"
)
