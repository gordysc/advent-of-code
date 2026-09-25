// Advent of Code 2017, day 5: A Maze of Twisty Trampolines, All Alike.
// https://adventofcode.com/2017/day/5
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 5, part1, part2)
}

// part1 counts the steps to leave the list when every jump adds one to its
// own offset.
func part1(in string) any {
	return escape(input.IntLines(in), func(int) int { return 1 })
}

// part2 counts the steps when an offset of three or more goes down by one
// after its jump, and any other offset goes up by one.
func part2(in string) any {
	return escape(input.IntLines(in), func(offset int) int {
		if offset >= 3 {
			return -1
		}

		return 1
	})
}

// escape follows the jumps from the first offset until a jump leaves the
// list, and returns the number of steps. After each jump, change gives the
// amount to add to the offset that was just used.
func escape(offsets []int, change func(int) int) int {
	steps := 0

	for i := 0; i >= 0 && i < len(offsets); steps++ {
		offset := offsets[i]
		offsets[i] += change(offset)
		i += offset
	}

	return steps
}
