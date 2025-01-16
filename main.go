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
	Gauges gauges `json:"gauges"`
}

type gauges struct {
	MaxValue int `json:"max_value"`
	MinValue int `json:"min_value"`
	IncreaseStep int `json:"increase_step"`
	DecreaseStep int `json:"decrease_step"`
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
		protectedRoutes.GET("/GetConfig", getConfig)
		protectedRoutes.POST("/UpdateGauge", update)
		protectedRoutes.POST("/AddGauge", addGauge)
		protectedRoutes.POST("/RemoveGauge", removeGauge)
	}
	
	r.Run()
}

func removeGauge(c *gin.Context){
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	name := c.PostForm("name")

	if i, exists := findGauge(config.Gauges.Gauges, name); exists == true{
		copy(config.Gauges.Gauges[i:], config.Gauges.Gauges[i+1:])
		config.Gauges.Gauges[len(config.Gauges.Gauges)-1] = gauge{}
		config.Gauges.Gauges = config.Gauges.Gauges[:len(config.Gauges.Gauges)-1]
	} else {
		c.AbortWithStatus(404)		
	}
	writeConfig(&config)
}

func addGauge(c *gin.Context){
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	name := c.PostForm("name")

	if _, exists := findGauge(config.Gauges.Gauges, name); exists == true{
		c.AbortWithStatus(404)		
	} else {
		config.Gauges.Gauges = append(config.Gauges.Gauges, gauge{
			Name: name,
			Value: 0,
			LastIncrease: time.Now().Local().Format(TIME_FORMAT),
		})
	}

	writeConfig(&config)

}

func getConfig(c *gin.Context){

	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, &config)
}

func update(c *gin.Context){
	config, err := loadConfig()
	if err != nil {
		log.Println(err)
		return
	}

	name := c.Query("name")

	// search for the gauge and increase it if found
	if i, exists := findGauge(config.Gauges.Gauges, name); exists == true{
		err := increase(&config.Gauges.Gauges[i], config)
		if err != nil {
			c.AbortWithStatus(401)
		}
	} else {
		c.AbortWithStatus(404)		
	}
	writeConfig(&config)
}



func findGauge(g []gauge, name string) (int,bool){
	for i, e := range g {
		if e.Name == name {
			return i, true
		}
	}
	return 0, false
}

func increase(g *gauge, c Config) error {
	if isToday(g.LastIncrease){
		return errors.New("Forbidden")
	}

	g.LastIncrease = time.Now().Local().Format(TIME_FORMAT)

	if g.Value == c.Gauges.MaxValue{
		return nil
	}
	g.Value += c.Gauges.IncreaseStep 
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

func writeConfig(c *Config){

	d, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
		return
	}
	
	err = os.WriteFile(os.Getenv("wP_CONFIG"), d, 766)
	if err != nil{
		log.Println(err)
	}
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

