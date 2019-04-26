package conn_msg

import (
	"chat_group/src/connect"
	"chat_group/src/datebase"
	"chat_group/src/serialization"
	"github.com/prometheus/common/log"
	"reflect"
	"time"
)

type StringMessage struct {
	Token   string
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
	//广播消息
	connect.GetConnectionPoolInstant().SendToOthers(*conn, stringMessageBytes)
	datebase.DO(func() {
		//记录消息
		chatRecordDO := datebase.NewChatRecordDO(msg.Token, conn.RemoteAddress, msg.Message, time.Now())
		insertSuccess := chatRecordDO.Insert()
		if !insertSuccess {
			log.Error("insert chat record fail")
		}
	})
	return nil
}

func NewStringMessage(message string) StringMessage {
	return StringMessage{
		Token:   "614482989",
		Content: MessageContent{MessageType: "STRING"},
		Message: message,
	}
}
