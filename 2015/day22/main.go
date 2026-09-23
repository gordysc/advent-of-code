// Advent of Code 2015, day 22: Wizard Simulator 20XX.
// https://adventofcode.com/2015/day/22
//
// The worked fights in the puzzle text give the player 10 hit points and 250
// mana, not the 50 and 500 of the real puzzle, so they do not fit this code.
// The example.txt file here is a handmade boss (51 hit points, 9 damage).
// Part 1 gives 900 and part 2 gives 1216.
package main

import (
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// Your starting stats are part of the puzzle text and are the same for every
// input.
const (
	playerHP   = 50
	playerMana = 500
)

// spell is one spell from the puzzle text. An effect spell has a duration and
// sets its timer in the state. An instant spell has a duration of 0.
type spell struct {
	cost     int
	damage   int
	heal     int
	duration int
}

// The spell list. The index of each spell is also its kind in cast.
const (
	magicMissile = iota
	drain
	shield
	poison
	recharge
)

var spells = []spell{
	magicMissile: {cost: 53, damage: 4},
	drain:        {cost: 73, damage: 2, heal: 2},
	shield:       {cost: 113, duration: 6},
	poison:       {cost: 173, duration: 6},
	recharge:     {cost: 229, duration: 5},
}

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 22, part1, part2)
}

// part1 finds the least mana you can spend and still win the fight.
func part1(in string) any {
	return leastMana(parseBoss(in), false)
}

// part2 is the same search on hard mode: you lose 1 hit point at the start of
// each of your turns.
func part2(in string) any {
	return leastMana(parseBoss(in), true)
}

// boss holds the stats from the input.
type boss struct {
	hp     int
	damage int
}

// parseBoss reads the two numbers in the input: "Hit Points: 51" and
// "Damage: 9".
func parseBoss(in string) boss {
	nums := input.Ints(in)
	if len(nums) != 2 {
		panic("expected hit points and damage")
	}

	return boss{hp: nums[0], damage: nums[1]}
}

// state is the fight at the start of one of your turns. The three timers are
// the turns left on each effect.
type state struct {
	playerHP int
	mana     int
	bossHP   int
	shield   int
	poison   int
	recharge int
}

// leastMana runs Dijkstra's algorithm over fight states. Each edge is one
// round (your turn, then the boss's turn), and its cost is the mana of the
// spell you cast.
func leastMana(b boss, hard bool) int {
	start := state{playerHP: playerHP, mana: playerMana, bossHP: b.hp}

	next := func(s state) []search.Edge[state] {
		s, alive := startTurn(s, hard)
		if !alive {
			return nil
		}

		// Poison can kill the boss before you cast anything, which costs no mana.
		if s.bossHP <= 0 {
			return []search.Edge[state]{{To: s, Cost: 0}}
		}

		var edges []search.Edge[state]

		for kind := range spells {
			if after, ok := castAndDefend(s, kind, b); ok {
				edges = append(edges, search.Edge[state]{To: after, Cost: spells[kind].cost})
			}
		}

		return edges
	}

	won := func(s state) bool {
		return s.bossHP <= 0
	}

	cost, ok := search.Dijkstra(start, next, won)
	if !ok {
		panic("no way to win the fight")
	}

	return cost
}

// startTurn runs the start of your turn: the hard mode hit point loss, then
// the effects. It returns false when you die.
func startTurn(s state, hard bool) (state, bool) {
	if hard {
		s.playerHP--
		if s.playerHP <= 0 {
			return s, false
		}
	}

	return applyEffects(s), true
}

// castAndDefend casts the spell kind and then plays the boss's turn. It
// returns false when you cannot cast the spell or when you die. When the boss
// dies, the returned state has bossHP <= 0 and the fight stops there.
func castAndDefend(s state, kind int, b boss) (state, bool) {
	sp := spells[kind]
	if sp.cost > s.mana || !cast(&s, kind) {
		return s, false
	}

	s.mana -= sp.cost
	s.bossHP -= sp.damage
	s.playerHP += sp.heal
	if s.bossHP <= 0 {
		return s, true
	}

	armor := 0
	if s.shield > 0 {
		armor = 7
	}

	s = applyEffects(s)
	if s.bossHP <= 0 {
		return s, true
	}

	s.playerHP -= max(1, b.damage-armor)

	return s, s.playerHP > 0
}

// cast starts the effect of an effect spell. It returns false when that effect
// is still active, because you cannot cast it again until it ends.
func cast(s *state, kind int) bool {
	var timer *int

	switch kind {
	case shield:
		timer = &s.shield
	case poison:
		timer = &s.poison
	case recharge:
		timer = &s.recharge
	default:
		return true
	}

	if *timer > 0 {
		return false
	}

	*timer = spells[kind].duration

	return true
}

// applyEffects runs the active effects at the start of a turn and counts down
// their timers. Shield has no direct effect here: its armor applies only while
// its timer is above 0, which the boss's turn checks.
func applyEffects(s state) state {
	if s.poison > 0 {
		s.bossHP -= 3
		s.poison--
	}

	if s.recharge > 0 {
		s.mana += 101
		s.recharge--
	}

	if s.shield > 0 {
		s.shield--
	}

	return s
}
