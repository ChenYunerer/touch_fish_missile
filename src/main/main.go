package main

import (
	"chat_group/src/client"
	"chat_group/src/config"
	"chat_group/src/datebase"
	"chat_group/src/log"
	"chat_group/src/server"
	"flag"
	"os"
	"os/signal"
)

var serverMode bool
var token string
var logLevel string

func cmd() {
	flag.BoolVar(&serverMode, "s", false, "start with server mode or client mode; true clientMode false serverMode")
	flag.StringVar(&token, "token", "unknown", "client token")
	flag.StringVar(&logLevel, "logLevel", "error", "log level: panic fatal error warn info debug trace")
	flag.Parse()
}

func main() {
	cmd()
	conf := config.GetInstance()
	conf.LogLevel = logLevel
	log.Init()
	if conf.SaveChatRecord {
		db := datebase.InitDB()
		defer db.Close()
	}
	if serverMode {
		log.Info("Now Start With Server Mode")
		server.StartServer()
	} else {
		log.Info("Now Start With Client Mode")
		client.StartClient(token)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Interrupt")
}
