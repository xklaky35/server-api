package main

import (
	"fmt"
	"os"

	"github.com/xklaky35/welcomePageAPI/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xklaky35/welcome-page-api"
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

	protectedRoutes := r.Group("/wP")

	// Setup auth
	protectedRoutes.Use(middleware.LoadValidUsers())
	protectedRoutes.Use(middleware.AuthMiddleware())

	err = welcomepageapi.Init(protectedRoutes)
	if err != nil{
		fmt.Println(err)
		return
	}

	r.Run(os.Getenv("PORT"))
}
