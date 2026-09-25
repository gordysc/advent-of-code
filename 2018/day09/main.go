// Advent of Code 2018, day 9: Marble Mania.
// https://adventofcode.com/2018/day/9
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 9, part1, part2)
}

// part1 plays the game with the numbers from the input and returns the
// winning score.
func part1(in string) any {
	nums := input.Ints(in)

	return play(nums[0], nums[1])
}

// part2 plays the same game, but the last marble is 100 times larger.
func part2(in string) any {
	nums := input.Ints(in)

	return play(nums[0], nums[1]*100)
}

// play places the marbles 1 to last in turn and returns the highest score.
//
// The circle is a doubly linked list kept in two slices: next[m] and prev[m]
// are the marbles clockwise and counter-clockwise of marble m. Each marble is
// its own index, so inserts and removals are fast and need no allocation.
func play(players, last int) int {
	next := make([]int, last+1)
	prev := make([]int, last+1)
	scores := make([]int, players)
	current := 0

	for m := 1; m <= last; m++ {
		// A multiple of 23 is kept. The player also takes the marble 7 places
		// counter-clockwise, and the marble clockwise of it becomes current.
		if m%23 == 0 {
			for range 7 {
				current = prev[current]
			}

			scores[m%players] += m + current
			next[prev[current]] = next[current]
			prev[next[current]] = prev[current]
			current = next[current]
			continue
		}

		// Any other marble goes between the marbles 1 and 2 places clockwise.
		a := next[current]
		b := next[a]

		next[a], prev[m] = m, a
		next[m], prev[b] = b, m
		current = m
	}

	return slices.Max(scores)
}
