package connect

import "net"

type Connection struct {
	RemoteAddress   string
	Conn            net.Conn
	SendMessageChan chan []byte
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		RemoteAddress:   conn.RemoteAddr().String(),
		Conn:            conn,
		SendMessageChan: make(chan []byte, 32),
	}
}
