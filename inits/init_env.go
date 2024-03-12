package inits

import (
	"log"

	"github.com/joho/godotenv"
)

func InitEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.\n", err.Error())
	}

}
