package ctrls

import (
	"errors"
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func GetTODOById(c *gin.Context) {

	id := c.Param("todo_id")

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var todo models.TODO

	findResult := db.DB.Preload("Tags").First(&todo, `id = ? AND user_id = ?`, id, user.ID)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no todo belonging to this user was found for id = %v", id))
		return
	}

	if findResult.Error != nil {
		fmt.Println(findResult.Error.Error())
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

func DeleteTODOById(c *gin.Context) {
	id := c.Param("todo_id")

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var todo models.TODO

	findResult := db.DB.Preload("Tags").First(&todo, `id = ? AND user_id = ?`, id, user.ID)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no todo belonging to this user was found for id = %v", id))
		return
	}

	deleteResult := db.DB.Delete(&todo)

	if deleteResult.Error != nil {
		fmt.Printf("result.Error: %v\n", deleteResult.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"deleted TODO": todo})
}
