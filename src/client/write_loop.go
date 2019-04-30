package client

import (
	"chat_group/src/config"
	"chat_group/src/connect"
	"chat_group/src/log"
	"time"
)

func writeLoop(conn *connect.Connection, quit chan struct{}) {
	conf := config.GetInstance()
	for {
		if conf.WriteTimeout.Seconds() != 0 {
			deadline := time.Now().Add(conf.WriteTimeout)
			err := conn.Conn.SetWriteDeadline(deadline)
			if err != nil {
				log.Error(err)
			}
		}
		select {
		case messageBytes, ok := <-conn.SendMessageChan:
			if !ok {
				log.Info("SendMessageChan Closed")
				return
			}
			log.Info("Send Message To ", conn.RemoteAddress)
			n, err := conn.Conn.Write(messageBytes)
			if err != nil {
				log.Error(err)
				conn.AddRetryTimes()
				connRetryTimes := conn.GetRetryTimes()
				log.Info("Write Conn ", conn.RemoteAddress, " Retry Times Is ", connRetryTimes)
				if connRetryTimes >= conf.RetryTimes {
					return
				}
				continue
			}
			if n == 0 {
				log.Error("Send Data Error")
				return
			}
		case <-quit:
			return
		}
	}
}
