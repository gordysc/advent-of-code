// Advent of Code 2020, day 10: Adapter Array.
// https://adventofcode.com/2020/day/10
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 10, part1, part2)
}

// part1 uses every adapter in a chain from the outlet to the device. The
// chain must go in sorted order, so the answer is the count of 1-jolt steps
// times the count of 3-jolt steps between neighbours in the sorted list.
//
// example.txt is the second, larger example. It gives 220 for part 1 and
// 19208 for part 2.
func part1(in string) any {
	chain := parse(in)

	var diffs [4]int
	for i := 1; i < len(chain); i++ {
		diffs[chain[i]-chain[i-1]]++
	}

	return diffs[1] * diffs[3]
}

// part2 counts the distinct chains from the outlet to the device. The count
// can be in the trillions, so it cannot list the chains. Instead, it walks
// the sorted list and counts the ways to reach each joltage: the sum of the
// ways to reach each earlier joltage at most 3 below it.
func part2(in string) any {
	chain := parse(in)
	ways := make([]int, len(chain))
	ways[0] = 1

	for i := 1; i < len(chain); i++ {
		for j := i - 1; j >= 0 && chain[i]-chain[j] <= 3; j-- {
			ways[i] += ways[j]
		}
	}

	return ways[len(ways)-1]
}

// parse reads the adapters and returns the full sorted chain. It adds the
// outlet (0 jolts) at the start and the device (3 more than the largest
// adapter) at the end.
func parse(in string) []int {
	chain := append([]int{0}, input.IntLines(in)...)
	slices.Sort(chain)

	return append(chain, chain[len(chain)-1]+3)
}
