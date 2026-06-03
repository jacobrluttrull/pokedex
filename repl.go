package main

import (
	"fmt"
	"math/rand"
	"os"
	"pokedex/internal/pokecache"
	"strings"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	return strings.Fields(lowered)
}

type config struct {
	Next     *string
	Previous *string
	Cache    pokecache.Cache
	Pokedex  map[string]Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
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
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Lists all pokemon in location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon by name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays pokemon information",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays pokemon by name",
			callback:    commandPokedex,
		},
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: pokedexcli help <command>")
	fmt.Println("")
	for name, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil

}
func commandMap(cfg *config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.Next != nil {
		url = *cfg.Next
	}
	data, err := fetchLocationAreas(url, &cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Next = data.Next
	cfg.Previous = data.Previous
	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	data, err := fetchLocationAreas(*cfg.Previous, &cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Next = data.Next
	cfg.Previous = data.Previous
	return nil
}
func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please provide a location name")
		return nil
	}
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]
	data, err := fetchExplore(url, &cfg.Cache)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, e := range data.PokemonEncounters {
		fmt.Println(" -", e.Pokemon.Name)
	}
	return nil
}
func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please provide a pokemon name")
		return nil
	}
	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	pokemon, err := fetchPokemon(name, &cfg.Cache)
	if err != nil {
		return err
	}
	if rand.Intn(pokemon.BaseExperience) < 50 {
		fmt.Printf("%s was caught!\n", name)
		cfg.Pokedex[name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}
func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please provide a pokemon name")
		return nil
	}
	pokemon, ok := cfg.Pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}
func commandPokedex(cfg *config, args []string) error {
	fmt.Println("Your Pokedex:")
	for name := range cfg.Pokedex {
		fmt.Println(" -", name)
	}
	return nil
}
