// Advent of Code 2022, day 25: Full of Hot Air.
// https://adventofcode.com/2022/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"aoc/lib/input"
	"aoc/lib/strx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 25, part1, nil)
}

// part1 adds all the SNAFU numbers and returns the sum as a SNAFU number.
func part1(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		total += fromSnafu(line)
	}

	return toSnafu(total)
}

// snafuDigits maps each SNAFU digit to its value.
var snafuDigits = map[rune]int{'=': -2, '-': -1, '0': 0, '1': 1, '2': 2}

// fromSnafu converts a SNAFU number to an int. SNAFU is base 5, but each
// digit has a value from -2 to 2.
func fromSnafu(s string) int {
	n := 0

	for _, r := range s {
		n = n*5 + snafuDigits[r]
	}

	return n
}

// toSnafu converts a positive int to a SNAFU number.
//
// It finds the digits from the right. A normal base-5 remainder of 3 or 4 is
// not a SNAFU digit, so it writes the remainder minus 5 ('=' or '-') and
// carries 1 to the next digit.
func toSnafu(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []byte

	for n > 0 {
		// Adding 2 before the division moves the remainders 3 and 4 into
		// the next place, which is the carry.
		rem := (n + 2) % 5
		digits = append(digits, "=-012"[rem])
		n = (n + 2) / 5
	}

	// The loop finds the lowest digit first, so reverse the result.
	return strx.Reverse(string(digits))
}
