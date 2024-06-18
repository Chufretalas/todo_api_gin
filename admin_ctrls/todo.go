package admin_ctrls

import (
	"errors"
	"fmt"
	"strconv"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllTODOs(c *gin.Context) {
	var todos []models.TODO

	result := db.DB.Preload("Tags").Find(&todos)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"todos": todos})
}

func GetTODOById(c *gin.Context) {

	id := c.Param("todo_id")

	var todo []models.TODO

	result := db.DB.Preload("Tags").First(&todo, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no todo with id = %v was found", id))
		return
	}

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"todo": todo})
}

func CreateTODO(c *gin.Context) {

	userId := c.Query("user_id")
	numericUserId, idParseError := strconv.ParseUint(userId, 10, 64)
	if idParseError != nil {
		fmt.Println(idParseError.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "user_id is not a valid unsigned integer"})
		return
	}

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

	tagsResult := db.DB.Where("id IN ? AND user_id = ?", body.Tags, userId).Find(&tags)

	if tagsResult.Error != nil {
		fmt.Println(tagsResult.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	todo := models.TODO{Title: body.Title, Description: body.Description, UserID: uint(numericUserId)}

	todoResult := db.DB.Create(&todo)
	if todoResult.Error != nil {
		fmt.Println(todoResult.Error.Error())
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

func UpdateTODObyId(c *gin.Context) {
	id := c.Param("todo_id")

	var todo models.TODO

	findResult := db.DB.Preload("Tags").First(&todo, `id = ?`, id)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no todo with id = %v was found", id))
		return
	}

	var body struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "body must be valid json"})
		return
	}

	var tags []models.Tag

	tagsResult := db.DB.Where("id IN ? AND user_id = ?", body.Tags, todo.UserID).Find(&tags)

	if tagsResult.Error != nil {
		fmt.Println(tagsResult.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	updateResult := db.DB.Preload("Tags").Model(&todo).Updates(models.TODO{Title: body.Title, Description: body.Description})

	if updateResult.Error != nil {
		fmt.Printf("result.Error: %v\n", updateResult.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	err = db.DB.Model(&todo).Association("Tags").Replace(&tags)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"todo": todo})
}

func DeleteTODOById(c *gin.Context) {
	id := c.Param("todo_id")

	var todo models.TODO

	findResult := db.DB.Preload("Tags").First(&todo, id)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no todo with id = %v was found", id))
		return
	}

	deleteResult := db.DB.Unscoped().Delete(&todo)

	if deleteResult.Error != nil {
		fmt.Printf("result.Error: %v\n", deleteResult.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"deleted TODO": todo})
}
