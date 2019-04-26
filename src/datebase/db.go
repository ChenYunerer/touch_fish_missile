package datebase

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

const NameOfDB = "chat_record.db"

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", NameOfDB)
	if err != nil {
		log.Panic(err)
	}
	DB = db
	if !db.HasTable(&ChatRecordDO{}) {
		db.CreateTable(&ChatRecordDO{})
	}
	return db
}
