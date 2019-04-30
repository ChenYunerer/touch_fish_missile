package config

import "time"

type Config struct {
	Network        string        //网络类型
	Ip             string        //监听地址
	Port           string        //监听端口
	WriteTimeout   time.Duration //写入超时时间
	ReadTimeout    time.Duration //读取超时时间
	RetryTimes     uint32        //重试次数
	PingDuration   time.Duration //心跳间隔
	SaveChatRecord bool          //是否保存聊天记录
}

var config *Config

func initConfig() {
	config = &Config{
		Network:        "tcp",
		Ip:             "127.0.0.1",
		Port:           "8888",
		WriteTimeout:   time.Duration(5) * time.Second,
		ReadTimeout:    time.Duration(5) * time.Second,
		RetryTimes:     3,
		PingDuration:   time.Duration(2) * time.Second,
		SaveChatRecord: false,
	}
}

func GetInstance() *Config {
	if config == nil {
		initConfig()
	}
	return config
}

func (c *Config) GetAddress() string {
	return c.Ip + ":" + c.Port
}
