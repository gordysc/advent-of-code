// Advent of Code 2024, day 4: Ceres Search.
// https://adventofcode.com/2024/day/4
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 4, part1, part2)
}

// part1 counts each time the word XMAS occurs in the grid. The word can go in
// any of the eight directions, and the occurrences can overlap.
func part1(in string) any {
	g := grid.Parse(in)
	count := 0

	// Start only at an X, then try each direction from it. This counts each
	// occurrence one time, because an occurrence has one start and one
	// direction.
	for start, c := range g.All() {
		if c != 'X' {
			continue
		}

		for _, dir := range grid.Dirs8 {
			if spells(g, start, dir, "XMAS") {
				count++
			}
		}
	}

	return count
}

// part2 counts each X-MAS: two MAS words that cross on their A, one on each
// diagonal. Each MAS can go in either direction.
func part2(in string) any {
	g := grid.Parse(in)
	count := 0

	for center, c := range g.All() {
		if c != 'A' {
			continue
		}

		// Read each diagonal from its top end, through the A, to its bottom
		// end. The diagonal is a MAS if it spells MAS or SAM.
		topLeft := center.Add(grid.UpLeft)
		topRight := center.Add(grid.UpRight)

		if isMAS(g, topLeft, grid.DownRight) && isMAS(g, topRight, grid.DownLeft) {
			count++
		}
	}

	return count
}

// isMAS reports whether the three cells from start in direction dir spell
// MAS in either direction.
func isMAS(g grid.Grid[byte], start, dir grid.Point) bool {
	return spells(g, start, dir, "MAS") || spells(g, start, dir, "SAM")
}

// spells reports whether the cells from start in direction dir spell word.
// A word that goes off the edge of the grid does not match.
func spells(g grid.Grid[byte], start, dir grid.Point, word string) bool {
	p := start

	for i := 0; i < len(word); i++ {
		// Get returns false for a point outside the grid, so this also
		// checks the bounds.
		c, ok := g.Get(p)
		if !ok || c != word[i] {
			return false
		}

		p = p.Add(dir)
	}

	return true
}
