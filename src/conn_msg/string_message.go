package conn_msg

import (
	"github.com/prometheus/common/log"
	"reflect"
	"time"
	"touch_fish_missile/src/config"
	"touch_fish_missile/src/connect"
	"touch_fish_missile/src/datebase"
	"touch_fish_missile/src/serialization"
	"touch_fish_missile/src/util"
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
	//broadcast message
	sendToOtherChanObject := connect.NewSendToOtherChanObject(stringMessageBytes, conn)
	connect.GetConnectionPoolInstant().PrepareSendToOther(sendToOtherChanObject)
	if config.GetInstance().SaveChatRecord {
		datebase.DO(func() {
			//record message into db
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
	//print msg into cmd line
	util.PrintMsgToCmd(time.Now().Format("2006-01-02 15:04:05"), msg.Token, ": ", msg.Message)
	return nil
}

func NewStringMessage(token, message string) StringMessage {
	return StringMessage{
		Token:   token,
		Content: MessageContent{MessageType: "STRI"},
		Message: message,
	}
}
