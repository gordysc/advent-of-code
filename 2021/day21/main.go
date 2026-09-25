// Advent of Code 2021, day 21: Dirac Dice.
// https://adventofcode.com/2021/day/21
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 21, part1, part2)
}

// part1 plays with the deterministic die until a player gets 1000 points.
//
// The die rolls 1, 2, 3 and so on up to 100, then starts again at 1. Each
// turn, the player rolls three times and moves forward by the sum. The answer
// is the losing score multiplied by the number of rolls.
func part1(in string) any {
	pos := startPositions(in)
	score := [2]int{}
	rolls := 0

	for turn := 0; ; turn = 1 - turn {
		move := 0
		for range 3 {
			move += rolls%100 + 1
			rolls++
		}

		pos[turn] = (pos[turn]+move-1)%10 + 1
		score[turn] += pos[turn]

		if score[turn] >= 1000 {
			return score[1-turn] * rolls
		}
	}
}

// state is one position in the Dirac game, seen from the player whose turn it
// is. Index 0 is the player who moves next.
type state struct {
	pos, score [2]int
}

// wins holds the number of universes that each player wins.
type wins [2]int

// part2 plays with the Dirac die and counts the universes each player wins.
//
// Three rolls of the Dirac die give 27 universes, but only 7 different sums.
// countWins keeps a memo of the result for each state, so each state is
// calculated only once. The answer is the larger of the two win counts.
func part2(in string) any {
	pos := startPositions(in)
	memo := map[state]wins{}
	w := countWins(state{pos: pos}, memo)

	return max(w[0], w[1])
}

// rollCounts maps each sum of three Dirac rolls to the number of universes
// that give that sum.
var rollCounts = map[int]int{3: 1, 4: 3, 5: 6, 6: 7, 7: 6, 8: 3, 9: 1}

// countWins returns the universes that each player wins from state s. Index
// 0 of the result is the player who moves next in s.
func countWins(s state, memo map[state]wins) wins {
	if w, ok := memo[s]; ok {
		return w
	}

	var total wins
	for move, n := range rollCounts {
		pos := (s.pos[0]+move-1)%10 + 1
		score := s.score[0] + pos

		if score >= 21 {
			total[0] += n
			continue
		}

		// Swap the players, so that the other player moves next.
		next := state{
			pos:   [2]int{s.pos[1], pos},
			score: [2]int{s.score[1], score},
		}
		w := countWins(next, memo)

		total[0] += n * w[1]
		total[1] += n * w[0]
	}

	memo[s] = total

	return total
}

// startPositions reads the start positions of the two players.
func startPositions(in string) [2]int {
	lines := input.Lines(in)

	// The first number on each line is the player number, the last is the
	// start position.
	a := input.Ints(lines[0])
	b := input.Ints(lines[1])

	return [2]int{a[len(a)-1], b[len(b)-1]}
}
