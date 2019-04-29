package connect

import (
	log "github.com/sirupsen/logrus"
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
	timeOut := time.NewTicker(time.Duration(1) * time.Second)
	for remountAddress, conn := range connPool.connections {
		if remountAddress != object.conn.RemoteAddress {
			select {
			case conn.SendMessageChan <- object.data:
				break
			case <-timeOut.C:
				close(object.conn.SendMessageChan)
				log.Info("超时 关闭连接")
				break
			default:
				log.Info("信道已经被关闭")
			}
		}
	}
}

func (connPool *ConnectionPool) AddConnection(conn *Connection) {
	connPool.Mutex.Lock()
	defer connPool.Mutex.Unlock()
	log.Info(conn.RemoteAddress, " 加入连接池")
	connPool.connections[conn.RemoteAddress] = conn
}

func (connPool *ConnectionPool) RemoveConnection(conn *Connection) {
	log.Info("RemoveConnection")
	connPool.Mutex.Lock()
	defer connPool.Mutex.Unlock()
	conn.Close()
	log.Info(conn.RemoteAddress, " 从连接池移除")
	delete(connPool.connections, conn.RemoteAddress)
}

func (connPool *ConnectionPool) PrepareSendToOther(object SendToOtherChanObject) {
	connPool.SendToOtherChan <- object
}
