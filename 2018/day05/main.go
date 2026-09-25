// Advent of Code 2018, day 5: Alchemical Reduction.
// https://adventofcode.com/2018/day/5
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 5, part1, part2)
}

// part1 counts the units left after the polymer fully reacts.
func part1(in string) any {
	return len(react(strings.TrimSpace(in), 0))
}

// part2 removes every unit of one type, in both polarities, and counts what is
// left after the rest reacts. It returns the shortest result over all types.
// Removing a type and then reacting gives the same result whether or not the
// polymer has reacted first, so each try starts from the shorter part 1
// result.
func part2(in string) any {
	reduced := string(react(strings.TrimSpace(in), 0))
	best := len(reduced)

	for unit := byte('a'); unit <= 'z'; unit++ {
		best = min(best, len(react(reduced, unit)))
	}

	return best
}

// react runs the polymer through a stack. Each unit either cancels the unit on
// top of the stack, when they are the same type with opposite polarity, or is
// pushed. Units of the type skip, in either case, are dropped. A skip of 0
// drops nothing.
func react(polymer string, skip byte) []byte {
	stack := make([]byte, 0, len(polymer))

	for i := range len(polymer) {
		c := polymer[i]
		if c|0x20 == skip {
			continue
		}

		// Upper and lower case ASCII letters differ only in bit 0x20.
		if n := len(stack); n > 0 && stack[n-1]^c == 0x20 {
			stack = stack[:n-1]
			continue
		}

		stack = append(stack, c)
	}

	return stack
}
