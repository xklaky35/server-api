package main

import (
	"github.com/xklaky35/welcomePageAPI/middleware"

	"github.com/gin-gonic/gin"
	"github.com/xklaky35/wpFileReader"
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


	{
		protectedRoutes.GET("/GetData", getData)
		protectedRoutes.POST("/UpdateGauge", update) //param
		protectedRoutes.POST("/AddGauge", addGauge) //body
		protectedRoutes.POST("/RemoveGauge", removeGauge) //body
		protectedRoutes.POST("/DailyCycle", dailyCycle)
	}
	
	r.Run()
}


