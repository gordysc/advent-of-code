// Advent of Code 2021, day 25: Sea Cucumber.
// https://adventofcode.com/2021/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 25, part1, nil)
}

// part1 finds the first step on which no sea cucumber moves.
//
// In each step, the east-facing herd ('>') moves first, then the
// south-facing herd ('v'). A sea cucumber moves one cell if that cell is
// empty at the start of its herd's move. The edges wrap around.
func part1(in string) any {
	g := grid.Parse(in)

	for step := 1; ; step++ {
		movedEast := moveHerd(&g, '>', grid.Right)
		movedSouth := moveHerd(&g, 'v', grid.Down)

		if !movedEast && !movedSouth {
			return step
		}
	}
}

// moveHerd moves all sea cucumbers of one herd at the same time, and reports
// if any of them moved.
//
// It reads from the old grid and writes to a copy, so a sea cucumber does
// not see the moves of its own herd in the same step.
func moveHerd(g *grid.Grid[byte], herd byte, dir grid.Point) bool {
	next := g.Clone()
	moved := false

	for p, c := range g.All() {
		if c != herd {
			continue
		}

		to := p.Add(dir)
		to = grid.P(to.X%g.W, to.Y%g.H)

		if g.At(to) == '.' {
			next.Set(to, herd)
			next.Set(p, '.')
			moved = true
		}
	}

	*g = next

	return moved
}
