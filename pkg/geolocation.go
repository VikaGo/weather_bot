package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getCityName(latitude, longitude float64) (string, error) {
	CageGeocodingAPIKey := os.Getenv("CageGeocodingAPIKey")
	apiURL := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%.4f+%.4f&key=%s", latitude, longitude, CageGeocodingAPIKey)

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
