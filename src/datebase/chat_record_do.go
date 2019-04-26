package datebase

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ChatRecordDO struct {
	gorm.Model
	Token      string
	IpAddress  string
	Message    string
	CreateTime time.Time
}

func NewChatRecordDO(token, ipAddress, message string, createTime time.Time) *ChatRecordDO {
	return &ChatRecordDO{
		Token:      token,
		IpAddress:  ipAddress,
		Message:    message,
		CreateTime: createTime,
	}
}

func (do *ChatRecordDO) Insert() bool {
	DB.Create(do)
	return true
}
