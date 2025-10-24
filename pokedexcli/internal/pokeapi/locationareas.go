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

func FetchLocationAreas(url string) (*LocationAreaResponse, error) {
	if url == "" {
		url = locationAreaBaseUrl + "?offset=0&limit=20"
	}

	result := &LocationAreaResponse{}
	data, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(data, result); err != nil {
			return nil, err
		}
		return result, nil
	}

	// Cache miss, query from PokeApi
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in FetchLocationAreas: unsuccessful response %s", res.Status)
	}

	cache.Add(url, data)
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return result, nil
}
