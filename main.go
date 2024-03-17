package main

import (
	"todo_api_gin/ctrls"
	"todo_api_gin/inits"
	"todo_api_gin/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	inits.InitEnv()
	inits.InitDB()
}

func main() {
	r := gin.Default()

	r.POST("/signup", ctrls.Signup)
	r.POST("/login", ctrls.Login)
	r.GET("/validate", middleware.RequireAuth, ctrls.Validate) // here RequireAuth is a middleware that we will be creating below. It protects the route

	r.GET("/api/users", ctrls.GetAllUsers)
	r.GET("/api/users/:id", ctrls.GetUserById)
	r.PUT("/api/users/:id", ctrls.UpdateUserById)
	r.PATCH("/api/users/:id", ctrls.UpdateUserById)
	r.DELETE("/api/users/:id", ctrls.DeleteUserById)

	r.GET("/api/tags", ctrls.GetAllTags)
	r.GET("/api/tags/:tag_id", ctrls.GetTagById)
	r.POST("/api/tags", ctrls.CreateTag)
	r.Run()
}
