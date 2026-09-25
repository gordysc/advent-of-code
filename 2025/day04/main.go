// Advent of Code 2025, day 4: Printing Department.
// https://adventofcode.com/2025/day/4
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// roll is the grid byte for a roll of paper.
const roll = '@'

// limit is the number of neighbouring rolls that blocks the forklift. A roll
// with fewer neighbouring rolls than this is accessible.
const limit = 4

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 4, part1, part2)
}

// part1 counts the rolls that a forklift can get to now.
func part1(in string) any {
	g := grid.Parse(in)
	count := 0

	for p, b := range g.All() {
		if b == roll && neighbourRolls(g, p) < limit {
			count++
		}
	}

	return count
}

// part2 removes accessible rolls until no more are accessible, and returns the
// number of rolls removed.
//
// When a roll goes, only its eight neighbours can change. So we do not scan
// the full grid again after each removal. We keep a count of the neighbouring
// rolls of each roll, and a queue of rolls to remove. A roll goes into the
// queue once: at the start if its count is already below the limit, or later
// at the moment its count falls from limit to limit-1.
func part2(in string) any {
	g := grid.Parse(in)
	counts := grid.New[int](g.W, g.H)

	var queue []grid.Point

	for p, b := range g.All() {
		if b != roll {
			continue
		}

		counts.Set(p, neighbourRolls(g, p))

		if counts.At(p) < limit {
			queue = append(queue, p)
		}
	}

	removed := 0

	// The loop reads len(queue) again on each pass, so it also visits the
	// rolls that we append while it runs.
	for i := 0; i < len(queue); i++ {
		p := queue[i]
		g.Set(p, '.')
		removed++

		for _, q := range g.Neighbors8(p) {
			if g.At(q) != roll {
				continue
			}

			counts.Set(q, counts.At(q)-1)

			if counts.At(q) == limit-1 {
				queue = append(queue, q)
			}
		}
	}

	return removed
}

// neighbourRolls counts the rolls in the eight cells around p.
func neighbourRolls(g grid.Grid[byte], p grid.Point) int {
	count := 0

	for _, q := range g.Neighbors8(p) {
		if g.At(q) == roll {
			count++
		}
	}

	return count
}
