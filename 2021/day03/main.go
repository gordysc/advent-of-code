// Advent of Code 2021, day 3: Binary Diagnostic.
// https://adventofcode.com/2021/day/3
package main

import (
	"strconv"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 3, part1, part2)
}

// part1 finds the power consumption. The gamma rate uses the most common bit
// in each position, and the epsilon rate uses the least common bit. Epsilon
// is gamma with all its bits flipped.
func part1(in string) any {
	lines := input.Lines(in)
	width := len(lines[0])

	gamma := 0

	for i := range width {
		gamma <<= 1
		if mostCommon(lines, i) == '1' {
			gamma |= 1
		}
	}

	epsilon := gamma ^ (1<<width - 1)

	return gamma * epsilon
}

// part2 finds the life support rating. Each rating starts with all numbers
// and removes numbers one bit position at a time, until one number remains.
// The oxygen rating keeps the numbers with the most common bit, and the CO2
// rating keeps the numbers with the least common bit.
func part2(in string) any {
	lines := input.Lines(in)

	oxygen := filter(lines, true)
	co2 := filter(lines, false)

	return oxygen * co2
}

// mostCommon returns the most common bit at position i. A tie gives '1'.
func mostCommon(lines []string, i int) byte {
	ones := 0

	for _, line := range lines {
		if line[i] == '1' {
			ones++
		}
	}

	if 2*ones >= len(lines) {
		return '1'
	}

	return '0'
}

// filter removes numbers bit by bit until one is left, and returns it as an
// int. When most is true it keeps the numbers with the most common bit (ties
// keep '1'). Otherwise it keeps the numbers with the least common bit (ties
// keep '0').
func filter(lines []string, most bool) int {
	for i := 0; len(lines) > 1; i++ {
		want := mostCommon(lines, i)
		if !most {
			want ^= '0' ^ '1' // flip '0' and '1'
		}

		var kept []string

		for _, line := range lines {
			if line[i] == want {
				kept = append(kept, line)
			}
		}

		lines = kept
	}

	n, _ := strconv.ParseInt(lines[0], 2, 64)

	return int(n)
}
