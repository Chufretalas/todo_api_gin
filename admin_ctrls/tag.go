package admin_ctrls

import (
	"errors"
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllTags(c *gin.Context) {

	var tags []models.Tag

	result := db.DB.Find(&tags)

	if result.Error != nil {
		fmt.Printf("result.Error: %v\n", result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"tags": tags})
}

// it receives the user id as a URL parameter: ?user_id=id
func CreateTag(c *gin.Context) {

	// =============== finding the user ===============
	userId := c.Request.URL.Query().Get("user_id")

	if userId == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "user_id must be passed by url parameter"})
		return
	}

	var user models.User

	result := db.DB.First(&user, "id = ?", userId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no user was found for id = %v", userId))
		return
	}

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}
	// =============== end finding the user ===============

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

func GetTagById(c *gin.Context) {

	tagId := c.Param("tag_id")

	var tag models.Tag

	result := db.DB.First(&tag, "id = ?", tagId)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"tag": tag})
}

func UpdateTagById(c *gin.Context) {

	id := c.Param("tag_id")

	var tag models.Tag

	findResult := db.DB.First(&tag, "id = ?", id)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no tag was found for id = %v", id))
		return
	}

	var body struct {
		Name string `json:"name"`
	}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "body must be valid json"})
		return
	}

	updateResult := db.DB.Model(&tag).Updates(models.Tag{Name: body.Name})

	if updateResult.Error != nil {
		fmt.Printf("result.Error: %v\n", updateResult.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"tag": tag})
}

func DeleteTagById(c *gin.Context) {
	id := c.Param("tag_id")

	result := db.DB.Delete(&models.Tag{}, "id = ?", id)

	if result.Error != nil {
		fmt.Printf("result.Error: %v\n", result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.Status(200)
}
