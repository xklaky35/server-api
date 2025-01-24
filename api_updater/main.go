package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)


func main(){
	godotenv.Load()
	for {
		// call at midnight
		if time.Now().Hour() == 0 && time.Now().Minute() == 0 {
			callEndpoint()
		}
	}
}


func callEndpoint(){
	req, err := http.NewRequest(http.MethodPost,"http://api:8080/wP/DailyCycle", nil)
	req.SetBasicAuth(os.Getenv("welcomePageUser"), os.Getenv("welcomePageSecret"))
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
