// Advent of Code 2022, day 3: Rucksack Reorganization.
// https://adventofcode.com/2022/day/3
package main

import (
	"math/bits"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 3, part1, part2)
}

// part1 finds the one item type that is in both halves of each rucksack and
// gives the sum of the priorities of these items.
func part1(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		half := len(line) / 2
		total += common(line[:half], line[half:])
	}

	return total
}

// part2 finds the badge of each group of three elves. The badge is the one
// item type that all three rucksacks of the group have. It gives the sum of
// the priorities of the badges.
func part2(in string) any {
	lines := input.Lines(in)
	total := 0

	for i := 0; i+2 < len(lines); i += 3 {
		total += common(lines[i], lines[i+1], lines[i+2])
	}

	return total
}

// common gives the priority of the item type that all the given strings
// contain.
//
// There are only 52 priorities, so each string becomes a 64-bit mask with
// one bit for each priority that it contains. An AND of the masks keeps only
// the shared bit, and the position of that bit is the priority.
func common(items ...string) int {
	shared := ^uint64(0)

	for _, s := range items {
		shared &= mask(s)
	}

	return bits.TrailingZeros64(shared)
}

// mask sets the bit at the priority of each item in s.
func mask(s string) uint64 {
	var m uint64

	for i := 0; i < len(s); i++ {
		m |= 1 << priority(s[i])
	}

	return m
}

// priority gives 1 to 26 for the letters a to z and 27 to 52 for the
// letters A to Z.
func priority(c byte) int {
	if c >= 'a' {
		return int(c-'a') + 1
	}

	return int(c-'A') + 27
}
