package main

import (
	"bufio"
	"chat_group/src/config"
	"chat_group/src/conn_msg"
	"chat_group/src/connect"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"time"
)


func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func main() {
	go listenerConn()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Interrupt")
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
	type BaseMessage  struct {
		Content struct{
			MessageType string
		}
	}
	bytes := make([]byte, 1024)
	for {
		reader := bufio.NewReader(conn.Conn)
		n, err := reader.Read(bytes)
		if err != nil {
			log.Error(err)
			return
		}
		if n == 0 {
			log.Error("no data read from reader")
			return
		}
		connect.GetConnectionPoolInstant().SendToOthers(conn, bytes)
		baseMessage := &BaseMessage{}
		err = json.Unmarshal(bytes[:n], baseMessage)
		if err != nil {
			log.Error(err)
		}
		log.Info("receive conn_msg from client ", baseMessage)
		messageType := conn_msg.MessageContentMap[baseMessage.Content.MessageType]
		message := reflect.New(messageType)
		err = json.Unmarshal(bytes[:n], message)
		if err != nil {
			log.Error(err)
		}

	}
}

func writeLoop(conn *connect.Connection, quit chan struct{}) {
	conf := config.GetInstance()
	pingTimer := time.NewTicker(conf.PingDuration)
	for {
		select {
		case messageJsonBytes := <-conn.SendMessageChan:
			messageJsonStr := string(messageJsonBytes)
			log.Info("send conn_msg to ", conn.RemoteAddress, " conn_msg : "+messageJsonStr)
			n, err := conn.Conn.Write(messageJsonBytes)
			if err != nil {
				log.Error(err)
				return
			}
			if n == 0 {
				log.Error("send data error")
				return
			}
		case <-pingTimer.C:
			pingMessage := conn_msg.NewPingMessage()
			messageJsonBytes, err := json.Marshal(pingMessage)
			if err != nil {
				log.Error(err)
			}
			conn.SendMessageChan <- messageJsonBytes
		case <-quit:
			return
		}
	}
}
