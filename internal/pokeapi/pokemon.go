package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Pokemon struct {
	NamedApiResource
	BaseExperience int `json:"base_experience"`
	Types          []struct {
		Slot int              `json:"slot"`
		Type NamedApiResource `json:"type"`
	} `json:"types"`
}

func FetchPokemon(name string) (*Pokemon, error) {
	if name == "" {
		return nil, fmt.Errorf("please specify a pokemon")
	}

	fullUrl := pokemonBaseUrl + name

	if data, ok := cache.Get(fullUrl); ok {
		pokemon := &Pokemon{}
		if err := json.Unmarshal(data, pokemon); err != nil {
			return nil, fmt.Errorf("error in FetchPokemon: %w", err)
		}
		return pokemon, nil
	}

	res, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("error in FetchPokemon: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("error in FetchPokemon: cannot find %s", name)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in FetchPokemon: %w", err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error in FetchPokemon: %w", err)
	}

	cache.Add(fullUrl, data)

	pokemon := &Pokemon{}
	if err := json.Unmarshal(data, pokemon); err != nil {
		return nil, fmt.Errorf("error in FetchPokemon: %w", err)
	}

	return pokemon, nil
}
