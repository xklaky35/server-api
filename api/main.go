package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xklaky35/welcomePageAPI/middleware"
)


const TIME_FORMAT string = time.RFC3339

type Config struct {
	MaxValue int `json:"max_value"`
	MinValue int `json:"min_value"`
	IncreaseStep int `json:"increase_step"`
	DecreaseStep int `json:"decrease_step"`
}

type Data struct {
	Gauges []gauge `json:"gauges"`
}

type gauge struct {
	Name string `json:"name"`
	Value int `json:"value"`
	LastIncrease string `json:"last_increase"`
}

func main() {
	
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error reading .env file")
	}

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


func dailyCycle(c *gin.Context){
	data, err := loadData()
	if err != nil {
		log.Fatal()
	}
	config, err := loadConfig()
	if err != nil {
		log.Fatal()
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
	writeData(&data)
}

func removeGauge(c *gin.Context){
	data, err := loadData()
	if err != nil {
		fmt.Println(err)
		return
	}

	name := c.PostForm("name")

	if i, exists := findGauge(data.Gauges, name); exists == true{
		copy(data.Gauges[i:], data.Gauges[i+1:])
		data.Gauges[len(data.Gauges)-1] = gauge{}
		data.Gauges = data.Gauges[:len(data.Gauges)-1]
	} else {
		c.AbortWithStatus(404)		
	}
	writeData(&data)
}

func addGauge(c *gin.Context){
	data, err := loadData()
	if err != nil {
		fmt.Println(err)
		return
	}

	name := c.PostForm("name")

	if _, exists := findGauge(data.Gauges, name); exists == true{
		c.AbortWithStatus(404)		
	} else {
		data.Gauges = append(data.Gauges, gauge{
			Name: name,
			Value: 0,
			LastIncrease: time.Now().Local().Format(TIME_FORMAT),
		})
	}

	writeData(&data)
}

func getData(c *gin.Context){

	data, err := loadData()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, &data)
}

func update(c *gin.Context){
	data, err := loadData()
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
	writeData(&data)
}



func findGauge(g []gauge, name string) (int,bool){
	for i, e := range g {
		if e.Name == name {
			return i, true
		}
	}
	return 0, false
}

func increase(g *gauge) error {

	config, err := loadConfig()
	if err != nil{
		log.Fatal()
	}

	if isToday(g.LastIncrease){
		return errors.New("Forbidden")
	}

	g.LastIncrease = time.Now().Local().Format(TIME_FORMAT)

	if g.Value == config.MaxValue{
		return nil
	}
	g.Value += config.IncreaseStep 
	return nil
}

func isToday(date string) bool{
	t, err := time.Parse(TIME_FORMAT,date)		
	if err != nil {
		log.Fatal(err)	
	}

	if t.Day() != time.Now().Day(){
		return false
	}
	return true
}

func writeData(d *Data){

	data, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
		return
	}
	
	err = os.WriteFile(os.Getenv("wP_DATA"), data, 766)
	if err != nil{
		log.Println(err)
	}
}

func loadData() (Data, error){
	var data Data  
	f, err := os.ReadFile(os.Getenv("wP_DATA"))
	if err != nil {
		return data, err
	}
		
	err = json.Unmarshal(f, &data)

	return data, nil
}

func loadConfig() (Config, error){
	var config Config 
	f, err := os.ReadFile(os.Getenv("wP_CONFIG"))
	if err != nil {
		return config, err
	}
		
	err = json.Unmarshal(f, &config)

	return config, nil
}

