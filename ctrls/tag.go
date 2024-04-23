package ctrls

import (
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
)

func GetAllTags(c *gin.Context) {
	user, _ := c.Get("user")

	var tags []models.Tag

	err := db.DB.Model(&user).Association("Tags").Find(&tags)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"tags": tags})
}

func CreateTag(c *gin.Context) {

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var body struct {
		Name string `json:"name" binding:"required"`
	}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "json body must contain non-empty 'name' fields"})
		return
	}

	tag := models.Tag{Name: body.Name}

	err = db.DB.Model(&user).Association("Tags").Append(&tag)

	if err != nil {
		fmt.Println("aqui", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"tag": tag})
}
