// Advent of Code 2023, day 13: Point of Incidence.
// https://adventofcode.com/2023/day/13
package main

import (
	"aoc/lib/input"
	"aoc/lib/strx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 13, part1, part2)
}

// part1 finds the mirror line in each pattern and adds up the summaries.
func part1(in string) any {
	return summarize(in, 0)
}

// part2 finds the new mirror line in each pattern after the smudge is fixed.
//
// The smudge is one wrong cell. So the new mirror line is the one where the
// two sides differ in exactly one cell. This also makes sure that the old
// line, where the sides differ in zero cells, is not found again.
func part2(in string) any {
	return summarize(in, 1)
}

// summarize finds the mirror line in each pattern where the two sides differ
// in exactly smudges cells. It adds the number of columns to the left of each
// vertical line, and 100 times the number of rows above each horizontal line.
func summarize(in string, smudges int) int {
	total := 0

	for _, rows := range input.Blocks(in) {
		if n := mirrorRow(rows, smudges); n > 0 {
			total += 100 * n
			continue
		}

		total += mirrorRow(transpose(rows), smudges)
	}

	return total
}

// mirrorRow finds a horizontal mirror line where the rows on the two sides
// differ in exactly smudges cells. It returns the number of rows above the
// line, or 0 when there is no such line.
func mirrorRow(rows []string, smudges int) int {
	for n := 1; n < len(rows); n++ {
		diff := 0

		// Compare each row above the line with its mirror row below the line.
		// Stop at the edge that is nearer, and stop early when there are too
		// many differences.
		for a, b := n-1, n; a >= 0 && b < len(rows) && diff <= smudges; a, b = a-1, b+1 {
			diff += strx.Hamming(rows[a], rows[b])
		}

		if diff == smudges {
			return n
		}
	}

	return 0
}

// transpose turns the columns of a pattern into rows, so that mirrorRow can
// also find vertical lines.
func transpose(rows []string) []string {
	cols := make([]string, len(rows[0]))

	for x := range cols {
		col := make([]byte, len(rows))
		for y, row := range rows {
			col[y] = row[x]
		}

		cols[x] = string(col)
	}

	return cols
}
