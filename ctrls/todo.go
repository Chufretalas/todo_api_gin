package ctrls

import (
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
)

// Fetch the tags of each todo
func GetAllTODOs(c *gin.Context) {
	user, _ := c.Get("user")

	var todos []models.TODO

	err := db.DB.Model(&user).Association("TODOs").Find(&todos)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"todos": todos})
}

// TODO: make this receive a body
func CreateTODO(c *gin.Context) {

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var tags []models.Tag

	err := db.DB.Model(&user).Association("Tags").Find(&tags, []string{"4", "5"})

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	todo := models.TODO{Title: "Titulo", Description: "Descrição do todo"}

	err = db.DB.Model(&user).Association("TODOs").Append(&todo)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	err = db.DB.Model(&todo).Association("Tags").Append(&tags)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"todo": todo})

}
