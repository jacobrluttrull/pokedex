package game

import (
	"encoding/json"
	"math/rand"
	"strings"
	"testing"
	"time"

	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
)

func testMon(t *testing.T, jsonStr string) pokeapi.Pokemon {
	t.Helper()
	var p pokeapi.Pokemon
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
	ember := pokeapi.Move{Name: "ember", Power: 40, Accuracy: 100}
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
	shock := pokeapi.Move{Name: "thunder-shock", Power: 40, Accuracy: 100}
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
	cache := pokecache.NewCache(time.Minute)
	noMoves := testMon(t, `{"name":"ditto"}`)
	moves := pickMoves(noMoves, cache)
	if len(moves) != 1 || moves[0].Name != "struggle" {
		t.Errorf("expected struggle fallback, got %v", moves)
	}
}
