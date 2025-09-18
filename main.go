package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{}
	registerCommands()

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		userCommand := input[0]
		args := input[1:]
		if len(args) > 0 {
			config.Args = args
		}
		command, ok := cmdRegistry[userCommand]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(&config)
		if err != nil {
			fmt.Println(fmt.Errorf("error from command: %w", err))
		}

		// Clean args after command
		config.Args = []string{}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(fmt.Errorf("error from scanner: %w", err))
	}
}
