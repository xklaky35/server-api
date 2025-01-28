package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)


func HeaderSetup() gin.HandlerFunc {
	return func(c *gin.Context) {
		//CORS
		c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Authorization")
        c.Header("Access-Control-Allow-Methods", "POST, GET")

		c.Header("Cache-Control", "no-store")
		c.Header("Content-Type", "application/json")
		c.Header("X-Content-Type-Options", "nosniff")

		if c.Request.Method == "OPTIONS" {
			log.Print(c.Request)
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := rate.NewLimiter(1,2)

		if !limit.Allow() {
			log.Print(c.Request)
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
