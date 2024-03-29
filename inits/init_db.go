package inits

import (
	"log"
	"todo_api_gin/db"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	var err error
	db.DB, err = gorm.Open(sqlite.Open("todos_data.db"), &gorm.Config{TranslateError: true, Logger: logger.Default})
	if err != nil {
		log.Fatal(err)
	}
}
