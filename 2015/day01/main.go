// Advent of Code 2015, day 1: Not Quite Lisp.
// https://adventofcode.com/2015/day/1
package main

import (
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 1, part1, part2)
}

// part1 follows every instruction and returns the floor Santa ends on.
// An opening parenthesis goes up one floor, a closing one goes down.
func part1(in string) any {
	floor := 0

	for i := 0; i < len(in); i++ {
		floor += step(in[i])
	}

	return floor
}

// part2 returns the 1-based position of the first instruction that takes Santa
// into the basement (floor -1).
func part2(in string) any {
	floor := 0

	for i := 0; i < len(in); i++ {
		floor += step(in[i])

		if floor == -1 {
			return i + 1
		}
	}

	return nil
}

// step converts one instruction character into a floor change.
func step(b byte) int {
	if b == '(' {
		return 1
	}

	return -1
}
