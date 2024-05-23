package middleware

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"todo_api_gin/db"
	"todo_api_gin/models"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"gorm.io/gorm"
)

// walkthrough I used for the auth stuff: https://medium.com/readytowork-org/secure-your-go-web-application-jwt-authentication-e65a5af7c049

func RequireAuth(c *gin.Context) {
	tokenHeader := c.Request.Header.Get("Authorization")
	if tokenHeader == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "incorrectly formatted authorization header"})
		return
	}

	split := strings.Split(tokenHeader, " ")
	if len(split) != 2 {
		c.AbortWithStatusJSON(401, gin.H{"error": "incorrectly formatted authorization header"})
		return
	}

	tokenString := split[1]

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKey(jwa.HS256, []byte(os.Getenv("JWT_SECRET"))))

	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "the token is not valid"})
		return
	}

	userId, _ := token.Get("user_id")

	// Find the user with token Subject
	var user models.User
	result := db.DB.First(&user, "id = ?", userId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(401, gin.H{"error": fmt.Sprintf("no user was found for id = %v", userId)})
		return
	}

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": "unknown error"})
		return
	}

	// Attach the request
	c.Set("user", user)

	//Continue
	c.Next()
}
