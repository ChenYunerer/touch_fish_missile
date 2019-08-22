package server

import (
	"chat_group/src/config"
	"chat_group/src/connect"
	"chat_group/src/log"
	"net"
	"sync"
)

func StartServer() {
	go listenerConn()
}

func listenerConn() {
	conf := config.GetInstance()
	listener, err := net.Listen(conf.Network, conf.GetListenAddress())
	if err != nil {
		log.Error(err)
	}
	log.Info("Listen Start At: ", conf.GetListenAddress())
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Accept Conn Error: ", err)
		}
		log.Info("Accept A Conn, Address Is: ", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer log.Info("HandleConn Over")
	connection := connect.NewConnection(conn)
	connPool := connect.GetConnectionPoolInstant()
	connPool.AddConnection(connection)
	defer connPool.RemoveConnection(connection)
	quit := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer log.Info("ReadLoop Over")
		defer wg.Done()
		defer close(quit)
		readLoop(connection)
	}()
	wg.Add(1)
	go func() {
		defer log.Info("WriteLoop Over")
		defer wg.Done()
		writeLoop(connection, quit)
	}()
	wg.Wait()
}
