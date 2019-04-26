package conn_msg

import (
	"chat_group/src/connect"
	"chat_group/src/serialization"
	"reflect"
)

type StringMessage struct {
	Content MessageContent
	Message string
}

func (msg *StringMessage) HandleMessage(conn *connect.Connection) error {
	t := reflect.TypeOf(msg).Elem()
	messageId := MessageTypeIdMap[t]
	stringMessageBytes, err := serialization.EncodeMessage(msg, messageId[:])
	if err != nil {
		return err
	}
	connect.GetConnectionPoolInstant().SendToOthers(*conn, stringMessageBytes)
	return nil
}

func NewStringMessage(message string) StringMessage {
	return StringMessage{
		Content: MessageContent{MessageType: "STRING"},
		Message: message,
	}
}
