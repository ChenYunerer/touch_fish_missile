package config

import "time"

type Config struct {
	Network      string
	Ip           string
	Port         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	RetryTimes   uint32
	PingDuration time.Duration
}

var config *Config

func initConfig() {
	config = &Config{
		Network:      "tcp",
		Ip:           "127.0.0.1",
		Port:         "8888",
		WriteTimeout: time.Duration(5) * time.Second,
		ReadTimeout:  time.Duration(5) * time.Second,
		RetryTimes:   3,
		PingDuration: time.Duration(2) * time.Second,
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
