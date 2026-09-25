// Advent of Code 2021, day 6: Lanternfish.
// https://adventofcode.com/2021/day/6
package main

import (
	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 6, part1, part2)
}

// part1 counts the lanternfish after 80 days.
func part1(in string) any {
	return simulate(in, 80)
}

// part2 counts the lanternfish after 256 days.
//
// The count grows too fast to keep each fish, so both parts only keep the
// number of fish for each timer value.
func part2(in string) any {
	return simulate(in, 256)
}

// simulate returns the number of fish after the given number of days.
//
// counts[t] is the number of fish with timer t. Each day, all timers go down
// by one. The fish at timer 0 reset to 6 and each makes a new fish at 8.
func simulate(in string, days int) int {
	var counts [9]int

	for _, t := range input.Ints(in) {
		counts[t]++
	}

	for range days {
		spawning := counts[0]

		copy(counts[:], counts[1:])
		counts[6] += spawning
		counts[8] = spawning
	}

	return slicesx.Sum(counts[:])
}
