// Advent of Code 2023, day 12: Hot Springs.
// https://adventofcode.com/2023/day/12
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 12, part1, part2)
}

// part1 adds up the number of possible arrangements of each row.
func part1(in string) any {
	return sumArrangements(in, 1)
}

// part2 unfolds each row to five copies and adds up the arrangements again.
//
// An unfolded row can have billions of arrangements, so a count of each one
// is not possible. The count uses memoized dynamic programming instead (see
// arrangements).
func part2(in string) any {
	return sumArrangements(in, 5)
}

// sumArrangements unfolds each row to copies copies and adds up the number of
// arrangements of all rows.
func sumArrangements(in string, copies int) int {
	total := 0

	for _, line := range input.Lines(in) {
		springs, groupText, _ := strings.Cut(line, " ")
		groups := input.Ints(groupText)

		// The copies of the springs are joined with a '?'. The copies of the
		// groups are joined with nothing.
		springs = strings.Repeat(springs+"?", copies)
		springs = springs[:len(springs)-1]
		groups = slices.Repeat(groups, copies)

		total += arrangements(springs, groups)
	}

	return total
}

// arrangements counts the ways to put the damaged groups onto the springs.
//
// count(i, j) is the number of ways to put groups j and after onto springs i
// and after. At spring i there are two choices. If the spring can be
// operational ('.' or '?'), skip it. If it can be damaged ('#' or '?'), put
// group j there. The group needs g springs that can all be damaged, and then
// a spring that can be operational (or the end of the row).
//
// Many paths get to the same (i, j), so a memo keeps each count. There are
// only len(springs) * len(groups) states, and each state takes O(g) work.
func arrangements(springs string, groups []int) int {
	n := len(springs)

	// memo has one row per spring position (0 to n) and one column per group
	// index (0 to len(groups)). -1 marks a state that is not known yet.
	memo := make([][]int, n+1)
	for i := range memo {
		memo[i] = slices.Repeat([]int{-1}, len(groups)+1)
	}

	// count is a closure, so it can see springs, groups and memo. It must be
	// declared before the assignment so that it can call itself.
	var count func(i, j int) int
	count = func(i, j int) int {
		// After a group at the end of the row, i can be one past the end.
		i = min(i, n)

		// When all groups are placed, no damaged spring must be left over.
		if j == len(groups) {
			if strings.Contains(springs[i:], "#") {
				return 0
			}

			return 1
		}

		if i == n {
			return 0
		}

		if memo[i][j] >= 0 {
			return memo[i][j]
		}

		ways := 0

		if springs[i] != '#' {
			ways += count(i+1, j)
		}

		if springs[i] != '.' {
			// The && operator stops at the first false part, so each index
			// is checked against n before it is read.
			end := i + groups[j]
			fits := end <= n && !strings.Contains(springs[i:end], ".")

			if fits && (end == n || springs[end] != '#') {
				ways += count(end+1, j+1)
			}
		}

		memo[i][j] = ways

		return ways
	}

	return count(0, 0)
}
