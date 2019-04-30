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

var startType string
var token string

var serverMode = true

func cmd() {
	flag.StringVar(&startType, "startType", "server", "start with server mode or client mode")
	flag.StringVar(&token, "token", "unknown", "client token")
	flag.Parse()
	switch startType {
	case "server":
		log.Info("Now Start With Server Mode")
		serverMode = true
		break
	case "client":
		log.Info("Now Start With Client Mode")
		log.Info("Client Token is ", token)
		serverMode = false
		break
	default:
		panic("Parameter Error")
	}
}

func main() {
	cmd()
	conf := config.GetInstance()
	if conf.SaveChatRecord {
		db := datebase.InitDB()
		defer db.Close()
	}
	if serverMode {
		server.StartServer()
	} else {
		client.StartClient()
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Interrupt")
}
