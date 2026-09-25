// Advent of Code 2021, day 15: Chiton.
// https://adventofcode.com/2021/day/15
package main

import (
	"aoc/lib/grid"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 15, part1, part2)
}

// part1 finds the lowest total risk of a path from the top-left corner to the
// bottom-right corner.
func part1(in string) any {
	return lowestRisk(parse(in))
}

// part2 finds the lowest total risk on a map that is five times larger in
// each direction.
//
// The full map is a 5×5 arrangement of copies of the input. Each copy one
// step to the right or down adds 1 to every risk level. A risk above 9 wraps
// back to 1.
func part2(in string) any {
	tile := parse(in)
	full := grid.New[int](tile.W*5, tile.H*5)

	for p := range full.Points() {
		base := tile.At(grid.P(p.X%tile.W, p.Y%tile.H))
		extra := p.X/tile.W + p.Y/tile.H

		// Take 1 before the modulo so the levels wrap from 9 to 1, not to 0.
		full.Set(p, (base+extra-1)%9+1)
	}

	return lowestRisk(full)
}

// parse reads the risk level of each position into a grid.
func parse(in string) grid.Grid[int] {
	return grid.Map(in, func(r rune) int { return int(r - '0') })
}

// lowestRisk runs Dijkstra from the top-left corner to the bottom-right
// corner. Moving into a position costs its risk level. The start position is
// not entered, so its risk does not count.
func lowestRisk(risk grid.Grid[int]) int {
	start := grid.P(0, 0)
	goal := grid.P(risk.W-1, risk.H-1)

	next := func(p grid.Point) []search.Edge[grid.Point] {
		var edges []search.Edge[grid.Point]

		for _, n := range risk.Neighbors4(p) {
			edges = append(edges, search.Edge[grid.Point]{To: n, Cost: risk.At(n)})
		}

		return edges
	}

	cost, _ := search.Dijkstra(start, next, func(p grid.Point) bool { return p == goal })

	return cost
}
