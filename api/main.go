package main

import (
	"github.com/xklaky35/welcomePageAPI/middleware"

	"github.com/gin-gonic/gin"
	"github.com/xklaky35/wpFileReader"
	"github.com/xklaky35/welcome-page-api"
)

var config filereader.Config

func main() {
	
	r := gin.Default()
	r.Use(middleware.RateLimiter())
	r.Use(middleware.HeaderSetup())

	protectedRoutes := r.Group("/wP")

	// Setup auth
	protectedRoutes.Use(middleware.LoadValidUsers())
	protectedRoutes.Use(middleware.AuthMiddleware())

	welcomepageapi.Init(protectedRoutes)


	r.Run()
}


