package main

import (
	"bufio"
	"bytes"
	"chat_group/src/config"
	"chat_group/src/conn_msg"
	"chat_group/src/connect"
	"chat_group/src/serialization"
	"errors"
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
	log.Info("Listen Start")
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
	defer log.Info("handleConn 结束")
	connection := connect.NewConnection(conn)
	connPool := connect.GetConnectionPoolInstant()
	connPool.AddConnection(connection)
	defer connPool.RemoveConnection(connection)
	quit := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer log.Info("read loop 结束")
		defer wg.Done()
		defer close(quit)
		readLoop(connection)
	}()
	wg.Add(1)
	go func() {
		defer log.Info("write loop 结束")
		defer wg.Done()
		writeLoop(connection, quit)
	}()
	wg.Wait()
}

func readLoop(conn *connect.Connection) {
	conf := config.GetInstance()
	bytes := make([]byte, 1024)
	reader := bufio.NewReader(conn.Conn)
	for {
		if conf.ReadTimeout.Seconds() != 0 {
			deadline := time.Now().Add(conf.WriteTimeout)
			err := conn.Conn.SetReadDeadline(deadline)
			if err != nil {
				log.Error(err)
			}
		}
		n, err := reader.Read(bytes)
		if err != nil {
			log.Error(err)
			conn.AddRetryTimes()
			connRetryTimes := conn.GetRetryTimes()
			log.Info("Read conn ", conn.RemoteAddress, " retry times is ", connRetryTimes)
			if connRetryTimes >= conf.RetryTimes {
				return
			}
			continue
		}
		if n == 0 {
			log.Error("no data read from reader")
			return
		}
		_, err = conn.Buffer.Write(bytes[:n])
		if err != nil {
			log.Error(err)
			return
		}
		bytess, err := decodeData(conn.Buffer)
		if err != nil {
			log.Error(err)
			return
		}
		messages, err := decodeMessage(bytess)
		if err != nil {
			log.Error(err)
			continue
		}
		for _, message := range messages {
			log.Info("receive conn_msg from client ", message)
			err = message.HandleMessage(conn)
			if err != nil {
				log.Error(err)
			}
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
		case messageBytes, ok := <-conn.SendMessageChan:
			if !ok {
				log.Info("SendMessageChan Closed")
				return
			}
			log.Info("Send Message To ", conn.RemoteAddress)
			n, err := conn.Conn.Write(messageBytes)
			if err != nil {
				log.Error(err)
				conn.AddRetryTimes()
				connRetryTimes := conn.GetRetryTimes()
				log.Info("Write Conn ", conn.RemoteAddress, " Retry Times Is ", connRetryTimes)
				if connRetryTimes >= conf.RetryTimes {
					return
				}
				continue
			}
			if n == 0 {
				log.Error("Send Data Error")
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

func decodeData(buf *bytes.Buffer) ([][]byte, error) {
	dataArray := [][]byte{}
	lenOfMessageID := conn_msg.LenOfMessageID
	lenOfLength := uint64(8)
	for uint64(buf.Len()) > lenOfLength {
		//logger.Debug("There is data in the buffer, extracting")
		lenBytes := buf.Bytes()[:lenOfLength]
		length := serialization.DecodeUint64(lenBytes)
		// logger.Debug("Length is %d", length)
		// Disconnect if we received an invalid length.
		if length < uint64(lenOfMessageID) {
			return [][]byte{}, errors.New("非法message长度")
		}

		if uint64(buf.Len())-lenOfLength < length {
			log.Info("Skipping, not enough data to read this")
			return dataArray, nil
		}

		buf.Next(int(lenOfLength)) // strip the length prefix
		data := make([]byte, length)
		_, err := buf.Read(data)
		if err != nil {
			return [][]byte{}, err
		}

		dataArray = append(dataArray, data)
	}
	return dataArray, nil
}

func decodeMessage(bytess [][]byte) ([]conn_msg.Message, error) {
	messages := make([]conn_msg.Message, 0)
	for _, bytes := range bytess {
		messageId := conn_msg.GetMessageIdFromMessageBytes(bytes)
		messageType := conn_msg.GetMessageTypeByMessageId(messageId)
		messageInterface, err := serialization.DecodeMessage(messageType, bytes, conn_msg.LenOfMessageID)
		if err != nil {
			return nil, err
		}
		message := messageInterface.(conn_msg.Message)
		messages = append(messages, message)
	}
	return messages, nil
}
