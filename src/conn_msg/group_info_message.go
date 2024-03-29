package conn_msg

import (
	"reflect"
	"touch_fish_missile/src/connect"
	"touch_fish_missile/src/log"
	"touch_fish_missile/src/serialization"
	"touch_fish_missile/src/util"
)

type GroupInfoMessage struct {
	Token   string
	Group   string
	Message string
	Content MessageContent
}

func (msg *GroupInfoMessage) ServerHandleMessage(conn *connect.Connection) error {
	conn.Token = msg.Token
	conn.Group = msg.Group
	returnMsg := "user: " + conn.Token + " has entered the " + conn.Group + " group"
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

func (msg *GroupInfoMessage) ClientHandleMessage(conn *connect.Connection) error {
	util.PrintSysNotifyToCmd(msg.Message)
	return nil
}

func NewGroupInfoMessage(token, group, message string) GroupInfoMessage {
	return GroupInfoMessage{
		Token:   token,
		Group:   group,
		Message: message,
		Content: MessageContent{MessageType: "INTR"},
	}
}
