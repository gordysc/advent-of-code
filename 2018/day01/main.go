// Advent of Code 2018, day 1: Chronal Calibration.
// https://adventofcode.com/2018/day/1
package main

import (
	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 1, part1, part2)
}

// part1 adds up every frequency change, starting from 0.
func part1(in string) any {
	return slicesx.Sum(input.Ints(in))
}

// part2 applies the changes over and over, from the start of the list each
// time, until the running frequency reaches a value it has had before.
func part2(in string) any {
	changes := input.Ints(in)
	seen := map[int]bool{0: true}
	freq := 0

	for {
		for _, c := range changes {
			freq += c
			if seen[freq] {
				return freq
			}

			seen[freq] = true
		}
	}
}
