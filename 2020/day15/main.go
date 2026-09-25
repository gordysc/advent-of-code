// Advent of Code 2020, day 15: Rambunctious Recitation.
// https://adventofcode.com/2020/day/15
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 15, part1, part2)
}

// part1 gives the 2020th number spoken in the memory game.
func part1(in string) any {
	return play(input.Ints(in), 2020)
}

// part2 gives the 30000000th number spoken in the memory game.
func part2(in string) any {
	return play(input.Ints(in), 30_000_000)
}

// play runs the memory game for the given number of turns and returns the
// last number spoken. Each turn, the next number is 0 if the current number
// is new. Otherwise it is the gap since the previous turn that spoke it.
//
// A spoken number is always smaller than the number of turns, so a flat
// slice indexed by number holds the last turn for each number. A map is
// much slower for 30 million turns. The slice uses int32 to halve the memory,
// and 0 means "not spoken yet", so the turns count from 1.
func play(start []int, turns int) int {
	last := make([]int32, max(turns, len(start)+1))
	for i, n := range start[:len(start)-1] {
		last[n] = int32(i + 1)
	}

	current := start[len(start)-1]
	for turn := len(start); turn < turns; turn++ {
		next := 0
		if seen := last[current]; seen != 0 {
			next = turn - int(seen)
		}

		last[current] = int32(turn)
		current = next
	}

	return current
}
