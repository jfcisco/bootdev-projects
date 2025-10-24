package pokeapi

import (
	"time"

	"github.com/jfcisco/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(5 * time.Minute)

const locationAreaBaseUrl string = "https://pokeapi.co/api/v2/location-area/"
const pokemonBaseUrl string = "https://pokeapi.co/api/v2/pokemon/"

type NamedApiResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
