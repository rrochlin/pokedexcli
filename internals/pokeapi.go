package internals

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Cache    *Cache
	Next     *string
	Previous *string
}

type LocationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetPokeLocations(conf *Config) (LocationArea, error) {
	var cleanedUrl string
	if conf.Next == nil {
		cleanedUrl = "https://pokeapi.co/api/v2/location-area/"
	} else {
		cleanedUrl = *conf.Next
	}
	var locationSearch LocationArea
	body, ok := conf.Cache.Get(cleanedUrl)
	if !ok {

		res, err := http.Get(cleanedUrl)
		if err != nil {
			return LocationArea{}, fmt.Errorf("could not get location data: %v", err)
		}

		defer res.Body.Close()
		if res.StatusCode > 299 {
			return LocationArea{}, fmt.Errorf("Server did not respond with OK, StatusCode:%v", res.StatusCode)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return LocationArea{}, fmt.Errorf("Could not read response with io: %v", err)
		}
	}

	if err := json.Unmarshal(body, &locationSearch); err != nil {
		return LocationArea{}, fmt.Errorf("Unmarshalling response failed: %v", err)
	}
	return locationSearch, nil
}
