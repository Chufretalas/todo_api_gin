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

func setupRouter(quiet bool) *gin.Engine {
	r := gin.New()

	if !quiet {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.POST("/signup", ctrls.Signup)
	r.POST("/login", ctrls.Login)
	r.GET("/validate", middleware.RequireAuth, ctrls.Validate)

	// the endpoints in admin have access to everything from all users
	admin := r.Group("/admin/api", gin.BasicAuth(gin.Accounts{
		adminUser: adminPass,
	}))

	admin.GET("/users", admin_ctrls.GetAllUsers)
	admin.GET("/users/:id", admin_ctrls.GetUserById)
	admin.PUT("/users/:id", admin_ctrls.UpdateUserById)
	admin.PATCH("/users/:id", admin_ctrls.UpdateUserById)
	admin.DELETE("/users/:id", admin_ctrls.DeleteUserById)

	admin.GET("/tags", admin_ctrls.GetAllTags)
	admin.GET("/tags/:tag_id", admin_ctrls.GetTagById)
	admin.POST("/tags", admin_ctrls.CreateTag)
	admin.PUT("/tags/:tag_id", admin_ctrls.UpdateTagById)
	admin.PATCH("/tags/:tag_id", admin_ctrls.UpdateTagById)
	admin.DELETE("/tags/:tag_id", admin_ctrls.DeleteTagById)

	// the endpoints in logged have access only to things of that one user
	logged := r.Group("/api", middleware.RequireAuth)

	logged.GET("/tags", ctrls.GetAllTags)
	logged.GET("/tags/:tag_id", ctrls.GetTagById)
	logged.POST("/tags", ctrls.CreateTag)
	logged.PUT("/tags/:tag_id", ctrls.UpdateTagById)
	logged.PATCH("/tags/:tag_id", ctrls.UpdateTagById)
	logged.DELETE("/tags/:tag_id", ctrls.DeleteTagById)

	logged.GET("/todos", ctrls.GetAllTODOs)
	logged.GET("/todos/:todo_id", ctrls.GetTODOById)
	logged.POST("/todos", ctrls.CreateTODO)
	logged.PUT("/todos/:todo_id", ctrls.UpdateTODObyId)
	logged.PATCH("/todos/:todo_id", ctrls.UpdateTODObyId)
	logged.DELETE("/todos/:todo_id", ctrls.DeleteTODOById)

	return r
}

func main() {
	r := setupRouter(false)
	r.Run()
}
