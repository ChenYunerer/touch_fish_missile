package conn_msg

import (
	"chat_group/src/connect"
	"chat_group/src/serialization"
)

type PingMessage struct {
	Content MessageContent
}

func (msg *PingMessage) HandleMessage(conn *connect.Connection) error {
	pongMessage := NewPongMessage()
	pongMessageBytes, err := serialization.EncodeMessage(&pongMessage)
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
