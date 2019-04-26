package connect

import (
	"net"
	"sync"
)

const SendMessageChanBuffer = 32

type Connection struct {
	RemoteAddress   string
	Conn            net.Conn
	SendMessageChan chan []byte
	RetryTimes      uint32
	sync.RWMutex
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		RemoteAddress:   conn.RemoteAddr().String(),
		Conn:            conn,
		SendMessageChan: make(chan []byte, SendMessageChanBuffer),
		RetryTimes:      0,
	}
}

func (conn *Connection) ResetRetryTimes() {
	if conn.RetryTimes == 0 {
		return
	}
	conn.RLock()
	defer conn.RUnlock()
	conn.RetryTimes = 0
}

func (conn *Connection) GetRetryTimes() uint32 {
	conn.RLock()
	defer conn.RUnlock()
	return conn.RetryTimes
}

func (conn *Connection) AddRetryTimes() {
	conn.Lock()
	defer conn.Unlock()
	conn.RetryTimes++
}
