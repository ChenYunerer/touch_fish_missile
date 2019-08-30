package conn_msg

import (
	"touch_fish_missile/src/connect"
	"touch_fish_missile/src/util"
)

type NotificationMessage struct {
	NotificationMessage string
	Content             MessageContent
}

func (msg *NotificationMessage) ServerHandleMessage(conn *connect.Connection) error {
	return nil
}

func (msg *NotificationMessage) ClientHandleMessage(conn *connect.Connection) error {
	//print msg into cmd line
	util.PrintSysNotifyToCmd(msg.NotificationMessage)
	return nil
}

func NewNotificationMessage(notificationMessage string) NotificationMessage {
	return NotificationMessage{
		NotificationMessage: notificationMessage,
		Content:             MessageContent{MessageType: "NOTI"},
	}
}
