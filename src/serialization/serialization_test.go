package serialization

import (
	"chat_group/src/conn_msg"
	"testing"
)

func TestSerialization(t *testing.T) {
	stringMessage := conn_msg.NewStringMessage("123")
	bytes, err := conn_msg.EncodeMessage(&stringMessage)
	if err != nil {
		print(err)
		return
	}
	message, err := conn_msg.DecodeMessage(bytes)
	if err != nil {
		print(err)
		return
	}
	print(message)
}
