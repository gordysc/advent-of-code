// Advent of Code 2022, day 1: Calorie Counting.
// https://adventofcode.com/2022/day/1
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 1, part1, part2)
}

// part1 finds the elf that carries the most calories and gives that total.
func part1(in string) any {
	return topTotal(in, 1)
}

// part2 gives the sum of the calories that the top three elves carry.
func part2(in string) any {
	return topTotal(in, 3)
}

// topTotal adds up the calories of each elf, sorts the totals from largest
// to smallest, and gives the sum of the first n totals.
func topTotal(in string, n int) int {
	var totals []int

	// Each block of lines is the list of items that one elf carries.
	for _, block := range input.Blocks(in) {
		totals = append(totals, slicesx.Sum(input.Ints(strings.Join(block, "\n"))))
	}

	// Sort from smallest to largest, then reverse to put the largest first.
	slices.Sort(totals)
	slices.Reverse(totals)

	return slicesx.Sum(totals[:n])
}
