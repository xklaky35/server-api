package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xklaky35/welcomePageAPI/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil{
		fmt.Println(err)
		return
	}
	
	r := gin.Default()
	r.Use(middleware.RateLimiter())
	r.Use(middleware.HeaderSetup())

	// Setup auth
	r.Use(middleware.LoadValidUsers())
	r.Use(middleware.AuthMiddleware())

	// Redirection of welcomePageAPI URLs
	r.GET("/wP/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://localhost:3001/%s", ctx.Param("endpoint"))
		ctx.Redirect(http.StatusPermanentRedirect, redirectURL)
	} )
	r.POST("/wP/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://localhost:3001/%s", ctx.Param("endpoint"))
		ctx.Redirect(http.StatusPermanentRedirect, redirectURL)
	} )

	// Redirection of speed URLs
	r.GET("/speed/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://localhost:3002/%s", ctx.Param("endpoint"))
		ctx.Redirect(http.StatusPermanentRedirect, redirectURL)
	} )
	r.POST("/speed/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://localhost:3002/%s", ctx.Param("endpoint"))
		ctx.Redirect(http.StatusPermanentRedirect, redirectURL)
	} )

	// Redirection of speed URLs
	r.GET("/p2g/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://localhost:3003/%s", ctx.Param("endpoint"))
		ctx.Redirect(http.StatusPermanentRedirect, redirectURL)
	})
	r.POST("/p2g/:endpoint", func(ctx *gin.Context) {
		redirectURL := fmt.Sprintf("http://localhost:3003/%s", ctx.Param("endpoint"))
		ctx.Redirect(http.StatusPermanentRedirect, redirectURL)
	})
	
	r.Run(os.Getenv("PORT"))
}
