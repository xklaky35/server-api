package main

import "github.com/xklaky35/welcomePageAPI/middleware"

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xklaky35/wpFileReader"
)


const TIME_FORMAT string = time.RFC3339


var config filereader.Config

func main() {
	initData()
	
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	protectedRoutes := r.Group("/wP")

	// Setup auth
	protectedRoutes.Use(middleware.LoadValidUsers())
	protectedRoutes.Use(middleware.AuthMiddleware())

	{
		protectedRoutes.GET("/GetData", getData)
		protectedRoutes.POST("/UpdateGauge", update)
		protectedRoutes.POST("/AddGauge", addGauge)
		protectedRoutes.POST("/RemoveGauge", removeGauge)
		protectedRoutes.POST("/DailyCycle", dailyCycle)
	}
	
	r.Run()
}


func initData() (bool, error) {
	err := godotenv.Load()
	if err != nil {
		return false, err
	}

	config, err = filereader.LoadConfig()
	if err != nil {
		return false, err
	}

	f, err := os.OpenFile("./logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)

	return true, nil
}


func dailyCycle(c *gin.Context){
	data, err := filereader.LoadData()
	if err != nil {
		log.Print()
	}

	for i, e := range data.Gauges {
		if isToday(e.LastIncrease) {
			continue
		}
		data.Gauges[i].Value -= config.DecreaseStep
		if data.Gauges[i].Value < 0 {
			data.Gauges[i].Value = 0
		}
	}
	filereader.WriteData(&data)
}

func removeGauge(c *gin.Context){
	data, err := filereader.LoadData()
	if err != nil {
		fmt.Println(err)
		return
	}

	name := c.PostForm("name")

	if i, exists := findGauge(data.Gauges, name); exists == true{
		copy(data.Gauges[i:], data.Gauges[i+1:])
		data.Gauges[len(data.Gauges)-1] = filereader.Gauge{}
		data.Gauges = data.Gauges[:len(data.Gauges)-1]
	} else {
		c.AbortWithStatus(404)		
	}
	filereader.WriteData(&data)
}

func addGauge(c *gin.Context){
	data, err := filereader.LoadData()
	if err != nil {
		fmt.Println(err)
		return
	}
	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		fmt.Println(err)
		return
	}

	name := c.PostForm("name")

	if _, exists := findGauge(data.Gauges, name); exists == true{
		c.AbortWithStatus(404)		
	} else {
		data.Gauges = append(data.Gauges, filereader.Gauge{
			Name: name,
			Value: 0,
			LastIncrease: time.Now().In(loc).Format(TIME_FORMAT),
		})
	}

	filereader.WriteData(&data)
}

func getData(c *gin.Context){
	data, err := filereader.LoadData()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, &data)
}

func update(c *gin.Context){
	data, err := filereader.LoadData()
	if err != nil {
		log.Println(err)
		return
	}

	name := c.Query("name")

	// search for the gauge and increase it if found
	if i, exists := findGauge(data.Gauges, name); exists == true{
		err := increase(&data.Gauges[i])
		if err != nil {
			c.AbortWithStatus(401)
		}
	} else {
		c.AbortWithStatus(404)		
	}
	filereader.WriteData(&data)
}

func findGauge(g []filereader.Gauge, name string) (int,bool){
	for i, e := range g {
		if e.Name == name {
			return i, true
		}
	}
	return 0, false
}

func increase(g *filereader.Gauge) error {
	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		return err
	}


	if isToday(g.LastIncrease){
		return errors.New("Forbidden")
	}

	g.LastIncrease = time.Now().In(loc).Format(TIME_FORMAT)

	if g.Value == config.MaxValue{
		return nil
	}
	g.Value += config.IncreaseStep 
	return nil
}

func isToday(date string) bool{
	t, err := time.Parse(TIME_FORMAT,date)		
	if err != nil {
		log.Print(err)	
	}
	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Print(err)
	}

	if t.Day() != time.Now().In(loc).Day(){
		return false
	}
	return true
}

