package config

import (
	"os"
	"time"
)

type Config struct {
	Network        string        //网络类型
	ListenIp       string        //监听地址
	ServerIp       string        //服务端地址
	Port           string        //监听端口
	WriteTimeout   time.Duration //写入超时时间
	ReadTimeout    time.Duration //读取超时时间
	RetryTimes     uint32        //重试次数
	PingDuration   time.Duration //心跳间隔
	SaveChatRecord bool          //是否保存聊天记录
	LogLevel       string        //日志打印等级
}

var config *Config

func initConfig() {
	config = &Config{
		Network:        "tcp",
		ListenIp:       "0.0.0.0",
		ServerIp:       "chat.yuner.fun",
		Port:           "8888",
		WriteTimeout:   time.Duration(15) * time.Second,
		ReadTimeout:    time.Duration(15) * time.Second,
		RetryTimes:     3,
		PingDuration:   time.Duration(5) * time.Second,
		SaveChatRecord: false,
		LogLevel:       "INFO",
	}
	//ListenIp Port 支持环境变量配置
	ip := os.Getenv("ip")
	if ip != "" {
		config.ListenIp = ip
	}
	port := os.Getenv("port")
	if port != "" {
		config.Port = port
	}
}

func GetInstance() *Config {
	if config == nil {
		initConfig()
	}
	return config
}

func (c *Config) GetListenAddress() string {
	return c.ListenIp + ":" + c.Port
}

func (c *Config) GetServerAddress() string {
	return c.ServerIp + ":" + c.Port
}
