package main

import (
	"os"
	"path/filepath"
)

func pokedexPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pokedex"), nil
}

func clearPokedex(cfg *config) error {
	cfg.Pokedex = map[string]Pokemon{}
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
