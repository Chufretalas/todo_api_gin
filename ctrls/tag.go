package ctrls

import (
	"errors"
	"fmt"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//TODO: only return data of that user, give an access denied (403) if someone tries to access someones elses stuff

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

// TODO: only accept arbitrary user_id if user is looged in as admin
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
		fmt.Printf("result.Error: %v\n", result.Error.Error())
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

// TODO: error if the current user id does not match the user_id of the tag
func GetTagById(c *gin.Context) {

	tagId := c.Param("tag_id")

	var tag models.Tag

	result := db.DB.First(&tag, "id = ?", tagId)

	if result.Error != nil {
		fmt.Println("aqui", result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	c.JSON(200, gin.H{"tag": tag})
}
