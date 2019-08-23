package client

import (
	"time"
	"touch_fish_missile/src/config"
	"touch_fish_missile/src/connect"
	"touch_fish_missile/src/log"
)

func writeLoop(conn *connect.Connection, token string, quit chan struct{}) {
	conf := config.GetInstance()
	//test: send test message per second
	t := time.NewTicker(time.Duration(1) * time.Second)
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
		case <-t.C:
			//stringMessage := conn_msg.NewStringMessage(token, "golang大法好")
			//t := reflect.TypeOf(stringMessage)
			//messageId := conn_msg.MessageTypeIdMap[t]
			//bytes, err := serialization.EncodeMessage(&stringMessage, messageId[:])
			//if err != nil {
			//	log.Error(err)
			//	continue
			//}
			//conn.SendMessageChan <- bytes
		case <-quit:
			return
		}
	}
}
