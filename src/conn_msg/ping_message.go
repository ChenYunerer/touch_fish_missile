package conn_msg

import "chat_group/src/connect"

type PingMessage struct {
	Content MessageContent
}

func (msg *PingMessage) HandleMessage(conn connect.Connection) error {
	return nil
}

func NewPingMessage() PingMessage {
	return PingMessage{
		Content: MessageContent{MessageType: "PING"},
	}
}
