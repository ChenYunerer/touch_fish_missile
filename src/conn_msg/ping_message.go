package conn_msg

import "chat_group/src/connect"

type PingMessage struct {
	Content MessageContent
}

func (msg *PingMessage) handleMessage(conn connect.Connection) {

}

func NewPingMessage() PingMessage {
	return PingMessage{
		Content: MessageContent{MessageType: "PING"},
	}
}
