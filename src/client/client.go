package client

import (
	"chat_group/src/config"
	"chat_group/src/conn_msg"
	"chat_group/src/connect"
	"chat_group/src/log"
	"chat_group/src/serialization"
	"net"
	"reflect"
	"sync"
	"time"
)

var token string
var group string

func StartClient(t string, g string) {
	token = t
	group = g
	go connectToServer()
}

func connectToServer() {
	conf := config.GetInstance()
	conn, err := net.DialTimeout(conf.Network, conf.GetServerAddress(), time.Duration(5)*time.Second)
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
	sendIntroduceMessage(connection)
	go listenCmd(connection)
	wg.Wait()
}

func sendIntroduceMessage(conn *connect.Connection) {
	introduceMessage := conn_msg.NewIntroduceMessage(token, group, "")
	t := reflect.TypeOf(introduceMessage)
	messageId := conn_msg.MessageTypeIdMap[t]
	bytes, err := serialization.EncodeMessage(&introduceMessage, messageId[:])
	if err != nil {
		log.Error(err)
	}
	conn.SendMessageChan <- bytes
}
