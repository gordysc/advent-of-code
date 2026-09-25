// Advent of Code 2020, day 3: Toboggan Trajectory.
// https://adventofcode.com/2020/day/3
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 3, part1, part2)
}

// part1 counts the trees on the slope right 3, down 1.
func part1(in string) any {
	return trees(input.Lines(in), 3, 1)
}

// part2 counts the trees on five slopes and multiplies the counts.
func part2(in string) any {
	rows := input.Lines(in)
	slopes := [][2]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}}

	product := 1
	for _, s := range slopes {
		product *= trees(rows, s[0], s[1])
	}

	return product
}

// trees counts the '#' cells hit when moving right dx and down dy from the
// top-left corner until the bottom. The pattern repeats to the right, so the
// column wraps modulo the row width.
func trees(rows []string, dx, dy int) int {
	count := 0

	x := 0
	for y := 0; y < len(rows); y += dy {
		row := rows[y]
		if row[x%len(row)] == '#' {
			count++
		}

		x += dx
	}

	return count
}
