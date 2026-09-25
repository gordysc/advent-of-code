// Advent of Code 2021, day 1: Sonar Sweep.
// https://adventofcode.com/2021/day/1
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 1, part1, part2)
}

// part1 counts the depth readings that are larger than the reading before.
func part1(in string) any {
	return countIncreases(input.IntLines(in), 1)
}

// part2 counts the three-reading windows whose sum is larger than the sum of
// the window before.
//
// Two windows next to each other share two readings. So the second sum is
// larger only when its new reading is larger than the reading that the first
// window has and the second does not. These two readings are three apart.
func part2(in string) any {
	return countIncreases(input.IntLines(in), 3)
}

// countIncreases counts the readings that are larger than the reading gap
// positions before them.
func countIncreases(depths []int, gap int) int {
	count := 0

	for i := gap; i < len(depths); i++ {
		if depths[i] > depths[i-gap] {
			count++
		}
	}

	return count
}
