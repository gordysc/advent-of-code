// Advent of Code 2018, day 8: Memory Maneuver.
// https://adventofcode.com/2018/day/8
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 8, part1, part2)
}

// part1 adds up the metadata entries of every node in the tree.
func part1(in string) any {
	sum, _, _ := walk(input.Ints(in), 0)

	return sum
}

// part2 finds the value of the root node.
func part2(in string) any {
	_, value, _ := walk(input.Ints(in), 0)

	return value
}

// walk reads the node that starts at index i. A node is a header with the
// number of children and the number of metadata entries, then the children,
// then the metadata. It returns the metadata sum of the whole subtree, the
// value of the node, and the index just after the node.
//
// A node with no children has the sum of its own metadata as its value.
// Otherwise each metadata entry is a 1-based index into the children, and the
// value is the sum of the values of those children. Indexes that do not point
// to a child add nothing.
func walk(nums []int, i int) (sum, value, next int) {
	children, entries := nums[i], nums[i+1]
	i += 2

	values := make([]int, children)
	for c := range children {
		// childSum is declared first so that the result can go straight into
		// values[c] and i. The := form cannot assign to a slice element.
		var childSum int
		childSum, values[c], i = walk(nums, i)
		sum += childSum
	}

	for _, meta := range nums[i : i+entries] {
		sum += meta

		switch {
		case children == 0:
			value += meta
		case meta >= 1 && meta <= children:
			value += values[meta-1]
		}
	}

	return sum, value, i + entries
}
