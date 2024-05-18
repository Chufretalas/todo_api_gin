package ctrls

import (
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
)

func GetAllTODOs(c *gin.Context) {
	user, _ := c.Get("user")

	var todos []models.TODO

	err := db.DB.Model(&user).Preload("Tags").Association("TODOs").Find(&todos)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"todos": todos})
}

func GetTodoById(c *gin.Context) {

	id := c.Param("todo_id")

	user, _ := c.Get("user")

	var todo models.TODO

	err := db.DB.Model(&user).Preload("Tags").Association("TODOs").Find(&todo, id)

	if todo.ID == 0 && id != "0" {
		c.AbortWithStatusJSON(400, gin.H{"error": "This TODO does not exist or does not belong to this user"})
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"todo": todo})
}

func CreateTODO(c *gin.Context) {

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var body struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "json body must contain non-empty 'title' field"})
		return
	}

	var tags []models.Tag

	fmt.Println(tags)

	err = db.DB.Model(&user).Association("Tags").Find(&tags, body.Tags)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	todo := models.TODO{Title: body.Title, Description: body.Description}

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
