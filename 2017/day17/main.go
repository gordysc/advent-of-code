// Advent of Code 2017, day 17: Spinlock.
// https://adventofcode.com/2017/day/17
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 17, part1, part2)
}

// The number of values each part inserts into the buffer.
const (
	shortRun = 2017
	longRun  = 50_000_000
)

// part1 builds the buffer with 2017 inserts and returns the value just after
// the last one.
func part1(in string) any {
	step := input.Int(in)
	buf := []int{0}
	pos := 0

	for n := 1; n <= shortRun; n++ {
		pos = (pos+step)%len(buf) + 1
		buf = slices.Insert(buf, pos, n)
	}

	return buf[(pos+1)%len(buf)]
}

// part2 returns the value just after 0 once 50 million values are in. Value 0
// never moves from position 0, so the answer is the last value inserted at
// position 1. The buffer itself is never built; only the current position is.
func part2(in string) any {
	step := input.Int(in)
	pos, after := 0, 0

	// Before value n goes in, the buffer holds n values.
	for n := 1; n <= longRun; {
		pos = (pos+step)%n + 1
		if pos == 1 {
			after = n
		}

		// The next inserts that do not wrap past the end each move forward by
		// step+1 and can never land at position 1. Skip them all at once.
		skip := (n - pos) / step
		pos += skip * (step + 1)
		n += skip + 1
	}

	return after
}
