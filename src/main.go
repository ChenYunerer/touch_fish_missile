package main

import (
	"chat_group/src/datebase"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func main() {
	db := datebase.InitDB()
	defer db.Close()
	StartServer()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Interrupt")
}
