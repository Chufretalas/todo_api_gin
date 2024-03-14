package ctrls

import (
	"errors"
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)

	c.JSON(200, gin.H{"users": users})
}

func CreateUser(c *gin.Context) {

	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "json body must contain non-empty 'username' and 'password' fields"})
		return
	}

	//TODO: hash the password
	user := models.User{Username: body.Username, Passhash: body.Password}

	result := db.DB.Create(&user)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("username '%v' is already taken", body.Username)})
		return
	}

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unkown error"})
		return
	}

	c.JSON(200, gin.H{"user": user})
}
