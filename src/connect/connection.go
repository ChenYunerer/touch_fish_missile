package connect

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"touch_fish_missile/src/log"
)

const SendMessageChanBuffer = 64

type ConnectionUserInfo struct {
	Token string
	Group string
}

type Connection struct {
	RemoteAddress   string
	Conn            net.Conn
	SendMessageChan chan []byte
	RetryTimes      uint32
	Buffer          *bytes.Buffer
	Log             *logrus.Entry
	ConnectionUserInfo
	sync.RWMutex
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		RemoteAddress:   conn.RemoteAddr().String(),
		Conn:            conn,
		SendMessageChan: make(chan []byte, SendMessageChanBuffer),
		RetryTimes:      0,
		Buffer:          &bytes.Buffer{},
		Log:             logrus.WithFields(logrus.Fields{"Origin": "Connection-" + conn.RemoteAddr().String()}),
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

func (conn *Connection) Close() {
	conn.Buffer = &bytes.Buffer{}
	err := conn.Conn.Close()
	if err != nil {
		log.Error(err)
	}
	log.Info("conn ", conn.RemoteAddress, " closed")
}
