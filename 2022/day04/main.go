// Advent of Code 2022, day 4: Camp Cleanup.
// https://adventofcode.com/2022/day/4
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 4, part1, part2)
}

// part1 counts the pairs where one range fully contains the other range.
func part1(in string) any {
	return countPairs(in, func(a1, a2, b1, b2 int) bool {
		return (a1 <= b1 && b2 <= a2) || (b1 <= a1 && a2 <= b2)
	})
}

// part2 counts the pairs where the two ranges overlap at all.
//
// Two ranges overlap when each one starts before or at the end of the other.
func part2(in string) any {
	return countPairs(in, func(a1, a2, b1, b2 int) bool {
		return a1 <= b2 && b1 <= a2
	})
}

// countPairs reads each line as two ranges a1-a2 and b1-b2 and counts the
// lines where match is true. input.UInts reads the dashes as separators and
// not as minus signs.
func countPairs(in string, match func(a1, a2, b1, b2 int) bool) int {
	count := 0

	for _, line := range input.Lines(in) {
		n := input.UInts(line)

		if match(n[0], n[1], n[2], n[3]) {
			count++
		}
	}

	return count
}
