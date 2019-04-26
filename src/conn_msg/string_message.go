package conn_msg

import (
	"chat_group/src/connect"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type StringMessage struct {
	Content MessageContent
	Message string
}

func (msg *StringMessage) HandleMessage(conn *connect.Connection) error {
	messageJsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
	}
	connect.GetConnectionPoolInstant().SendToOthers(conn, messageJsonBytes)
	return nil
}

func NewStringMessage(message string) StringMessage {
	return StringMessage{
		Content: MessageContent{MessageType: "STRING"},
		Message: message,
	}
}
