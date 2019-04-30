package conn_msg

import (
	"bytes"
	"chat_group/src/log"
	"chat_group/src/serialization"
	"errors"
)

func DecodeData(buf *bytes.Buffer) ([][]byte, error) {
	dataArray := [][]byte{}
	lenOfMessageID := LenOfMessageID
	lenOfLength := uint64(8)
	for uint64(buf.Len()) > lenOfLength {
		//logger.Debug("There is data in the buffer, extracting")
		lenBytes := buf.Bytes()[:lenOfLength]
		length := serialization.DecodeUint64(lenBytes)
		// logger.Debug("Length is %d", length)
		// Disconnect if we received an invalid length.
		if length < uint64(lenOfMessageID) {
			return [][]byte{}, errors.New("非法message长度")
		}

		if uint64(buf.Len())-lenOfLength < length {
			log.Info("Skipping, not enough data to read this")
			return dataArray, nil
		}

		buf.Next(int(lenOfLength)) // strip the length prefix
		data := make([]byte, length)
		_, err := buf.Read(data)
		if err != nil {
			return [][]byte{}, err
		}

		dataArray = append(dataArray, data)
	}
	return dataArray, nil
}

func DecodeMessage(bytess [][]byte) ([]Message, error) {
	messages := make([]Message, 0)
	for _, bytes := range bytess {
		messageId := GetMessageIdFromMessageBytes(bytes)
		messageType := GetMessageTypeByMessageId(messageId)
		messageInterface, err := serialization.DecodeMessage(messageType, bytes, LenOfMessageID)
		if err != nil {
			return nil, err
		}
		message := messageInterface.(Message)
		messages = append(messages, message)
	}
	return messages, nil
}
