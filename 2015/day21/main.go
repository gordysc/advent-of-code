// Advent of Code 2015, day 21: RPG Simulator 20XX.
// https://adventofcode.com/2015/day/21
//
// The puzzle text has one worked fight, but no shop example. You cannot lose
// to the boss in that fight, so part 2 has no answer for it. The example.txt
// file here is a stronger handmade boss (100 hit points, 8 damage, 2 armor).
// Part 1 gives 91 and part 2 gives 158.
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// playerHP is the number of hit points you start with. It is part of the
// puzzle text and is the same for every input.
const playerHP = 100

// item is one thing from the shop.
type item struct {
	cost   int
	damage int
	armor  int
}

// The shop is part of the puzzle text and is the same for every input.
var (
	weapons = []item{
		{8, 4, 0},  // Dagger
		{10, 5, 0}, // Shortsword
		{25, 6, 0}, // Warhammer
		{40, 7, 0}, // Longsword
		{74, 8, 0}, // Greataxe
	}

	// armors starts with an empty slot, because armor is optional.
	armors = []item{
		{0, 0, 0},   // No armor
		{13, 0, 1},  // Leather
		{31, 0, 2},  // Chainmail
		{53, 0, 3},  // Splintmail
		{75, 0, 4},  // Bandedmail
		{102, 0, 5}, // Platemail
	}

	// rings starts with two empty slots, because you can wear zero, one or two
	// rings, and the shop has only one of each real ring.
	rings = []item{
		{0, 0, 0},   // No ring
		{0, 0, 0},   // No ring
		{25, 1, 0},  // Damage +1
		{50, 2, 0},  // Damage +2
		{100, 3, 0}, // Damage +3
		{20, 0, 1},  // Defense +1
		{40, 0, 2},  // Defense +2
		{80, 0, 3},  // Defense +3
	}
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 21, part1, part2)
}

// part1 finds the least gold you can spend and still win the fight.
func part1(in string) any {
	boss := parseBoss(in)
	best := -1

	for _, l := range loadouts() {
		if wins(l, boss) && (best < 0 || l.cost < best) {
			best = l.cost
		}
	}

	return best
}

// part2 finds the most gold the shopkeeper can make you spend and still lose.
func part2(in string) any {
	boss := parseBoss(in)
	worst := 0

	for _, l := range loadouts() {
		if !wins(l, boss) {
			worst = max(worst, l.cost)
		}
	}

	return worst
}

// fighter holds the stats of one side of the fight.
type fighter struct {
	hp     int
	damage int
	armor  int
}

// parseBoss reads the three numbers in the input:
// "Hit Points: 104", "Damage: 8" and "Armor: 1".
func parseBoss(in string) fighter {
	nums := input.Ints(in)
	if len(nums) != 3 {
		panic("expected hit points, damage and armor")
	}

	return fighter{hp: nums[0], damage: nums[1], armor: nums[2]}
}

// loadouts returns the sum of every legal purchase: one weapon, zero or one
// armor, and zero to two different rings.
func loadouts() []item {
	var all []item

	for _, w := range weapons {
		for _, a := range armors {
			for i := 0; i < len(rings); i++ {
				for j := i + 1; j < len(rings); j++ {
					r1, r2 := rings[i], rings[j]

					all = append(all, item{
						cost:   w.cost + a.cost + r1.cost + r2.cost,
						damage: w.damage + a.damage + r1.damage + r2.damage,
						armor:  w.armor + a.armor + r1.armor + r2.armor,
					})
				}
			}
		}
	}

	return all
}

// wins reports whether you beat the boss with the given gear. Each attack does
// at least 1 damage, so the fight is decided by who needs fewer turns. You
// attack first, so you also win a tie.
func wins(gear item, boss fighter) bool {
	yourTurns := mathx.DivCeil(boss.hp, max(1, gear.damage-boss.armor))
	bossTurns := mathx.DivCeil(playerHP, max(1, boss.damage-gear.armor))

	return yourTurns <= bossTurns
}
