package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Gauges gauges `json:"gauges"`
}

type gauges struct {
	MaxValue int `json:"max_value"`
	MinValue int `json:"min_value"`
	GaugeArt gauge `json:"gauge_art"`
	GaugeCoding gauge `json:"gauge_coding"`
	GaugeMusic gauge `json:"gauge_music"`
	GaugeJapanese gauge `json:"gauge_japanese"`
}

type gauge struct {
	Value int `json:"value"`
	LastIncrease string `json:"last_increase"`
}

func main() {
	

	r := gin.Default()

	r.Handle(http.MethodGet, "/wPGC", getConfig)
	r.Handle(http.MethodPost, "/wPUA", updateArt)
	r.Handle(http.MethodPost, "/wPUC", updateCoding)
	r.Handle(http.MethodPost, "/wPUM", updateMusic)
	r.Handle(http.MethodPost, "/wPUJ", updateJapanese)
	
	r.Run()
}

func getConfig(c *gin.Context){

	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, config)
}
func updateArt(c *gin.Context){
	update("art")
}
func updateCoding(c *gin.Context){
	update("coding")
}
func updateMusic(c *gin.Context){
	update("music")
}
func updateJapanese(c *gin.Context){
	update("japanese")
}


func update(gauge string){
	f, err := os.OpenFile("logs/log.log", 0, 766)
	if err != nil{
		fmt.Println("Log could not be opened")
		return
	}
	l := log.New(f, "> ", 0)


	config, err := loadConfig()
	if err != nil {
		l.Println(err)
		return
	}
	
	switch(gauge){
		case "art":{
			config.Gauges.GaugeArt.Value++
			config.Gauges.GaugeArt.LastIncrease = time.Now().Format(time.RFC1123)
		}
		case "coding":{
			config.Gauges.GaugeCoding.Value++
			config.Gauges.GaugeCoding.LastIncrease = time.Now().Format(time.RFC1123)
		}
		case "music":{
			config.Gauges.GaugeMusic.Value++
			config.Gauges.GaugeMusic.LastIncrease = time.Now().Format(time.RFC1123)
		}
		case "japanese":{
			config.Gauges.GaugeJapanese.Value++
			config.Gauges.GaugeJapanese.LastIncrease = time.Now().Format(time.RFC1123)
		}
	}



	d, err := json.Marshal(&config)
	if err != nil {
		l.Println(err)
		return
	}
	
	err = os.WriteFile("data/config.json", d, 766)
	if err != nil{
		l.Println(err)
	}
}



func loadConfig() (Config, error){
	var config Config  
	f, err := os.ReadFile("data/config.json")
	if err != nil {
		return config, err
	}
		
	err = json.Unmarshal(f, &config)

	return config, nil
}

