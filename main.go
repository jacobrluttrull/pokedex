package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal/pokecache"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// make a simple repl system for the pokedex
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	cfg := &config{
		Cache:   pokecache.NewCache(5 * time.Minute),
		Pokedex: map[string]Pokemon{},
	}
	err := loadPokedex(cfg)
	if err != nil {
		fmt.Printf("error loading pokedex: %v\n", err)
		return
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := cleanInput(scanner.Text())
		if len(text) == 0 {
			continue
		}
		cmd, exists := commands[text[0]]
		if !exists {
			fmt.Println("Command not found")
			continue
		}
		err := cmd.callback(cfg, text[1:])
		if err != nil {
			return
		}
	}
}
