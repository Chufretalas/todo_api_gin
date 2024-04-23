package main

import (
	"log"
	"os"
	"todo_api_gin/admin_ctrls"
	"todo_api_gin/ctrls"
	"todo_api_gin/inits"
	"todo_api_gin/middleware"

	"github.com/gin-gonic/gin"
)

var adminUser string
var adminPass string

func init() {
	inits.InitEnv()
	inits.InitDB()

	adminUser = os.Getenv("ADMIN_USER")
	adminPass = os.Getenv("ADMIN_PASS")

	if adminUser == "" || adminPass == "" {
		log.Fatal("please set up the ADMIN_USER and ADMIN_PASS environment variables")
	}
}

func main() {
	r := gin.Default()

	r.POST("/signup", ctrls.Signup)
	r.POST("/login", ctrls.Login)
	r.GET("/validate", middleware.RequireAuth, ctrls.Validate)

	// the endpoints in admin have access to everything from all users
	admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		adminUser: adminPass,
	}))

	admin.GET("/api/users", admin_ctrls.GetAllUsers)
	admin.GET("/api/users/:id", admin_ctrls.GetUserById)
	admin.PUT("/api/users/:id", admin_ctrls.UpdateUserById)
	admin.PATCH("/api/users/:id", admin_ctrls.UpdateUserById)
	admin.DELETE("/api/users/:id", admin_ctrls.DeleteUserById)

	admin.GET("/api/tags", admin_ctrls.GetAllTags)
	admin.GET("/api/tags/:tag_id", admin_ctrls.GetTagById)
	admin.POST("/api/tags", admin_ctrls.CreateTag)
	admin.PUT("/api/tags/:tag_id", admin_ctrls.UpdateTagById)
	admin.PATCH("/api/tags/:tag_id", admin_ctrls.UpdateTagById)
	admin.DELETE("/api/tags/:tag_id", admin_ctrls.DeleteTagById)

	// the endpoints in logged have access only to things of that one user
	logged := r.Group("/api", middleware.RequireAuth)

	logged.GET("/tags", ctrls.GetAllTags)
	logged.POST("/tags", ctrls.CreateTag)

	logged.GET("/todos", ctrls.GetAllTODOs)
	logged.POST("/todos", ctrls.CreateTODO)

	r.Run()
}
