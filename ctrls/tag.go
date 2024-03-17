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
