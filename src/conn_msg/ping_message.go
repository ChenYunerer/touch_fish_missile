package conn_msg

import (
	"chat_group/src/connect"
	"chat_group/src/serialization"
	"reflect"
)

type PingMessage struct {
	Content MessageContent
}

func (msg *PingMessage) ServerHandleMessage(conn *connect.Connection) error {
	pongMessage := NewPongMessage()
	t := reflect.TypeOf(pongMessage)
	messageId := MessageTypeIdMap[t]
	pongMessageBytes, err := serialization.EncodeMessage(&pongMessage, messageId[:])
	if err != nil {
		return err
	}
	conn.SendMessageChan <- pongMessageBytes
	return nil
}

func (msg *PingMessage) ClientHandleMessage(conn *connect.Connection) error {
	return nil
}

func NewPingMessage() PingMessage {
	return PingMessage{
		Content: MessageContent{MessageType: "PING"},
	}
}
