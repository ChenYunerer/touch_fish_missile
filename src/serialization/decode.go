package serialization

import (
	"bytes"
	"chat_group/src/conn_msg"
	"encoding/gob"
	"errors"
	"reflect"
)

func DecodeMessage(messageBytes []byte) (conn_msg.Message, error) {
	if len(messageBytes) <= conn_msg.LenOfMessageID {
		return nil, errors.New("message invalid: len of message too short")
	}
	messageIdBytes := messageBytes[:conn_msg.LenOfMessageID]
	messageId := conn_msg.MessageId{}
	for i, v := range messageIdBytes {
		messageId[i] = v
	}
	messageType := conn_msg.MessageIdTypeMap[messageId]
	value := reflect.New(messageType)
	message := value.Interface().(conn_msg.Message)
	err := decode(messageBytes[conn_msg.LenOfMessageID:], message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func decode(messageBytes []byte, data interface{}) error {
	var buf bytes.Reader
	buf.Reset(messageBytes)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
