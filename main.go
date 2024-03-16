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
	r.GET("/api/users/:id", ctrls.GetUserById)
	r.POST("/api/users", ctrls.CreateUser)
	r.PUT("/api/users/:id", ctrls.UpdateUserById)
	r.PATCH("/api/users/:id", ctrls.UpdateUserById)
	r.DELETE("/api/users/:id", ctrls.DeleteUserById)

	r.GET("/api/tags", ctrls.GetAllTags)
	r.POST("/api/tags", ctrls.CreateTag)
	r.Run()
}
