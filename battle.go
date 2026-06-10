package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"pokedex/internal/pokecache"
)

const battleLevel = 50

func getStat(p Pokemon, name string) int {
	for _, s := range p.Stats {
		if s.Stat.Name == name {
			return s.BaseStat
		}
	}
	return 1
}

func struggleMove() Move {
	m := Move{Name: "struggle", Power: 50, Accuracy: 100}
	m.Type.Name = "normal"
	return m
}

// pickMoves fetches up to 4 damaging moves for a pokemon. Status moves
// (power 0) are skipped; fetches are capped so a long move list can't
// hammer the API. Falls back to struggle if nothing usable is found.
func pickMoves(p Pokemon, cache *pokecache.Cache) []Move {
	moves := []Move{}
	attempts := 0
	for _, m := range p.Moves {
		if len(moves) == 4 || attempts == 8 {
			break
		}
		attempts++
		mv, err := fetchMove(m.Move.URL, cache)
		if err != nil || mv.Power == 0 {
			continue
		}
		moves = append(moves, mv)
	}
	if len(moves) == 0 {
		moves = append(moves, struggleMove())
	}
	return moves
}

// damageRoll resolves one attack: accuracy check, Gen 4 damage formula
// at a flat level, STAB, type effectiveness, crit chance, and variance.
// Returns the damage dealt and the battle text to print.
func damageRoll(attacker, defender Pokemon, mv Move, rng *rand.Rand) (int, []string) {
	msgs := []string{fmt.Sprintf("%s used %s!", attacker.Name, mv.Name)}

	accuracy := mv.Accuracy
	if accuracy == 0 {
		accuracy = 100
	}
	if rng.Intn(100) >= accuracy {
		return 0, append(msgs, "...but it missed!")
	}

	eff := typeEffect(mv.Type.Name, typeNames(defender))
	if eff == 0 {
		return 0, append(msgs, fmt.Sprintf("It doesn't affect %s...", defender.Name))
	}

	atk := getStat(attacker, "attack")
	def := getStat(defender, "defense")
	dmg := float64((2*battleLevel/5+2)*mv.Power*atk)/float64(def)/50.0 + 2

	for _, t := range attacker.Types {
		if t.Type.Name == mv.Type.Name {
			dmg *= 1.5
			break
		}
	}
	dmg *= eff
	if rng.Intn(16) == 0 {
		dmg *= 2
		msgs = append(msgs, "A critical hit!")
	}
	dmg *= 0.85 + rng.Float64()*0.15

	if eff > 1 {
		msgs = append(msgs, "It's super effective!")
	} else if eff < 1 {
		msgs = append(msgs, "It's not very effective...")
	}
	if int(dmg) < 1 {
		return 1, msgs
	}
	return int(dmg), msgs
}

func hpBar(name string, cur, max int) string {
	if cur < 0 {
		cur = 0
	}
	filled := 0
	if max > 0 {
		filled = cur * 10 / max
	}
	return fmt.Sprintf("%-12s [%s%s] %d/%d", name, strings.Repeat("=", filled), strings.Repeat("-", 10-filled), cur, max)
}

func runBattle(yours, wild Pokemon, cache *pokecache.Cache) bool {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	yourMoves := pickMoves(yours, cache)
	wildMoves := pickMoves(wild, cache)
	yourHP := getStat(yours, "hp")
	wildHP := getStat(wild, "hp")
	yourMax, wildMax := yourHP, wildHP

	fmt.Printf("--- %s vs %s ---\n", yours.Name, wild.Name)
	for {
		fmt.Println(hpBar(yours.Name, yourHP, yourMax))
		fmt.Println(hpBar(wild.Name, wildHP, wildMax))
		for i, m := range yourMoves {
			fmt.Printf("  %d) %s (%s, %d power)\n", i+1, m.Name, m.Type.Name, m.Power)
		}
		fmt.Print("Choose a move or 'run': ")
		var input string
		fmt.Scan(&input)
		if input == "run" {
			fmt.Println("You fled the battle!")
			return false
		}
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(yourMoves) {
			fmt.Println("invalid choice")
			continue
		}
		yourMove := yourMoves[choice-1]
		wildMove := wildMoves[rng.Intn(len(wildMoves))]

		yourSpeed := getStat(yours, "speed")
		wildSpeed := getStat(wild, "speed")
		youFirst := yourSpeed > wildSpeed || (yourSpeed == wildSpeed && rng.Intn(2) == 0)

		if youFirst {
			if attack(yours, wild, yourMove, &wildHP, rng) {
				fmt.Printf("The wild %s fainted! You won!\n", wild.Name)
				return true
			}
			if attack(wild, yours, wildMove, &yourHP, rng) {
				fmt.Printf("%s fainted! You lost!\n", yours.Name)
				return false
			}
		} else {
			if attack(wild, yours, wildMove, &yourHP, rng) {
				fmt.Printf("%s fainted! You lost!\n", yours.Name)
				return false
			}
			if attack(yours, wild, yourMove, &wildHP, rng) {
				fmt.Printf("The wild %s fainted! You won!\n", wild.Name)
				return true
			}
		}
	}
}

// attack applies one move to the defender's hp and reports whether it fainted.
func attack(attacker, defender Pokemon, mv Move, defenderHP *int, rng *rand.Rand) bool {
	dmg, msgs := damageRoll(attacker, defender, mv, rng)
	for _, m := range msgs {
		fmt.Println(m)
	}
	*defenderHP -= dmg
	return *defenderHP <= 0
}
