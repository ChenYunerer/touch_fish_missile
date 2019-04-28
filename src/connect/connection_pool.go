package connect

import (
	"sync"
)

type ConnectionPool struct {
	connections map[string]*Connection
	sync.RWMutex
}

var connectionPool *ConnectionPool

func initConnectionPool() *ConnectionPool {
	connectionPool = &ConnectionPool{
		connections: map[string]*Connection{},
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
	connPool.RWMutex.Lock()
	defer connPool.RWMutex.Unlock()
	connPool.connections[conn.RemoteAddress] = conn
}

func (connPool *ConnectionPool) RemoveConnection(conn *Connection) {
	connPool.RWMutex.Lock()
	defer connPool.RWMutex.Unlock()
	conn.Close()
	delete(connPool.connections, conn.RemoteAddress)
}

func (connPool *ConnectionPool) SendToOthers(me Connection, message []byte) {
	connPool.RWMutex.RLock()
	defer connPool.RWMutex.RUnlock()
	for remountAddress, conn := range connPool.connections {
		if remountAddress != me.RemoteAddress {
			//conn.SendMessageChan.c
			conn.SendMessageChan <- message
		}
	}
}
