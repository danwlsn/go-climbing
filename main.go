package main

import (
	"os"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"weather/api"
	"weather/location"
	"github.com/mailgun/mailgun-go/v4"
)

func findWeekendDates() (time.Time, time.Time ) {
	now := time.Now()
	daysTillSunday := 7 - int(now.Weekday())
	daysTillWeekend := daysTillSunday - 2
	start := now.AddDate(0,0,daysTillWeekend)
	end := now.AddDate(0,0,daysTillSunday)

	return start, end
}

type CragManager struct {
	crags []location.Location
	api api.Api
}

func BuildCragManagerFromJson(json_file string) CragManager{
	weather_api := api.Api{BaseUrl: "https://api.open-meteo.com/v1/ecmwf?"}
	cragManager := CragManager{api: weather_api}
	var crags []location.Location

	data, err := ioutil.ReadFile(json_file)
	if err != nil {
		log.Fatal("couldn't read file", err)
	}

	err = json.Unmarshal(data, &crags)
	if err != nil {
		log.Fatal("couldn't build crags", err)
	}

	cragManager.crags = crags

	return cragManager
}

func (manager *CragManager) collectWeatherInfo(start time.Time, end time.Time) {
	for i:= 0; i < len(manager.crags); i++ {
		crag := &manager.crags[i]
		weather_json := manager.api.GetWeatherAsJson(crag.Latitude, crag.Longitude, start, end)
		crag.UpdateForecastFromJson(weather_json)
	}
}

func (manager *CragManager) LowestRainfall() location.Location {
	bestCrag := manager.crags[0]

	for _, crag := range manager.crags{
		if (crag.TotalRainfall() < bestCrag.TotalRainfall()){
			bestCrag = crag
		}
	}

	return bestCrag

}


func (manager *CragManager) PrintWeather() string {
	for _, crag := range manager.crags {
		fmt.Println(fmt.Sprintf("Weather for: %s", crag.Name))
		fmt.Println(fmt.Sprintf("Total rainfall: %f", crag.TotalRainfall()))
		crag.PrintWeather()
	}

	output := fmt.Sprintf("%s is the best", manager.LowestRainfall().Name)
	return output
}

func (manager *CragManager) BuildWeatherEmail() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("%s is the best", manager.LowestRainfall().Name))
	b.WriteString("\n")

	for _, crag := range manager.crags {
		b.WriteString(fmt.Sprintf("Weather for: %s", crag.Name))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("Total rainfall: %f", crag.TotalRainfall()))
		b.WriteString("\n")
		crag.WriteWeatherToBuffer(&b)
		b.WriteString("\n")
	}


	return b.String()
}

func main() {
	fmt.Println("go climbing")
	cragManager := BuildCragManagerFromJson("./crags.json")

	start, end := findWeekendDates()

	cragManager.collectWeatherInfo(start, end)
	body := cragManager.BuildWeatherEmail()
	fmt.Println(body)

	// Create an instance of the Mailgun Client
	yourDomain := os.Getenv("MAILGUN_DOMAIN")
	privateAPIKey := os.Getenv("MAILGUN_API_KEY")
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	sender := fmt.Sprintf("sender@%s", yourDomain)
	subject := "Weather for the weekend"
	recipient := os.Getenv("MAILGUN_EMAIL")

	m := mg.NewMessage(
		sender,
		subject,
		body,
		recipient,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, id, err := mg.Send(ctx, m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
