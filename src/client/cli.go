package client

import (
	"fmt"
	"reflect"
	"touch_fish_missile/src/conn_msg"
	"touch_fish_missile/src/connect"
	"touch_fish_missile/src/log"
	"touch_fish_missile/src/serialization"
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
