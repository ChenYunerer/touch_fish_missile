package connect

import (
	"chat_group/src/config"
	"chat_group/src/log"
	"strings"
	"sync"
	"time"
)

type SendToOtherChanObject struct {
	data []byte
	conn *Connection
}

func NewSendToOtherChanObject(b []byte, conn *Connection) SendToOtherChanObject {
	return SendToOtherChanObject{
		data: b,
		conn: conn,
	}
}

type ConnectionPool struct {
	connections     map[string]*Connection
	SendToOtherChan chan SendToOtherChanObject
	sync.Mutex
}

var connectionPool *ConnectionPool

func startConnectionPool() {
	connectionPool = &ConnectionPool{
		connections:     map[string]*Connection{},
		SendToOtherChan: make(chan SendToOtherChanObject, 1024),
	}
	go func() {
		for {
			select {
			case object := <-connectionPool.SendToOtherChan:
				connectionPool.sendToOthers(object)
			}
		}
	}()
}

func GetConnectionPoolInstant() *ConnectionPool {
	if connectionPool == nil {
		startConnectionPool()
	}
	return connectionPool
}

func (connPool *ConnectionPool) sendToOthers(object SendToOtherChanObject) {
	conf := config.GetInstance()
	timeOut := time.NewTicker(conf.WriteTimeout)
	for remountAddress, conn := range connPool.connections {
		if !strings.EqualFold(conn.Group, object.conn.Group) {
			continue
		}
		if remountAddress != object.conn.RemoteAddress {
			select {
			case conn.SendMessageChan <- object.data:
				break
			case <-timeOut.C:
				close(object.conn.SendMessageChan)
				log.Info("timeout, close connection")
				break
			default:
				log.Info("chan has bean closed")
			}
		}
	}
}

func (connPool *ConnectionPool) AddConnection(conn *Connection) {
	connPool.Mutex.Lock()
	defer connPool.Mutex.Unlock()
	log.Info(conn.RemoteAddress, " add into connection pool")
	connPool.connections[conn.RemoteAddress] = conn
}

func (connPool *ConnectionPool) RemoveConnection(conn *Connection) {
	log.Info("RemoveConnection")
	connPool.Mutex.Lock()
	defer connPool.Mutex.Unlock()
	conn.Close()
	log.Info(conn.RemoteAddress, " remove from connection poll")
	delete(connPool.connections, conn.RemoteAddress)
}

func (connPool *ConnectionPool) PrepareSendToOther(object SendToOtherChanObject) {
	connPool.SendToOtherChan <- object
}

func (connPool *ConnectionPool) SearchConnectionsByGroup(group string) []Connection {
	//todo
}
