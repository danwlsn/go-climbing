package api

import (
	"fmt"
	"net/http"
	"log"
	"io"
	"time"
)

type Api struct {
	BaseUrl string
}

func (api *Api) Get(query string) *http.Response{

	url := fmt.Sprintf("%s%s", api.BaseUrl, query)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal("GET", err)
	}

	return resp

}

func (api *Api) GetWeatherAsJson(lat float32, long float32, start time.Time, end time.Time) []byte{
	start_fmt := start.Format("2006-01-02")
	end_fmt := end.Format("2006-01-02")
	query := fmt.Sprintf("latitude=%f&longitude=%f&current_weather=true&start_date=%s&end_date=%s&hourly=precipitation&elevation=505", lat, long, start_fmt, end_fmt)
	resp := api.Get(query)

	return readJsonFromResp(resp)

}

func readJsonFromResp(resp *http.Response) []byte {
	read_body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("read json", err)
	}

	return read_body
}
