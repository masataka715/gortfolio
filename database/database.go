package database

import (
	"gortfolio/config"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func Open() *gorm.DB {
	db, err := gorm.Open(config.Config.SQLDriver, config.Config.DbName)
	if err != nil {
		log.Println(err)
	}
	return db
}

func Migrate(i interface{}) {
	db := Open()
	db.AutoMigrate(i)
	defer db.Close()
}
