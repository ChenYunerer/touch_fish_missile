package client

import (
	"chat_group/src/config"
	"chat_group/src/connect"
	"chat_group/src/log"
	"net"
	"sync"
	"time"
)

var token string

func StartClient(t string) {
	token = t
	go connectToServer()
}

func connectToServer() {
	conf := config.GetInstance()
	conn, err := net.DialTimeout(conf.Network, conf.GetAddress(), time.Duration(5)*time.Second)
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()
	handleConn(conn)
}

func handleConn(conn net.Conn) {
	defer log.Info("HandleConn Over")
	connection := connect.NewConnection(conn)
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
		writeLoop(connection, token, quit)
	}()
	wg.Wait()
}
