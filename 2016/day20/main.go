// Advent of Code 2016, day 20: Firewall Rules.
// https://adventofcode.com/2016/day/20
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 20, part1, part2)
}

// maxIP is the highest IP address in the puzzle. The worked example in the
// puzzle text uses 9 instead. The lowest allowed IP in the example is 3 with
// both limits. With a limit of 9 the example allows 2 IPs (3 and 9), but with
// this limit part 2 gives 4294967288 for the example (3, and 9 to maxIP).
const maxIP = 4294967295

// part1 finds the lowest IP that no range blocks.
func part1(in string) any {
	open := gaps(in)
	if len(open) == 0 {
		return nil
	}

	return open[0][0]
}

// part2 counts every IP that no range blocks.
func part2(in string) any {
	count := 0
	for _, g := range gaps(in) {
		count += g[1] - g[0] + 1
	}

	return count
}

// gaps reads the blocked ranges and returns the ranges of allowed IPs, lowest
// first. Each range is a pair of the first and last IP, both included.
//
// The blocked ranges are sorted by their start. A single pass then tracks the
// first IP that no range seen so far covers. A range that starts after that
// IP leaves a gap in front of it.
func gaps(in string) [][2]int {
	var blocked [][2]int
	for _, line := range input.Lines(in) {
		n := input.UInts(line)
		blocked = append(blocked, [2]int{n[0], n[1]})
	}

	slices.SortFunc(blocked, func(a, b [2]int) int { return a[0] - b[0] })

	var open [][2]int
	next := 0

	for _, r := range blocked {
		if r[0] > next {
			open = append(open, [2]int{next, r[0] - 1})
		}

		next = max(next, r[1]+1)
	}

	if next <= maxIP {
		open = append(open, [2]int{next, maxIP})
	}

	return open
}
