package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jfcisco/pokedexcli/internal/pokeapi"
)

// Represents current state of REPL
type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandHelp(c *config) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:`)

	printOrder := []string{"help", "map", "mapb", "exit"}

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

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
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
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
