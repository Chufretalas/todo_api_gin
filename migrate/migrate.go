package main

import (
	"todo_api_gin/db"
	"todo_api_gin/inits"
	"todo_api_gin/models"
)

func init() {
	inits.InitEnv()
	inits.InitDB()
}

func main() {
	db.DB.AutoMigrate(&models.User{}, &models.TODO{}, &models.Tag{})
}
