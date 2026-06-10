package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
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

// helpers

func seedPokemon(cfg *config, name string) {
	cfg.Pokedex[name] = Pokemon{Name: name, BaseExperience: 50, Height: 3, Weight: 18}
}

func seedCache(cfg *config, name string) {
	body, _ := json.Marshal(Pokemon{Name: name, BaseExperience: 1})
	cfg.Cache.Add(name, body)
}

// commandHelp

func TestCommandHelp(t *testing.T) {
	cfg := newTestCfg()
	if err := commandHelp(cfg, []string{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// commandMap / commandMapb

func TestCommandMap(t *testing.T) {
	next := "https://pokeapi.co/api/v2/location-area/?offset=20"
	base := "https://pokeapi.co/api/v2/location-area/"
	body := []byte(`{"next":"` + next + `","previous":null,"results":[{"name":"test-area","url":"http://example.com"}]}`)

	tests := []struct {
		name    string
		seedURL string
		nextSet string
	}{
		{"first page", base, next},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newTestCfg()
			cfg.Cache.Add(tt.seedURL, body)
			if err := commandMap(cfg, []string{}); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if cfg.Next == nil || *cfg.Next != tt.nextSet {
				t.Errorf("expected Next=%q, got %v", tt.nextSet, cfg.Next)
			}
		})
	}
}

func TestCommandMapb(t *testing.T) {
	next := "https://pokeapi.co/api/v2/location-area/?offset=20"
	prev := "https://pokeapi.co/api/v2/location-area/?offset=0"
	body := []byte(`{"next":"` + next + `","previous":null,"results":[{"name":"test-area","url":"http://example.com"}]}`)

	t.Run("no previous", func(t *testing.T) {
		cfg := newTestCfg()
		if err := commandMapb(cfg, []string{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("with previous", func(t *testing.T) {
		cfg := newTestCfg()
		cfg.Previous = &prev
		cfg.Cache.Add(prev, body)
		if err := commandMapb(cfg, []string{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if cfg.Next == nil || *cfg.Next != next {
			t.Errorf("expected Next=%q, got %v", next, cfg.Next)
		}
	})
}

// commandExplore

func TestCommandExplore(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		cfg := newTestCfg()
		if err := commandExplore(cfg, []string{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("cached location", func(t *testing.T) {
		cfg := newTestCfg()
		location := "pallet-town-area"
		url := "https://pokeapi.co/api/v2/location-area/" + location
		cfg.Cache.Add(url, []byte(`{"pokemon_encounters":[{"pokemon":{"name":"rattata"}},{"pokemon":{"name":"pidgey"}}]}`))
		if err := commandExplore(cfg, []string{location}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("empty encounters", func(t *testing.T) {
		cfg := newTestCfg()
		location := "empty-area"
		url := "https://pokeapi.co/api/v2/location-area/" + location
		cfg.Cache.Add(url, []byte(`{"pokemon_encounters":[]}`))
		if err := commandExplore(cfg, []string{location}); err != nil {
			t.Errorf("unexpected error on empty encounters: %v", err)
		}
	})
}

// commandCatch

func TestCommandCatch(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		cfg := newTestCfg()
		if err := commandCatch(cfg, []string{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("guaranteed catch default ball", func(t *testing.T) {
		cfg := newTestCfg()
		seedCache(cfg, "magikarp")
		if err := commandCatch(cfg, []string{"magikarp"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if _, ok := cfg.Pokedex["magikarp"]; !ok {
			t.Error("expected magikarp in Pokedex after guaranteed catch")
		}
	})
	t.Run("unknown pokemon returns error", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping network call in short mode")
		}
		cfg := newTestCfg()
		err := commandCatch(cfg, []string{"notarealpokemon12345"})
		if err == nil {
			t.Error("expected error for unknown pokemon, got nil")
		}
	})
}

func TestCommandCatchBallTypes(t *testing.T) {
	balls := []struct {
		name string
		flag string
	}{
		{"pokeball", "--ball=pokeball"},
		{"greatball", "--ball=greatball"},
		{"ultraball", "--ball=ultraball"},
	}
	for _, tt := range balls {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newTestCfg()
			seedCache(cfg, "magikarp")
			if err := commandCatch(cfg, []string{"magikarp", tt.flag}); err != nil {
				t.Errorf("unexpected error with %s: %v", tt.name, err)
			}
			// BaseExperience=1 guarantees catch for all ball types
			if _, ok := cfg.Pokedex["magikarp"]; !ok {
				t.Errorf("expected magikarp caught with %s", tt.name)
			}
		})
	}
	t.Run("unknown ball type", func(t *testing.T) {
		cfg := newTestCfg()
		seedCache(cfg, "magikarp")
		if err := commandCatch(cfg, []string{"magikarp", "--ball=masterball"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if _, ok := cfg.Pokedex["magikarp"]; ok {
			t.Error("expected no catch with unknown ball type")
		}
	})
}

// commandInspect

func TestCommandInspect(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		seed    bool
		wantErr bool
	}{
		{"no args", []string{}, false, false},
		{"not caught", []string{"mewtwo"}, false, false},
		{"caught pokemon", []string{"pidgey"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newTestCfg()
			if tt.seed {
				seedPokemon(cfg, "pidgey")
			}
			err := commandInspect(cfg, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandInspect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// commandPokedex

func TestCommandPokedex(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		cfg := newTestCfg()
		if err := commandPokedex(cfg, []string{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("with multiple pokemon", func(t *testing.T) {
		cfg := newTestCfg()
		seedPokemon(cfg, "pidgey")
		seedPokemon(cfg, "caterpie")
		if err := commandPokedex(cfg, []string{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(cfg.Pokedex) != 2 {
			t.Errorf("expected 2 pokemon, got %d", len(cfg.Pokedex))
		}
	})
}

// commandClear

func TestCommandClear(t *testing.T) {
	t.Run("clears in-memory pokedex", func(t *testing.T) {
		cfg := newTestCfg()
		seedPokemon(cfg, "bulbasaur")
		t.Setenv("HOME", t.TempDir())
		if err := commandClear(cfg, []string{}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(cfg.Pokedex) != 0 {
			t.Errorf("expected empty Pokedex after clear, got %d entries", len(cfg.Pokedex))
		}
	})
	t.Run("no-op when no save file", func(t *testing.T) {
		cfg := newTestCfg()
		t.Setenv("HOME", t.TempDir())
		if err := commandClear(cfg, []string{}); err != nil {
			t.Errorf("unexpected error when no file exists: %v", err)
		}
	})
}

// save / load roundtrip

func TestSaveLoadRoundtrip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg := newTestCfg()
	seedPokemon(cfg, "charmander")

	if err := savePokedex(cfg); err != nil {
		t.Fatalf("savePokedex failed: %v", err)
	}

	cfg2 := newTestCfg()
	if err := loadPokedex(cfg2); err != nil {
		t.Fatalf("loadPokedex failed: %v", err)
	}
	if _, ok := cfg2.Pokedex["charmander"]; !ok {
		t.Error("expected charmander to load from disk")
	}
}

func TestLoadNoFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg := newTestCfg()
	if err := loadPokedex(cfg); err != nil {
		t.Errorf("loadPokedex with no file should be a no-op, got: %v", err)
	}
}

// battle system

func testMon(t *testing.T, jsonStr string) Pokemon {
	t.Helper()
	var p Pokemon
	if err := json.Unmarshal([]byte(jsonStr), &p); err != nil {
		t.Fatalf("bad test pokemon json: %v", err)
	}
	return p
}

func TestTypeEffect(t *testing.T) {
	cases := []struct {
		moveType string
		defTypes []string
		want     float64
	}{
		{"fire", []string{"grass"}, 2},
		{"fire", []string{"water"}, 0.5},
		{"electric", []string{"ground"}, 0},
		{"water", []string{"fire", "rock"}, 4},
		{"normal", []string{"normal"}, 1},
		{"ice", []string{"dragon", "flying"}, 4},
	}
	for _, c := range cases {
		if got := typeEffect(c.moveType, c.defTypes); got != c.want {
			t.Errorf("typeEffect(%s vs %v) = %v, want %v", c.moveType, c.defTypes, got, c.want)
		}
	}
}

func TestDamageRoll(t *testing.T) {
	charmander := testMon(t, `{"name":"charmander","types":[{"type":{"name":"fire"}}],"stats":[{"base_stat":52,"stat":{"name":"attack"}},{"base_stat":43,"stat":{"name":"defense"}}]}`)
	bulbasaur := testMon(t, `{"name":"bulbasaur","types":[{"type":{"name":"grass"}},{"type":{"name":"poison"}}],"stats":[{"base_stat":49,"stat":{"name":"attack"}},{"base_stat":49,"stat":{"name":"defense"}}]}`)
	ember := Move{Name: "ember", Power: 40, Accuracy: 100}
	ember.Type.Name = "fire"

	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 50; i++ {
		dmg, msgs := damageRoll(charmander, bulbasaur, ember, rng)
		// base ≈ (22*40*52/49)/50+2 ≈ 20.7, STAB 1.5, eff 2 → ~62 before
		// variance/crit; never a miss at 100 accuracy, never below 1
		if dmg < 40 || dmg > 160 {
			t.Fatalf("damage %d outside sane range", dmg)
		}
		joined := strings.Join(msgs, " ")
		if !strings.Contains(joined, "super effective") {
			t.Fatalf("expected super effective message, got %q", joined)
		}
		if strings.Contains(joined, "missed") {
			t.Fatalf("100-accuracy move missed")
		}
	}
}

func TestDamageRollImmune(t *testing.T) {
	pikachu := testMon(t, `{"name":"pikachu","types":[{"type":{"name":"electric"}}],"stats":[{"base_stat":55,"stat":{"name":"attack"}}]}`)
	diglett := testMon(t, `{"name":"diglett","types":[{"type":{"name":"ground"}}],"stats":[{"base_stat":25,"stat":{"name":"defense"}}]}`)
	shock := Move{Name: "thunder-shock", Power: 40, Accuracy: 100}
	shock.Type.Name = "electric"

	rng := rand.New(rand.NewSource(1))
	dmg, msgs := damageRoll(pikachu, diglett, shock, rng)
	if dmg != 0 {
		t.Errorf("expected 0 damage vs immune type, got %d", dmg)
	}
	if !strings.Contains(strings.Join(msgs, " "), "doesn't affect") {
		t.Errorf("expected immunity message, got %v", msgs)
	}
}

func TestPickMovesFallsBackToStruggle(t *testing.T) {
	cfg := newTestCfg()
	noMoves := testMon(t, `{"name":"ditto"}`)
	moves := pickMoves(noMoves, cfg.Cache)
	if len(moves) != 1 || moves[0].Name != "struggle" {
		t.Errorf("expected struggle fallback, got %v", moves)
	}
}

func TestClearDeletesFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	cfg := newTestCfg()
	seedPokemon(cfg, "squirtle")

	if err := savePokedex(cfg); err != nil {
		t.Fatalf("savePokedex failed: %v", err)
	}

	if err := clearPokedex(cfg); err != nil {
		t.Fatalf("clearPokedex failed: %v", err)
	}

	// file should be gone
	path, err := pokedexPath()
	if err != nil {
		t.Fatalf("pokedexPath failed: %v", err)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Error("expected save file to be deleted after clear")
	}
}
