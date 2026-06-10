package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"

	"github.com/chzyer/readline"
)

func main() {
	commands := getCommands()
	cfg := &config{
		Cache:   pokecache.NewCache(5 * time.Minute),
		Pokedex: map[string]pokeapi.Pokemon{},
	}
	err := loadPokedex(cfg)
	if err != nil {
		fmt.Printf("error loading pokedex: %v\n", err)
		return
	}
	logFile, err := os.OpenFile("repl.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening log: %v\n", err)
		return
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	rl, err := readline.New("Pokedex > ")
	if err != nil {
		fmt.Printf("error starting repl: %v\n", err)
		return
	}
	defer rl.Close()
	for {
		line, err := rl.Readline()
		if err != nil {
			return
		}
		text := cleanInput(line)
		if len(text) == 0 {
			continue
		}
		logger.Println(line)
		cmd, exists := commands[text[0]]
		if !exists {
			fmt.Println("Command not found")
			continue
		}
		err = cmd.callback(cfg, text[1:])
		if err != nil {
			fmt.Println(err)
		}
	}
}
