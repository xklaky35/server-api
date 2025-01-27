package main

import (
	"log"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/joho/godotenv"
	"github.com/xklaky35/wpFileReader"
)

var config filereader.Config

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config, err = filereader.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	
	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// call at midnight
		if time.Now().In(loc).Hour() == 23 && time.Now().In(loc).Minute() == 59 {
			err := callEndpoint()
			if err != nil {
				log.Print(err)
			}
		}
	}
}


func callEndpoint() error {
	req, err := http.NewRequest(http.MethodPost,"http://api:8080/wP/DailyCycle", nil)
	req.SetBasicAuth(os.Getenv("welcomePageUser"), os.Getenv("welcomePageSecret"))

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
