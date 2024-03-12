package main

import (
	"todo_api_gin/inits"

	"github.com/gin-gonic/gin"
)

func init() {
	inits.InitEnv()
	inits.InitDB()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
