package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		if err := validate(c); err == false {
			log.Print(c.Request)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func LoadValidUsers() gin.HandlerFunc{

	return gin.BasicAuth(gin.Accounts{
		os.Getenv("apiUser"):os.Getenv("apiSecret"),
	})
}

func validate(c *gin.Context) bool {

	_, exists := c.Get(gin.AuthUserKey)
	if exists == false {
		return exists 
	}

	return true
}
