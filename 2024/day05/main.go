// Advent of Code 2024, day 5: Print Queue.
// https://adventofcode.com/2024/day/5
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 5, part1, part2)
}

// part1 adds the middle page numbers of the updates that are already in the
// correct order.
func part1(in string) any {
	rules, updates := parse(in)
	total := 0

	for _, pages := range updates {
		// IsSortedFunc compares only pages that are next to each other. That
		// is enough, because the rules order every pair of pages in an
		// update without a cycle (see part2).
		if slices.IsSortedFunc(pages, rules.compare) {
			total += pages[len(pages)/2]
		}
	}

	return total
}

// part2 puts each update that is in the wrong order into the correct order,
// and adds the middle page numbers of only those updates.
//
// The rules are not a total order of all pages: in real inputs they contain
// cycles. But for the pages of one update, the rules give an order for every
// pair. So a normal sort with the rules as the compare function gives the
// one correct order.
func part2(in string) any {
	rules, updates := parse(in)
	total := 0

	for _, pages := range updates {
		if slices.IsSortedFunc(pages, rules.compare) {
			continue
		}

		// SortFunc sorts in place. The update is not used again, so that is
		// safe here.
		slices.SortFunc(pages, rules.compare)
		total += pages[len(pages)/2]
	}

	return total
}

// rule holds one ordering rule: page before must come before page after.
type rule struct {
	before, after int
}

// ruleSet holds all the ordering rules. A struct key lets a map act as a set
// of pairs, and an empty struct value takes no memory.
type ruleSet map[rule]struct{}

// compare orders two pages for the slices sort functions. It returns a
// negative number when a must come before b, a positive number when b must
// come before a, and 0 when no rule orders them.
func (r ruleSet) compare(a, b int) int {
	if _, ok := r[rule{a, b}]; ok {
		return -1
	}

	if _, ok := r[rule{b, a}]; ok {
		return 1
	}

	return 0
}

// parse reads the ordering rules from the first block and the updates from
// the second block.
func parse(in string) (ruleSet, [][]int) {
	blocks := input.Blocks(in)
	if len(blocks) < 2 {
		return ruleSet{}, nil
	}

	rules := ruleSet{}
	for _, line := range blocks[0] {
		nums := input.Ints(line)
		rules[rule{nums[0], nums[1]}] = struct{}{}
	}

	var updates [][]int
	for _, line := range blocks[1] {
		updates = append(updates, input.Ints(line))
	}

	return rules, updates
}
