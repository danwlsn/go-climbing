package location

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func parseDate(date string) string {
	datetime, err := time.Parse("2006-01-02T15:04", date)

	if err != nil {
		log.Fatal("date parse", err)
	}

	return datetime.Format("Mon, 02-01 15:04")
}

type Hourly struct {
	Time []string
	Precipitation []float32
}

func (forecast *Hourly) formatOutput(index int) string {
	date_format := parseDate(forecast.Time[index])
	msg := fmt.Sprintf("%f at %s", forecast.Precipitation[index], date_format)
	return msg
}

type Location struct {
	Name string
	Latitude float32
	Longitude float32
	Hourly Hourly
}

func (location *Location) TotalRainfall() float32 {

	var rain float32

	for _, hourly := range location.Hourly.Precipitation{
		rain += hourly
	}

	return rain
}

func (location *Location) UpdateForecastFromJson(weather_json []byte) {
	err := json.Unmarshal(weather_json, &location)
	if err != nil {
		log.Fatal("json", err)
	}
}
func (location *Location) PrintWeather(){
	for i:=0; (i<len(location.Hourly.Time)); i++ {
		fmt.Println(location.Hourly.formatOutput(i))
	}
}

func (location *Location) WriteWeatherToBuffer(buff *bytes.Buffer){
	for i:=0; (i<len(location.Hourly.Time)); i++ {
		buff.WriteString(fmt.Sprint(location.Hourly.formatOutput(i)))
		buff.WriteString("\n")
	}
}
