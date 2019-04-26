package serialization

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func DecodeMessage(messageType reflect.Type, messageBytes []byte, lenOfMessageId uint32) (interface{}, error) {
	value := reflect.New(messageType).Interface()
	err := decode(messageBytes[lenOfMessageId:], value)
	if err != nil {
		return nil, err
	}
	return value, nil
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
