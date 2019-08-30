package conn_msg

import (
	"reflect"
	"touch_fish_missile/src/connect"
	"touch_fish_missile/src/log"
	"touch_fish_missile/src/serialization"
	"touch_fish_missile/src/util"
)

type IntroduceMessage struct {
	Token   string
	Group   string
	Message string
	Content MessageContent
}

func (msg *IntroduceMessage) ServerHandleMessage(conn *connect.Connection) error {
	conn.Token = msg.Token
	conn.Group = msg.Group
	returnMsg := "user: " + conn.Token + " has entered the " + conn.Group + " group"
	returnMsg = returnMsg + "\n"
	returnMsg = returnMsg + "群组用户: "
	groupConns := connect.GetConnectionPoolInstant().SearchConnectionsByGroup(msg.Group)
	for _, groupConn := range groupConns {
		returnMsg = returnMsg + groupConn.Token + " "
	}
	introduceMessage := NewIntroduceMessage("", "", returnMsg)
	t := reflect.TypeOf(introduceMessage)
	messageId := MessageTypeIdMap[t]
	bytes, err := serialization.EncodeMessage(&introduceMessage, messageId[:])
	if err != nil {
		log.Error(err)
	}
	conn.SendMessageChan <- bytes
	return nil
}

func (msg *IntroduceMessage) ClientHandleMessage(conn *connect.Connection) error {
	util.PrintSysNotifyToCmd(msg.Message)
	return nil
}

func NewIntroduceMessage(token, group, message string) IntroduceMessage {
	return IntroduceMessage{
		Token:   token,
		Group:   group,
		Message: message,
		Content: MessageContent{MessageType: "INTR"},
	}
}
