package main

import (
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	return strings.Fields(lowered)
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays this help message",
			callback:    commandHelp,
		},
	}
}
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: pokedexcli help <command>")
	fmt.Println("")
	for name, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil

}
