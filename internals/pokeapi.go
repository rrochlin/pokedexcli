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
	Pokedex  map[string]Pokemon
}

const (
	locationURL = "https://pokeapi.co/api/v2/location-area/"
	pokemonURL  = "https://pokeapi.co/api/v2/pokemon/"
)

func GetPokeLocations(conf *Config) (LocationArea, error) {
	var cleanedUrl string
	if conf.Next == nil {
		cleanedUrl = locationURL
	} else {
		cleanedUrl = *conf.Next
	}
	var locationSearch LocationArea
	if err := handleGetRequest(conf, cleanedUrl, &locationSearch); err != nil {
		return LocationArea{}, err
	}
	return locationSearch, nil
}

func GetPokeArea(conf *Config, area string) (PokeArea, error) {
	url := locationURL + area
	var pokeArea PokeArea
	if err := handleGetRequest(conf, url, &pokeArea); err != nil {
		return PokeArea{}, err
	}
	return pokeArea, nil
}

func GetPokemon(conf *Config, pokemon string) (Pokemon, error) {
	url := pokemonURL + pokemon
	var mon Pokemon
	if err := handleGetRequest(conf, url, &mon); err != nil {
		return Pokemon{}, err
	}
	return mon, nil
}

func handleGetRequest[T any](conf *Config, url string, object *T) error {
	body, ok := conf.Cache.Get(url)
	if !ok {

		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("could not get %T: %v", *new(T), err)
		}

		defer res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("Server did not respond with OK, StatusCode:%v", res.StatusCode)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Could not read response with io: %v", err)
		}
	}

	if err := json.Unmarshal(body, &object); err != nil {
		return fmt.Errorf("Unmarshalling response failed: %v", err)
	}

	return nil
}
