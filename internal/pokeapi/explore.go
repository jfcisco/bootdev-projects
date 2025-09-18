package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type ExploreAreaResponse struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func ExploreArea(area string) (*ExploreAreaResponse, error) {
	if area == "" {
		return nil, fmt.Errorf("error in ExploreArea: empty area argument")
	}

	fullUrl := locationAreaBaseUrl + area

	result := &ExploreAreaResponse{}
	data, ok := cache.Get(fullUrl)
	if ok {
		if err := json.Unmarshal(data, result); err != nil {
			return nil, fmt.Errorf("error in ExploreArea: %w", err)
		}
		return result, nil
	}

	res, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("error in ExploreArea: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		// Area is not found
		emptyRes := &ExploreAreaResponse{PokemonEncounters: []PokemonEncounter{}}
		return emptyRes, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in ExploreArea: unsuccessful response %s", res.Status)
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error in ExploreArea: %w", err)
	}

	if err := json.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("error in ExploreArea: %w", err)
	}
	return result, nil
}
