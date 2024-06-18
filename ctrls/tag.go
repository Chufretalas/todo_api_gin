package ctrls

import (
	"errors"
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func GetTagById(c *gin.Context) {
	id := c.Param("tag_id")

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var tag models.Tag

	findResult := db.DB.First(&tag, `id = ? AND user_id = ?`, id, user.ID)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no tag belonging to this user was found for id = %v", id))
		return
	}

	if findResult.Error != nil {
		fmt.Println(findResult.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"tag": tag})
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
		c.AbortWithStatusJSON(400, gin.H{"error": "json body must contain non-empty 'name' field"})
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

func UpdateTagById(c *gin.Context) {
	id := c.Param("tag_id")

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var tag models.Tag

	findResult := db.DB.First(&tag, `id = ? AND user_id = ?`, id, user.ID)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no tag belonging to this user was found for id = %v", id))
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

	untyped_user, _ := c.Get("user")
	user := untyped_user.(models.User)

	var tag models.Tag

	findResult := db.DB.First(&tag, `id = ? AND user_id = ?`, id, user.ID)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, fmt.Sprintf("no tag belonging to this user was found for id = %v", id))
		return
	}

	result := db.DB.Unscoped().Delete(&tag)

	if result.Error != nil {
		fmt.Printf("result.Error: %v\n", result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"deleted tag": tag})
}
