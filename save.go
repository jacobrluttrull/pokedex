package main

import (
	"encoding/json"
	"os"
)

func savePokedex(cfg *config) error {
	data, err := json.Marshal(cfg.Pokedex)
	if err != nil {
		return err
	}
	path, err := pokedexPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadPokedex(cfg *config) error {
	path, err := pokedexPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &cfg.Pokedex)
}
