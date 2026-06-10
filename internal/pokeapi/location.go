package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"pokedex/internal/pokecache"
)

type LocationAreaResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func FetchLocationAreas(url string, cache *pokecache.Cache) (LocationAreaResp, error) {
	if body, ok := cache.Get(url); ok {
		var data LocationAreaResp
		err := json.Unmarshal(body, &data)
		if err != nil {
			return LocationAreaResp{}, err
		}
		for _, area := range data.Results {
			fmt.Println(area.Name)
		}
		return data, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResp{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResp{}, err
	}

	cache.Add(url, body)

	var data LocationAreaResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		return LocationAreaResp{}, err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	return data, nil
}

type ExploreResp struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func FetchExplore(url string, cache *pokecache.Cache) (ExploreResp, error) {
	if body, ok := cache.Get(url); ok {
		var data ExploreResp
		return data, json.Unmarshal(body, &data)
	}
	resp, err := http.Get(url)
	if err != nil {
		return ExploreResp{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExploreResp{}, err
	}
	cache.Add(url, body)
	var data ExploreResp
	return data, json.Unmarshal(body, &data)
}
