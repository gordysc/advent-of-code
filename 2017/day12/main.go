// Advent of Code 2017, day 12: Digital Plumber.
// https://adventofcode.com/2017/day/12
package main

import (
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 12, part1, part2)
}

// part1 counts the programs in the group that holds program 0.
func part1(in string) any {
	pipes := parse(in)

	return len(search.Flood(0, func(p int) []int { return pipes[p] }))
}

// part2 counts the groups. Each program not yet in a group starts a new one,
// and a flood fill from it marks every program in that group.
func part2(in string) any {
	pipes := parse(in)
	next := func(p int) []int { return pipes[p] }
	seen := map[int]bool{}
	groups := 0

	for p := range pipes {
		if seen[p] {
			continue
		}

		groups++
		for q := range search.Flood(p, next) {
			seen[q] = true
		}
	}

	return groups
}

// parse reads each "2 <-> 0, 3, 4" line into a map from a program to the
// programs it has a pipe to.
func parse(in string) map[int][]int {
	pipes := map[int][]int{}

	for _, line := range input.Lines(in) {
		nums := input.Ints(line)
		pipes[nums[0]] = nums[1:]
	}

	return pipes
}
