package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"pokedex/internal/pokeapi"
)

func pokedexPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pokedex"), nil
}

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

func clearPokedex(cfg *config) error {
	cfg.Pokedex = map[string]pokeapi.Pokemon{}
	path, err := pokedexPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
