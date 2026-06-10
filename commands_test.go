package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
)

func newTestCfg() *config {
	return &config{
		Cache:   pokecache.NewCache(5 * time.Minute),
		Pokedex: map[string]pokeapi.Pokemon{},
	}
}

// helpers

func seedPokemon(cfg *config, name string) {
	cfg.Pokedex[name] = pokeapi.Pokemon{Name: name, BaseExperience: 50, Height: 3, Weight: 18}
}

func seedCache(cfg *config, name string) {
	body, _ := json.Marshal(pokeapi.Pokemon{Name: name, BaseExperience: 1})
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

// commandRename / commandRelease

func TestCommandRename(t *testing.T) {
	t.Run("missing args", func(t *testing.T) {
		cfg := newTestCfg()
		if err := commandRename(cfg, []string{"pidgey"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("renames caught pokemon", func(t *testing.T) {
		cfg := newTestCfg()
		seedPokemon(cfg, "pidgey")
		if err := commandRename(cfg, []string{"pidgey", "skybird"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if cfg.Pokedex["pidgey"].Nickname != "skybird" {
			t.Errorf("expected nickname skybird, got %q", cfg.Pokedex["pidgey"].Nickname)
		}
	})
}

func TestCommandRelease(t *testing.T) {
	t.Run("by name", func(t *testing.T) {
		cfg := newTestCfg()
		seedPokemon(cfg, "pidgey")
		if err := commandRelease(cfg, []string{"pidgey"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if _, ok := cfg.Pokedex["pidgey"]; ok {
			t.Error("expected pidgey to be released")
		}
	})
	t.Run("by nickname", func(t *testing.T) {
		cfg := newTestCfg()
		seedPokemon(cfg, "pidgey")
		p := cfg.Pokedex["pidgey"]
		p.Nickname = "skybird"
		cfg.Pokedex["pidgey"] = p
		if err := commandRelease(cfg, []string{"skybird"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if _, ok := cfg.Pokedex["pidgey"]; ok {
			t.Error("expected release by nickname to work")
		}
	})
	t.Run("not caught", func(t *testing.T) {
		cfg := newTestCfg()
		if err := commandRelease(cfg, []string{"mewtwo"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
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
