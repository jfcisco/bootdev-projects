package pokeapi

import (
	"time"

	"github.com/jfcisco/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(5 * time.Minute)

const locationAreaBaseUrl string = "https://pokeapi.co/api/v2/location-area/"