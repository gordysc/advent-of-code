// Advent of Code 2017, day 2: Corruption Checksum.
// https://adventofcode.com/2017/day/2
package main

import (
	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 2, part1, part2)
}

// part1 adds up the difference between the largest and smallest value of
// each row.
func part1(in string) any {
	sum := 0

	for _, line := range input.Lines(in) {
		lo, hi := slicesx.MinMax(input.Ints(line))
		sum += hi - lo
	}

	return sum
}

// part2 adds up the result of the one even division in each row.
func part2(in string) any {
	sum := 0

	for _, line := range input.Lines(in) {
		sum += evenQuotient(input.Ints(line))
	}

	return sum
}

// evenQuotient finds the only pair of values in a row where one divides the
// other with no remainder, and returns the result of that division. A row
// with no such pair gives 0.
func evenQuotient(row []int) int {
	for i, a := range row {
		for j, b := range row {
			if i != j && a%b == 0 {
				return a / b
			}
		}
	}

	return 0
}
