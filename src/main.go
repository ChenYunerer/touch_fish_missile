package main

import (
	"chat_group/src/datebase"
	"chat_group/src/log"
	"os"
	"os/signal"
)

func main() {
	db := datebase.InitDB()
	defer db.Close()
	StartServer()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Interrupt")
}
