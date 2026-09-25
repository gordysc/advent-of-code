// Advent of Code 2020, day 1: Report Repair.
// https://adventofcode.com/2020/day/1
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// target is the sum that the expense entries must make.
const target = 2020

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 1, part1, part2)
}

// part1 finds the two entries that sum to 2020 and multiplies them.
func part1(in string) any {
	nums := input.IntLines(in)

	a, b, ok := pairSum(nums, target)
	if !ok {
		return nil
	}

	return a * b
}

// part2 finds the three entries that sum to 2020 and multiplies them. For
// each entry, it looks for a pair in the remaining entries that makes up the
// rest of the sum. This is O(n^2), which is fast for a list of about 200.
func part2(in string) any {
	nums := input.IntLines(in)

	for i, a := range nums {
		b, c, ok := pairSum(nums[i+1:], target-a)
		if ok {
			return a * b * c
		}
	}

	return nil
}

// pairSum finds two entries in nums that sum to want. It keeps a set of the
// entries it has seen, so each entry needs one lookup for its complement.
func pairSum(nums []int, want int) (int, int, bool) {
	seen := map[int]bool{}

	for _, n := range nums {
		if seen[want-n] {
			return want - n, n, true
		}

		seen[n] = true
	}

	return 0, 0, false
}
