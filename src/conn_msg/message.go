package conn_msg

import (
	"chat_group/src/connect"
	"reflect"
)

type MessageContent struct {
	MessageType string
}

type Message interface {
	handleMessage(conn connect.Connection)
}

var MessageContentMap map[string]reflect.Type

func init() {
	MessageContentMap = make(map[string]reflect.Type)
	MessageContentMap["PING"] = reflect.TypeOf(PingMessage{})
	MessageContentMap["STRING"] = reflect.TypeOf(StringMessage{})
}
