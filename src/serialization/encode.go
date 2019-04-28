package serialization

import (
	"bytes"
	"encoding/gob"
)

func EncodeMessage(message interface{}, messageId []byte) ([]byte, error) {
	messageContentBytes, err := encode(message)
	if err != nil {
		return nil, err
	}
	messageBytes := make([]byte, 0)
	lenBytes := EncodeUint64(uint64(len(messageId) + len(messageContentBytes)))
	messageBytes = append(messageBytes, lenBytes[:]...)
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

func EncodeUint64(v uint64) [8]byte {
	var b [8]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	return b
}

func DecodeUint64(bytes []byte) uint64 {
	return uint64(bytes[0]) | uint64(bytes[1])<<8 | uint64(bytes[2])<<16 | uint64(bytes[3])<<24 |
		uint64(bytes[4])<<32 | uint64(bytes[5])<<40 | uint64(bytes[6])<<48 | uint64(bytes[7])<<56
}
