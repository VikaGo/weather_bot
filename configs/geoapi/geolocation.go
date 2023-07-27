package geoapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func GetCityName(latitude, longitude float64) (string, error) {

	baseURL, err := url.Parse("https://api.opencagedata.com/geocode/v1/json")
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("q", fmt.Sprintf("%.4f %.4f", latitude, longitude))
	params.Add("key", os.Getenv("APIKEY"))

	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var locationResp struct {
		Results []struct {
			Components struct {
				City string `json:"city"`
			} `json:"components"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return "", err
	}

	if len(locationResp.Results) > 0 {
		return locationResp.Results[0].Components.City, nil
	}

	return "", fmt.Errorf("city not found")
}
