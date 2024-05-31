package ctrls

import (
	"errors"
	"fmt"
	"os"
	"time"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "json body must contain non-empty 'username' and 'password' fields"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Failed to hash password.",
		})
		return
	}

	user := models.User{Username: body.Username, Passhash: string(hash)}

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

	c.Status(200)
}

func Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(400, gin.H{"error": "json body must contain non-empty 'username' and 'password' fields"})
		return
	}

	// find user
	var user models.User

	result := db.DB.First(&user, "username = ?", body.Username)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("no user was found for username = %v", body.Username)})
		return
	}

	if result.Error != nil {
		fmt.Printf("result.Error: %v\n", result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}
	// end find user

	// compare hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Passhash), []byte(body.Password))

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid username or password",
		})
		return
	}
	// end compare hash

	// generate JWT
	token, err := jwt.NewBuilder().
		Claim("user_id", user.ID).
		Claim("exp", time.Now().Add(time.Hour*1).Unix()).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"error": "failed to create token",
		})
		return
	}

	tokenString, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(os.Getenv("JWT_SECRET"))))

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"error": "failed to sign token",
		})
		return
	}
	// end generate JWT

	c.JSON(200, gin.H{"token": string(tokenString)})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	// user.(models.User).Username    -->   to access specific data

	c.JSON(200, gin.H{
		"message": user,
	})
}
