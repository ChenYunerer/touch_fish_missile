package client

import (
	"chat_group/src/config"
	"chat_group/src/conn_msg"
	"chat_group/src/connect"
	"chat_group/src/log"
	"chat_group/src/serialization"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"
)

var token string
var groupTag string

func StartClient(t string, gt string) {
	token = t
	groupTag = gt
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
	go listenCmd(connection)
	wg.Wait()
}

var cmdInputStr string

func listenCmd(conn *connect.Connection) {
	for {
		num, err := fmt.Scanln(&cmdInputStr)
		if num == 0 || err != nil {
			log.Error("listenCmd err num: ", string(num), " err: ", err)
			continue
		}
		stringMessage := conn_msg.NewStringMessage(token, cmdInputStr)
		t := reflect.TypeOf(stringMessage)
		messageId := conn_msg.MessageTypeIdMap[t]
		bytes, err := serialization.EncodeMessage(&stringMessage, messageId[:])
		if err != nil {
			log.Error(err)
			continue
		}
		conn.SendMessageChan <- bytes
	}
}
