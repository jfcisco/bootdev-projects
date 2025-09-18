package pokeapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func FetchLocationAreas(url string) (LocationAreaResponse, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return LocationAreaResponse{}, errors.New("fetchLocationAreas: unsuccessful response")
	}

	var result LocationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&result); err != nil {
		return LocationAreaResponse{}, err
	}
	return result, nil
}
