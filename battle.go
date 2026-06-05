package main

import "fmt"

func getStat(p Pokemon, name string) int {
	for _, s := range p.Stats {
		if s.Stat.Name == name {
			return s.BaseStat
		}
	}
	return 1
}

func runBattle(yours, wild Pokemon) bool {
	yourHP := getStat(yours, "hp")
	wildHP := getStat(wild, "hp")
	yourAtk := getStat(yours, "attack")
	wildAtk := getStat(wild, "attack")

	fmt.Printf("--- %s vs %s ---\n", yours.Name, wild.Name)
	for yourHP > 0 && wildHP > 0 {
		wildHP -= yourAtk
		fmt.Printf("%s attacks! %s HP: %d\n", yours.Name, wild.Name, max(wildHP, 0))
		if wildHP <= 0 {
			break
		}
		yourHP -= wildAtk
		fmt.Printf("%s attacks! %s HP: %d\n", wild.Name, yours.Name, max(yourHP, 0))
	}
	if yourHP > 0 {
		fmt.Printf("%s won!\n", yours.Name)
		return true
	}
	fmt.Printf("%s lost!\n", yours.Name)
	return false
}
