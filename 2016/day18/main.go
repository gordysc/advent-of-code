// Advent of Code 2016, day 18: Like a Rogue.
// https://adventofcode.com/2016/day/18
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 18, part1, part2)
}

// rowsSmall and rowsLarge are the numbers of rows in the two parts of the
// puzzle. The worked example in the puzzle text uses 10 rows instead, where
// ".^^.^.^^^^" gives 38 safe tiles. With 40 or 400000 rows the example gives
// other counts.
const (
	rowsSmall = 40
	rowsLarge = 400000
)

// part1 counts the safe tiles in the first 40 rows.
func part1(in string) any {
	return countSafe(strings.TrimSpace(in), rowsSmall)
}

// part2 counts the safe tiles in the first 400000 rows.
func part2(in string) any {
	return countSafe(strings.TrimSpace(in), rowsLarge)
}

// countSafe builds the rows one after another and counts the safe tiles.
//
// The four trap rules in the puzzle come down to one test: a tile is a trap
// when exactly one of the tiles above-left and above-right is a trap. So a
// trap is left XOR right, and the tile directly above does not matter.
//
// A row is a slice of 0 (safe) and 1 (trap), with a safe tile at each end.
// These two extra tiles stand for the walls, so the edges need no special case.
func countSafe(first string, rows int) int {
	width := len(first)
	row := make([]byte, width+2)
	next := make([]byte, width+2)

	for i := range width {
		if first[i] == '^' {
			row[i+1] = 1
		}
	}

	safe := 0

	for range rows {
		traps := 0
		for i := 1; i <= width; i++ {
			traps += int(row[i])
			next[i] = row[i-1] ^ row[i+1]
		}

		safe += width - traps
		row, next = next, row
	}

	return safe
}
