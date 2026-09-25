// Advent of Code 2025, day 3: Lobby.
// https://adventofcode.com/2025/day/3
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 3, part1, part2)
}

// part1 adds the largest joltage of each bank when we turn on two batteries.
func part1(in string) any {
	return totalJoltage(in, 2)
}

// part2 adds the largest joltage of each bank when we turn on twelve
// batteries.
func part2(in string) any {
	return totalJoltage(in, 12)
}

// totalJoltage adds the largest joltage of each bank with count batteries on.
func totalJoltage(in string, count int) int {
	total := 0

	for _, line := range input.Lines(in) {
		if len(line) < count {
			continue
		}

		total += maxJoltage(line, count)
	}

	return total
}

// maxJoltage returns the largest number we can make from count digits of bank,
// with the digits kept in their order.
//
// A larger first digit always gives a larger number, whatever the digits after
// it are. So we pick greedily: each digit is the largest digit that still
// leaves enough digits after it for the rest. If two digits are equal, the
// leftmost one is the best, because it leaves the most choices after it.
func maxJoltage(bank string, count int) int {
	value := 0
	start := 0

	for left := count; left > 0; left-- {
		// The last position we can use still leaves left-1 digits after it.
		best := start
		for i := start + 1; i <= len(bank)-left; i++ {
			if bank[i] > bank[best] {
				best = i
			}
		}

		// bank[best] is a byte such as '7'. Subtracting '0' gives the digit 7.
		value = value*10 + int(bank[best]-'0')
		start = best + 1
	}

	return value
}
