package main

import (
	"todo_api_gin/ctrls"
	"todo_api_gin/inits"

	"github.com/gin-gonic/gin"
)

func init() {
	inits.InitEnv()
	inits.InitDB()
}

func main() {
	r := gin.Default()
	r.GET("/api/users", ctrls.GetAllUsers)
	r.POST("/api/users", ctrls.CreateUser)
	r.Run() // listen and serve on 0.0.0.0:8080
}
