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
		command, ok := cmdRegistry[userCommand]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(&config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error from command:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error from scanner:", err)
	}
}
