package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
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
		url = locationAreaBaseUrl
	}

	var data []byte
	data, ok := cache.Get(url)
	if ok {
		var result LocationAreaResponse
		if err := json.Unmarshal(data, &result); err != nil {
			return LocationAreaResponse{}, err
		}
		return result, nil
	}

	// Cache miss, query from PokeApi
	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)

	if err != nil {
		return LocationAreaResponse{}, err
	} else if res.StatusCode != http.StatusOK {
		return LocationAreaResponse{}, fmt.Errorf("error in FetchLocationAreas: unsuccessful response %s", res.Status)
	}

	cache.Add(url, data)
	var result LocationAreaResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return LocationAreaResponse{}, err
	}
	return result, nil
}
