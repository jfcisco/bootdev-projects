package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/jfcisco/pokedexcli/internal/pokeapi"
)

// Represents current state of REPL
type config struct {
	Next     string
	Previous string
	Args     []string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandHelp(c *config) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:`)

	printOrder := []string{"help", "map", "mapb", "explore", "catch", "inspect", "exit"}

	for _, key := range printOrder {
		cmd := cmdRegistry[key]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(c *config) error {
	onLastPage := c.Next == "" && c.Previous != ""
	if onLastPage {
		fmt.Println("you're on the last page")
		return nil
	}

	res, err := pokeapi.FetchLocationAreas(c.Next)
	if err != nil {
		return err
	}

	for _, area := range res.Results {
		fmt.Println(area.Name)
	}

	c.Next = res.Next
	c.Previous = res.Previous
	return nil
}

func commandMapb(c *config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := pokeapi.FetchLocationAreas(c.Previous)
	if err != nil {
		return err
	}

	for _, area := range res.Results {
		fmt.Println(area.Name)
	}

	c.Next = res.Next
	c.Previous = res.Previous
	return nil
}

func commandExplore(c *config) error {
	if len(c.Args) == 0 {
		return fmt.Errorf("error in commandExplore: please specify an area to explore")
	}

	area := c.Args[0]
	fmt.Printf("Exploring %s...\n", c.Args[0])

	data, err := pokeapi.ExploreArea(area)

	if err != nil {
		return err
	}

	if len(data.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, enc := range data.PokemonEncounters {
		fmt.Printf("- %v\n", enc.Pokemon.Name)
	}
	return nil
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandCatch(c *config) error {
	if len(c.Args) == 0 {
		return fmt.Errorf("error in commandCatch: please specify a Pokemon to catch")
	}

	name := c.Args[0]
	creature, err := pokeapi.FetchPokemon(name)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	// Lower base EXP > higher chance
	roll := rand.Intn(100)
	dc := min(creature.BaseExperience/4, 80)

	if roll < dc {
		fmt.Printf("%s escaped!\n", creature.Name)
	} else {
		fmt.Printf("%s was caught!\n", creature.Name)
		caught[creature.Name] = *creature
		fmt.Println("You may now inspect it with the inspect command")
	}
	return nil
}

func commandInspect(c *config) error {
	if len(c.Args) == 0 {
		return fmt.Errorf("error in commandInspect: please specify a Pokemon to inspect")
	}

	name := c.Args[0]
	creature, ok := caught[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	PrintDetails(&creature)
	return nil
}

func commandPokedex(c *config) error {
	fmt.Println("Your Pokedex:")
	if len(caught) == 0 {
		fmt.Println("No pokemon registered yet. Try catching one!")
		return nil
	}

	for _, p := range caught {
		fmt.Printf("\t- %s\n", p.Name)
	}
	return nil
}

var cmdRegistry map[string]cliCommand

func registerCommands() {
	cmdRegistry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas in the world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas in the world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explores the named area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to capture the given pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays information for a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows a list of captured pokemon",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
