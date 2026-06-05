package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedex/internal/pokecache"
)

type Pokemon struct {
	Name           string `json:"name"`
	Nickname       string `json:"nickname,omitempty"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func fetchPokemon(name string, cache *pokecache.Cache) (Pokemon, error) {
	if body, ok := cache.Get(name); ok {
		var data Pokemon
		err := json.Unmarshal(body, &data)
		if err != nil {
			return Pokemon{}, err
		}
		return data, nil
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	resp, err := http.Get(url)

	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return Pokemon{}, fmt.Errorf("pokemon not found: %s", name)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}
	cache.Add(name, body)
	var data Pokemon
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Pokemon{}, err

	}
	return data, nil
}
