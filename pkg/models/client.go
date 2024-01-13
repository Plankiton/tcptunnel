package models

import "net"

// Client representa um cliente conectado
type Client struct {
	ID       string
	Conn     net.Conn
	Messages chan Message
}
