package ctrls

import (
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)

	c.JSON(200, gin.H{"users": users})
}

func CreateUser(c *gin.Context) {

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.BindJSON(&body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(400)
		return
	}

	//TODO: hash the password
	user := models.User{Username: body.Username, Passhash: body.Password}

	result := db.DB.Create(&user)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{"user": user})
}
