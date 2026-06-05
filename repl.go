package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"pokedex/internal/pokecache"
	"strings"
)

var ballThresholds = map[string]int{
	"pokeball":  50,
	"greatball": 75,
	"ultraball": 90,
}

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
		"clear": {
			name:        "clear",
			description: "Clears all pokemon",
			callback:    commandClear,
		},
		"rename": {
			name:        "rename",
			description: "Renames a pokemon",
			callback:    commandRename,
		},
		"release": {
			name:        "release",
			description: "Releases a pokemon",
			callback:    commandRelease,
		},
		"wander": {
			name:        "wander",
			description: "Wanders are a pokemon",
			callback:    commandWander,
		},
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	err := savePokedex(cfg)
	if err != nil {
		return err
	}
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
	fs := flag.NewFlagSet("catch", flag.ContinueOnError)
	ball := fs.String("ball", "pokeball", "type of ball to use")
	fs.Parse(args[1:])

	threshold, ok := ballThresholds[*ball]
	if !ok {
		fmt.Printf("unknown ball type: %s\n", *ball)
		return nil
	}

	fmt.Printf("Throwing a %s at %s...\n", *ball, name)
	pokemon, err := fetchPokemon(name, &cfg.Cache)
	if err != nil {
		return err
	}
	if rand.Intn(pokemon.BaseExperience) < threshold {
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
	displayName := pokemon.Name
	if pokemon.Nickname != "" {
		displayName = fmt.Sprintf("%s (%s)", pokemon.Nickname, pokemon.Name)
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", displayName, pokemon.Height, pokemon.Weight)
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
func commandClear(cfg *config, args []string) error {
	err := clearPokedex(cfg)
	if err != nil {
		return err
	}
	fmt.Println("Pokedex cleared!")
	return nil

}

func commandRename(cfg *config, args []string) error {
	if len(args) < 2 {
		fmt.Println("usage: rename <pokemon> <nickname>")
		return nil
	}
	name := args[0]
	pokemon, ok := cfg.Pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	pokemon.Nickname = args[1]
	cfg.Pokedex[name] = pokemon
	fmt.Printf("%s has been renamed to %s\n", name, pokemon.Nickname)
	return nil
}

func commandRelease(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please provide a pokemon name or nickname")
		return nil
	}
	query := args[0]
	keyToDelete := ""

	for key, p := range cfg.Pokedex {
		if key == query || p.Nickname == query {
			keyToDelete = key
			break
		}
	}
	if keyToDelete != "" {
		fmt.Printf("you have not caught this pokemon")
		return nil
	}
	delete(cfg.Pokedex, keyToDelete)
	fmt.Printf("%s has been released!\n", query)
	return nil
}
func commandWander(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please provide a location name")
		return nil
	}
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]
	data, err := fetchExplore(url, &cfg.Cache)
	if err != nil {
		return err
	}
	encounters := data.PokemonEncounters
	if len(encounters) == 0 {
		fmt.Println("no pokemon found in this area")
		return nil
	}
	random := encounters[rand.Intn(len(encounters))]
	name := random.Pokemon.Name
	fmt.Printf("A wild %s appeared!\n", name)
	wild, err := fetchPokemon(name, &cfg.Cache)
	if err != nil {
		return err
	}

	var input string
	fmt.Print("Battle? (y/n): ")
	fmt.Scan(&input)

	wonBattle := false
	if input == "y" {
		if len(cfg.Pokedex) == 0 {
			fmt.Println("you have no pokemon to battle with!")
		} else {
			fmt.Println("Your pokemon:")
			for k := range cfg.Pokedex {
				fmt.Println(" -", k)
			}
			fmt.Print("Which pokemon?: ")
			fmt.Scan(&input)
			yours, ok := cfg.Pokedex[input]
			if !ok {
				fmt.Println("you don't have that pokemon")
				return nil
			}
			wonBattle = runBattle(yours, wild)
		}
	}

	if wonBattle {
		fmt.Print("Catch it? (y/n): ")
		fmt.Scan(&input)
		if input == "y" {
			return commandCatch(cfg, []string{name})
		}
	} else {
		fmt.Print("Catch it anyway? (y/n): ")
		fmt.Scan(&input)
		if input == "y" {
			return commandCatch(cfg, []string{name})
		}
	}
	return nil
}
