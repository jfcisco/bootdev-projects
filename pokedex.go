package main

import (
	"fmt"

	"github.com/jfcisco/pokedexcli/internal/pokeapi"
)

func PrintDetails(p *pokeapi.Pokemon) {
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %v\n", p.Height)
	fmt.Printf("Weight: %v\n", p.Weight)

	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("\t- %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("\t- %s\n", t.Type.Name)
	}
}
