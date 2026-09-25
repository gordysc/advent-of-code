// Advent of Code 2023, day 6: Wait For It.
// https://adventofcode.com/2023/day/6
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 6, part1, part2)
}

// part1 multiplies together the number of ways to win each race.
func part1(in string) any {
	lines := input.Lines(in)
	times := input.UInts(lines[0])
	records := input.UInts(lines[1])
	product := 1

	for i := range times {
		product *= waysToWin(times[i], records[i])
	}

	return product
}

// part2 counts the ways to win one long race. The spaces between the digits
// are not real, so the part removes them before it reads the two numbers.
func part2(in string) any {
	lines := input.Lines(in)
	time := input.UInts(strings.ReplaceAll(lines[0], " ", ""))[0]
	record := input.UInts(strings.ReplaceAll(lines[1], " ", ""))[0]

	return waysToWin(time, record)
}

// waysToWin counts the hold times that beat the record in a race of the given
// time.
//
// A hold of h milliseconds moves the boat h*(time-h) millimetres. This
// distance goes up until h reaches time/2 and then goes down in the same way,
// because h and time-h give the same distance. So a binary search on the
// first half finds the shortest winning hold, and every hold from it up to
// time minus it also wins. The part 2 race is too long to test each hold.
func waysToWin(time, record int) int {
	half := time / 2

	first := search.BinarySearch(0, half, func(h int) bool {
		return h*(time-h) > record
	})
	if first > half {
		return 0
	}

	return time - 2*first + 1
}
