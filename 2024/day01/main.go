// Advent of Code 2024, day 1: Historian Hysteria.
// https://adventofcode.com/2024/day/1
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 1, part1, part2)
}

// part1 pairs the smallest number of the left list with the smallest number
// of the right list, then the second smallest of each, and so on. It returns
// the sum of the distances in all pairs.
func part1(in string) any {
	left, right := parse(in)

	// slices.Sort sorts in place. The pairs only need the order, so the
	// original order of the lists does not matter.
	slices.Sort(left)
	slices.Sort(right)

	total := 0
	for i := range left {
		total += mathx.Abs(left[i] - right[i])
	}

	return total
}

// part2 returns the similarity score. Each number in the left list adds
// itself times the number of times it occurs in the right list.
func part2(in string) any {
	left, right := parse(in)

	// Count the right list once, so each lookup is fast. A missing key in a
	// Go map gives the zero value, so a number that never occurs counts 0.
	counts := map[int]int{}
	for _, n := range right {
		counts[n]++
	}

	total := 0
	for _, n := range left {
		total += n * counts[n]
	}

	return total
}

// parse reads the two columns of location IDs into two lists.
func parse(in string) ([]int, []int) {
	var left, right []int

	for _, line := range input.Lines(in) {
		nums := input.Ints(line)
		if len(nums) != 2 {
			continue
		}

		left = append(left, nums[0])
		right = append(right, nums[1])
	}

	return left, right
}
