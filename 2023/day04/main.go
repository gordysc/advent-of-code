// Advent of Code 2023, day 4: Scratchcards.
// https://adventofcode.com/2023/day/4
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 4, part1, part2)
}

// part1 adds the points of all cards. A card with n matches is worth
// 2^(n-1) points, and a card with no matches is worth nothing.
func part1(in string) any {
	total := 0

	for _, m := range matchCounts(in) {
		if m > 0 {
			// A left shift by m-1 doubles 1 that many times, so it gives
			// 2^(m-1).
			total += 1 << (m - 1)
		}
	}

	return total
}

// part2 counts the cards at the end, originals and copies together. A card
// with n matches wins one copy of each of the next n cards.
//
// A card only wins cards below it. So one pass from top to bottom is enough:
// when the pass reaches a card, the count of its copies is already final, and
// every one of these copies wins the same next cards.
func part2(in string) any {
	matches := matchCounts(in)
	copies := make([]int, len(matches))

	for i := range copies {
		copies[i] = 1
	}

	total := 0

	for i, m := range matches {
		total += copies[i]

		for j := i + 1; j <= i+m && j < len(copies); j++ {
			copies[j] += copies[i]
		}
	}

	return total
}

// matchCounts gives, for each card in order, how many of its numbers are
// also winning numbers.
func matchCounts(in string) []int {
	var counts []int

	for _, line := range input.Lines(in) {
		_, numbers, _ := strings.Cut(line, ": ")
		winning, have, _ := strings.Cut(numbers, " | ")

		wins := set.From(input.UInts(winning))
		count := 0

		for _, n := range input.UInts(have) {
			if wins.Has(n) {
				count++
			}
		}

		counts = append(counts, count)
	}

	return counts
}
