package ctrls

import (
	"errors"
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//TODO: hash the password in create and update

func GetAllUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)

	c.JSON(200, gin.H{"users": users})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	result := db.DB.First(&user, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no user was found for id = %v", id))
		return
	}

	if result.Error != nil {
		fmt.Printf("result.Error: %v\n", result.Error.Error())
		c.AbortWithStatusJSON(500, "unknown error")
		return
	}

	c.JSON(200, gin.H{"user": user})
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

	user := models.User{Username: body.Username, Passhash: body.Password}

	result := db.DB.Create(&user)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("username '%v' is already taken", body.Username)})
		return
	}

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"user": user})
}

func UpdateUserById(c *gin.Context) {

	id := c.Param("id")

	var user models.User

	findResult := db.DB.First(&user, "id = ?", id)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no user was found for id = %v", id))
		return
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "body must be valid json"})
		return
	}

	updateResult := db.DB.Model(&user).Updates(models.User{Username: body.Username, Passhash: body.Password})

	if updateResult.Error != nil {
		fmt.Printf("result.Error: %v\n", updateResult.Error.Error())
		c.AbortWithStatusJSON(500, "unknown error")
		return
	}

	c.JSON(200, gin.H{"user": user})
}

func DeleteUserById(c *gin.Context) {

	id := c.Param("id")

	result := db.DB.Delete(&models.User{}, "id = ?", id)

	if result.Error != nil {
		fmt.Printf("result.Error: %v\n", result.Error.Error())
		c.AbortWithStatusJSON(500, "unknown error")
		return
	}

	c.Status(200)
}
