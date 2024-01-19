package models

import (
	"net"

	"golang.org/x/net/context"
)

// Client representa um cliente conectado
type Client struct {
	ID       string
	Conn     net.Conn
	Messages chan Message

	Ctx context.Context
}
