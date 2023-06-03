package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type weatherForecast struct {
	City        string
	WeatherList []weatherData
}

type weatherData struct {
	Date        string
	Temperature float64
	Humidity    int
	Description string
}

func getWeatherForecastByCity(city string) (*weatherForecast, error) {
	baseURL := "https://api.openweathermap.org/data/2.5/forecast"
	apiURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("q", city)
	params.Add("appid", os.Getenv("WAPI"))
	params.Add("units", "metric")

	apiURL.RawQuery = params.Encode()

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherResp struct {
		City struct {
			Name string `json:"name"`
		} `json:"city"`
		List []struct {
			Date string `json:"dt_txt"`
			Main struct {
				Temperature float64 `json:"temp"`
				Humidity    int     `json:"humidity"`
			} `json:"main"`
			Weather []struct {
				Description string `json:"description"`
			} `json:"weather"`
		} `json:"list"`
	}

	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return nil, err
	}

	forecast := &weatherForecast{
		City:        weatherResp.City.Name,
		WeatherList: make([]weatherData, 0),
	}

	currentDate := time.Now().Format("2006-01-02")
	for _, data := range weatherResp.List {
		if strings.HasPrefix(data.Date, currentDate) {
			weather := weatherData{
				Date:        data.Date,
				Temperature: data.Main.Temperature,
				Humidity:    data.Main.Humidity,
				Description: data.Weather[0].Description,
			}
			forecast.WeatherList = append(forecast.WeatherList, weather)
		}
	}
	return forecast, nil
}

func getWeatherForecastByLocation(latitude, longitude float64) (*weatherForecast, error) {
	city, err := getCityName(latitude, longitude)
	if err != nil {
		return nil, err
	}

	return getWeatherForecastByCity(city)
}
