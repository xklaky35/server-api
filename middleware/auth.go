package middleware

import (
	"errors"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		if err := validate(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func LoadValidUsers() gin.HandlerFunc{

	return gin.BasicAuth(gin.Accounts{
		os.Getenv("welcomePageUser"):os.Getenv("welcomePageSecret"),
	})
}

func validate(c *gin.Context) error{

	_, err := c.Get(gin.AuthUserKey)
	if err == false {
		return errors.New("Wrong auth")
	}

	return nil
}
