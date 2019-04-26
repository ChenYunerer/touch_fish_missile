package serialization

import (
	"bytes"
	"chat_group/src/conn_msg"
	"encoding/gob"
	"reflect"
)

func EncodeMessage(message conn_msg.Message) ([]byte, error) {
	t := reflect.TypeOf(message).Elem()
	messageId := conn_msg.MessageTypeIdMap[t]
	messageContentBytes, err := encode(message)
	if err != nil {
		return nil, err
	}
	messageBytes := make([]byte, 0)
	messageBytes = append(messageBytes, messageId[:]...)
	messageBytes = append(messageBytes, messageContentBytes...)
	return messageBytes, nil
}

func encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
