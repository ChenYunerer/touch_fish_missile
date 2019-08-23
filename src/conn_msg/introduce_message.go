package conn_msg

import (
	"chat_group/src/connect"
)

type IntroduceMessage struct {
	Token   string
	Content MessageContent
	Group   string
}

func (msg *IntroduceMessage) ServerHandleMessage(conn *connect.Connection) error {
	conn.ConnectionUserInfo.Token = msg.Token
	conn.ConnectionUserInfo.Group = msg.Group
	return nil
}

func (msg *IntroduceMessage) ClientHandleMessage(conn *connect.Connection) error {
	return nil
}

func NewIntroduceMessage(token, group string) IntroduceMessage {
	return IntroduceMessage{
		Token:   token,
		Group:   group,
		Content: MessageContent{MessageType: "INTR"},
	}
}
