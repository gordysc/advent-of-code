// Advent of Code 2021, day 7: The Treachery of Whales.
// https://adventofcode.com/2021/day/7
package main

import (
	"math"
	"slices"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 7, part1, part2)
}

// part1 finds the least fuel to align all crabs when each step costs 1.
//
// The sum of distances is smallest at the median position, so the crabs
// move there.
func part1(in string) any {
	crabs := input.Ints(in)
	slices.Sort(crabs)

	target := crabs[len(crabs)/2]

	fuel := 0
	for _, c := range crabs {
		fuel += mathx.Abs(c - target)
	}

	return fuel
}

// part2 finds the least fuel to align all crabs when the nth step costs n.
//
// A move of distance d costs the triangular number of d. The positions span
// only a few thousand values, so the code tries each position from the
// lowest crab to the highest and keeps the cheapest.
func part2(in string) any {
	crabs := input.Ints(in)
	lo, hi := slicesx.MinMax(crabs)

	best := math.MaxInt

	for target := lo; target <= hi; target++ {
		fuel := 0
		for _, c := range crabs {
			fuel += mathx.Triangular(mathx.Abs(c - target))
		}

		best = min(best, fuel)
	}

	return best
}
