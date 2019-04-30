package conn_msg

import (
	"chat_group/src/config"
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

func (msg *StringMessage) ServerHandleMessage(conn *connect.Connection) error {
	t := reflect.TypeOf(msg).Elem()
	messageId := MessageTypeIdMap[t]
	stringMessageBytes, err := serialization.EncodeMessage(msg, messageId[:])
	if err != nil {
		return err
	}
	//广播消息
	sendToOtherChanObject := connect.NewSendToOtherChanObject(stringMessageBytes, conn)
	connect.GetConnectionPoolInstant().PrepareSendToOther(sendToOtherChanObject)
	if config.GetInstance().SaveChatRecord {
		datebase.DO(func() {
			//记录消息
			chatRecordDO := datebase.NewChatRecordDO(msg.Token, conn.RemoteAddress, msg.Message, time.Now())
			insertSuccess := chatRecordDO.Insert()
			if !insertSuccess {
				log.Error("insert chat record fail")
			}
		})
	}
	return nil
}

func (msg *StringMessage) ClientHandleMessage(conn *connect.Connection) error {
	return nil
}

func NewStringMessage(token, message string) StringMessage {
	return StringMessage{
		Token:   token,
		Content: MessageContent{MessageType: "STRING"},
		Message: message,
	}
}
