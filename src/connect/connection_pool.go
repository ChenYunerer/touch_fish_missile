package connect

import (
	"sync"
)

type ConnectionPool struct {
	Connections map[string]*Connection
	sync.Mutex
}

var connectionPool *ConnectionPool

func initConnectionPool() *ConnectionPool {
	connectionPool = &ConnectionPool{
		Connections: map[string]*Connection{},
	}
	return connectionPool
}

func GetConnectionPoolInstant() *ConnectionPool {
	if connectionPool == nil {
		initConnectionPool()
	}
	return connectionPool
}

func (connPool *ConnectionPool) AddConnection(conn *Connection) {
	connPool.Mutex.Lock()
	defer connPool.Mutex.Unlock()
	connPool.Connections[conn.RemoteAddress] = conn
}

func (connPool *ConnectionPool) RemoveConnection(conn *Connection) {
	connPool.Mutex.Lock()
	defer connPool.Mutex.Unlock()
	delete(connPool.Connections, conn.RemoteAddress)
}

func (connPool *ConnectionPool) SendToOthers(me Connection, message []byte) {
	for remountAddress, conn := range connPool.Connections {
		if remountAddress != me.RemoteAddress {
			conn.SendMessageChan <- message
		}
	}
}
