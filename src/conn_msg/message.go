package conn_msg

import (
	"chat_group/src/connect"
	log "github.com/sirupsen/logrus"
	"reflect"
)

const LenOfMessageID = 4

type MessageId [LenOfMessageID]byte

type MessageContent struct {
	MessageType string
}

type Message interface {
	HandleMessage(conn connect.Connection) error
}

var MessageMap map[string]interface{}
var MessageIdTypeMap map[MessageId]reflect.Type
var MessageTypeIdMap map[reflect.Type]MessageId

func init() {
	MessageMap = make(map[string]interface{})
	MessageMap["PING"] = PingMessage{}
	MessageMap["STRI"] = StringMessage{}

	MessageIdTypeMap = make(map[MessageId]reflect.Type)
	MessageTypeIdMap = make(map[reflect.Type]MessageId)
	for messageIdStr, value := range MessageMap {
		t := reflect.TypeOf(value)
		messageId := messageIdFromString(messageIdStr)
		MessageIdTypeMap[MessageId(messageId)] = t
		MessageTypeIdMap[t] = messageId
	}
}

func messageIdFromString(messageIdStr string) MessageId {
	messageId := MessageId{}
	if len(messageIdStr) > LenOfMessageID {
		log.Panic("message register err: message too long, ", messageIdStr)
	}
	for i, c := range messageIdStr {
		messageId[i] = byte(c)
	}
	//不够位数进行填充
	for i := len(messageIdStr); i < LenOfMessageID; i++ {
		messageId[i] = 0x00
	}
	return messageId
}
