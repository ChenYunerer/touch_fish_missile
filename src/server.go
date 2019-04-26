package main

import (
	"bufio"
	"chat_group/src/config"
	"chat_group/src/conn_msg"
	"chat_group/src/connect"
	"chat_group/src/serialization"
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
	"sync"
	"time"
)

func StartServer() {
	go listenerConn()
}

func listenerConn() {
	conf := config.GetInstance()
	listener, err := net.Listen(conf.Network, conf.GetAddress())
	if err != nil {
		log.Error(err)
	}
	log.Info("listen start")
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
		}
		log.Info("accept a conn")
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	connection := connect.NewConnection(conn)
	connPool := connect.GetConnectionPoolInstant()
	connPool.AddConnection(connection)
	defer connection.Conn.Close()
	defer connPool.RemoveConnection(connection)
	log.Info("handle conn address is ", connection.RemoteAddress)
	quit := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(quit)
		readLoop(connection)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		writeLoop(connection, quit)
	}()
	wg.Wait()
}

func readLoop(conn *connect.Connection) {
	conf := config.GetInstance()
	bytes := make([]byte, 1024)
	for {
		if conf.ReadTimeout.Seconds() != 0 {
			deadline := time.Now().Add(conf.WriteTimeout)
			err := conn.Conn.SetReadDeadline(deadline)
			if err != nil {
				log.Error(err)
			}
		}
		reader := bufio.NewReader(conn.Conn)
		n, err := reader.Read(bytes)
		if err != nil {
			log.Error(err)
			conn.AddRetryTimes()
			connRetryTimes := conn.GetRetryTimes()
			log.Info("conn ", conn.RemoteAddress, " retry times is ", connRetryTimes)
			if connRetryTimes >= conf.RetryTimes {
				return
			}
			continue
		}
		if n == 0 {
			log.Error("no data read from reader")
			return
		}
		messageId := conn_msg.GetMessageIdFromMessageBytes(bytes)
		messageType := conn_msg.GetMessageTypeByMessageId(messageId)
		messageInterface, err := serialization.DecodeMessage(messageType, bytes, conn_msg.LenOfMessageID)
		message := messageInterface.(conn_msg.Message)
		if err != nil {
			log.Info(err)
			continue
		}
		log.Info("receive conn_msg from client ", message)
		err = message.HandleMessage(conn)
		if err != nil {
			log.Error(err)
		}
	}
}

func writeLoop(conn *connect.Connection, quit chan struct{}) {
	conf := config.GetInstance()
	pingTimer := time.NewTicker(conf.PingDuration)
	for {
		if conf.WriteTimeout.Seconds() != 0 {
			deadline := time.Now().Add(conf.WriteTimeout)
			err := conn.Conn.SetWriteDeadline(deadline)
			if err != nil {
				log.Error(err)
			}
		}
		select {
		case messageBytes := <-conn.SendMessageChan:
			log.Info("send conn_msg to ", conn.RemoteAddress)
			n, err := conn.Conn.Write(messageBytes)
			if err != nil {
				log.Error(err)
				conn.AddRetryTimes()
				connRetryTimes := conn.GetRetryTimes()
				log.Info("conn ", conn.RemoteAddress, " retry times is ", connRetryTimes)
				if connRetryTimes >= conf.RetryTimes {
					return
				}
				continue
			} else {
				//发送消息不重制RetryTimes
				//conn.ResetRetryTimes()
			}
			if n == 0 {
				log.Error("send data error")
				return
			}
		case <-pingTimer.C:
			pingMessage := conn_msg.NewPingMessage()
			t := reflect.TypeOf(pingMessage)
			messageId := conn_msg.MessageTypeIdMap[t]
			bytes, err := serialization.EncodeMessage(&pingMessage, messageId[:])
			if err != nil {
				log.Error(err)
				continue
			}
			conn.SendMessageChan <- bytes
		case <-quit:
			return
		}
	}
}
