package datebase

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

const NameOfDB = "chat_record.db"

var DB *gorm.DB
var DBProcessFuncChan chan func()

func InitDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", NameOfDB)
	if err != nil {
		log.Panic(err)
	}
	DB = db
	if !db.HasTable(&ChatRecordDO{}) {
		db.CreateTable(&ChatRecordDO{})
	}
	go handleDBInsertProcess()
	return db
}

func handleDBInsertProcess() {
	DBProcessFuncChan = make(chan func(), 1024)
	for {
		select {
		case dbFunc := <-DBProcessFuncChan:
			dbFunc()
		}
	}
}

func DO(dbFunc func()) {
	DBProcessFuncChan <- dbFunc
}
