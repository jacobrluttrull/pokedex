package main

import (
	"testing"
	"time"

	"pokedex/internal/pokecache"
)

func newTestCfg() *config {
	return &config{
		Cache:   pokecache.NewCache(5 * time.Minute),
		Pokedex: map[string]Pokemon{},
	}
}

// commandHelp

func TestCommandHelp(t *testing.T) {
	cfg := newTestCfg()
	if err := commandHelp(cfg, []string{}); err != nil {
		t.Errorf("commandHelp returned unexpected error: %v", err)
	}
}

// commandMapb

func TestCommandMapbNoPrevious(t *testing.T) {
	cfg := newTestCfg()
	if err := commandMapb(cfg, []string{}); err != nil {
		t.Errorf("commandMapb with no previous returned error: %v", err)
	}
}

func TestCommandMapbWithPrevious(t *testing.T) {
	cfg := newTestCfg()
	next := "https://pokeapi.co/api/v2/location-area/?offset=20"
	prev := "https://pokeapi.co/api/v2/location-area/?offset=0"
	cfg.Previous = &prev
	body := []byte(`{"next":"` + next + `","previous":null,"results":[{"name":"test-area","url":"http://example.com"}]}`)
	cfg.Cache.Add(prev, body)

	if err := commandMapb(cfg, []string{}); err != nil {
		t.Errorf("commandMapb with previous returned error: %v", err)
	}
	if cfg.Next == nil || *cfg.Next != next {
		t.Errorf("expected cfg.Next=%q, got %v", next, cfg.Next)
	}
}

// commandMap

func TestCommandMapFirstPage(t *testing.T) {
	cfg := newTestCfg()
	next := "https://pokeapi.co/api/v2/location-area/?offset=20"
	base := "https://pokeapi.co/api/v2/location-area/"
	body := []byte(`{"next":"` + next + `","previous":null,"results":[{"name":"test-area","url":"http://example.com"}]}`)
	cfg.Cache.Add(base, body)

	if err := commandMap(cfg, []string{}); err != nil {
		t.Errorf("commandMap returned error: %v", err)
	}
	if cfg.Next == nil || *cfg.Next != next {
		t.Errorf("expected cfg.Next=%q, got %v", next, cfg.Next)
	}
}

// commandExplore

func TestCommandExploreNoArgs(t *testing.T) {
	cfg := newTestCfg()
	if err := commandExplore(cfg, []string{}); err != nil {
		t.Errorf("commandExplore with no args returned error: %v", err)
	}
}

func TestCommandExploreWithCache(t *testing.T) {
	cfg := newTestCfg()
	location := "pallet-town-area"
	url := "https://pokeapi.co/api/v2/location-area/" + location
	body := []byte(`{"pokemon_encounters":[{"pokemon":{"name":"rattata"}},{"pokemon":{"name":"pidgey"}}]}`)
	cfg.Cache.Add(url, body)

	if err := commandExplore(cfg, []string{location}); err != nil {
		t.Errorf("commandExplore returned error: %v", err)
	}
}

// commandCatch

func TestCommandCatchNoArgs(t *testing.T) {
	cfg := newTestCfg()
	if err := commandCatch(cfg, []string{}); err != nil {
		t.Errorf("commandCatch with no args returned error: %v", err)
	}
}

func TestCommandCatchGuaranteedCatch(t *testing.T) {
	cfg := newTestCfg()
	// BaseExperience=1 means rand.Intn(1)==0, always < 50 → always caught
	body := []byte(`{"name":"magikarp","base_experience":1,"height":9,"weight":100,"stats":[],"types":[]}`)
	cfg.Cache.Add("magikarp", body)

	if err := commandCatch(cfg, []string{"magikarp"}); err != nil {
		t.Errorf("commandCatch returned error: %v", err)
	}
	if _, ok := cfg.Pokedex["magikarp"]; !ok {
		t.Errorf("expected magikarp to be in Pokedex after guaranteed catch")
	}
}

// commandInspect

func TestCommandInspectNoArgs(t *testing.T) {
	cfg := newTestCfg()
	if err := commandInspect(cfg, []string{}); err != nil {
		t.Errorf("commandInspect with no args returned error: %v", err)
	}
}

func TestCommandInspectNotCaught(t *testing.T) {
	cfg := newTestCfg()
	if err := commandInspect(cfg, []string{"mewtwo"}); err != nil {
		t.Errorf("commandInspect for uncaught pokemon returned error: %v", err)
	}
}

func TestCommandInspectCaught(t *testing.T) {
	cfg := newTestCfg()
	cfg.Pokedex["pidgey"] = Pokemon{
		Name:           "pidgey",
		BaseExperience: 50,
		Height:         3,
		Weight:         18,
	}
	if err := commandInspect(cfg, []string{"pidgey"}); err != nil {
		t.Errorf("commandInspect for caught pokemon returned error: %v", err)
	}
}

// commandPokedex

func TestCommandPokedexEmpty(t *testing.T) {
	cfg := newTestCfg()
	if err := commandPokedex(cfg, []string{}); err != nil {
		t.Errorf("commandPokedex with empty pokedex returned error: %v", err)
	}
}

func TestCommandPokedexWithPokemon(t *testing.T) {
	cfg := newTestCfg()
	cfg.Pokedex["pidgey"] = Pokemon{
		Name:           "pidgey",
		BaseExperience: 50,
		Height:         3,
		Weight:         18,
	}
	cfg.Pokedex["caterpie"] = Pokemon{
		Name:           "caterpie",
		BaseExperience: 39,
		Height:         3,
		Weight:         29,
	}
	if err := commandPokedex(cfg, []string{}); err != nil {
		t.Errorf("commandPokedex with pokemon returned error: %v", err)
	}
}
