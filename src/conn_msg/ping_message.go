package conn_msg

import (
	"reflect"
	"touch_fish_missile/src/connect"
	"touch_fish_missile/src/serialization"
)

type PingMessage struct {
	Content MessageContent
}

func (msg *PingMessage) ServerHandleMessage(conn *connect.Connection) error {
	return nil
}

func (msg *PingMessage) ClientHandleMessage(conn *connect.Connection) error {
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

func NewPingMessage() PingMessage {
	return PingMessage{
		Content: MessageContent{MessageType: "PING"},
	}
}
