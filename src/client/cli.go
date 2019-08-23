package client

import (
	"chat_group/src/conn_msg"
	"chat_group/src/connect"
	"chat_group/src/log"
	"chat_group/src/serialization"
	"fmt"
	"reflect"
)

var cmdInputStr string

func listenCmd(conn *connect.Connection) {
	for {
		num, err := fmt.Scanln(&cmdInputStr)
		if num == 0 || err != nil {
			log.Info("listenCmd err num: ", string(num), " err: ", err)
			continue
		}
		//if strings.HasPrefix(cmdInputStr, ":") {
		//	switch cmdInputStr {
		//	case ":":
		//
		//	}
		//}
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
