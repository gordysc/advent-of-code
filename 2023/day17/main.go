// Advent of Code 2023, day 17: Clumsy Crucible.
// https://adventofcode.com/2023/day/17
package main

import (
	"aoc/lib/grid"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 17, part1, part2)
}

// part1 gives the least heat loss from the top-left block to the
// bottom-right block for a crucible that moves 1 to 3 blocks before it turns.
func part1(in string) any {
	return leastHeatLoss(parse(in), 1, 3)
}

// part2 gives the least heat loss for an ultra crucible. It moves 4 to 10
// blocks before it turns, and it must move at least 4 blocks before it stops
// at the end.
//
// example.txt gives 102 for part 1 and 94 for part 2. The puzzle also has a
// second, smaller example for part 2 (a row of 1s above four rows of 9s that
// end in 1). It gives 71. To test it, put it in a file and use -input.
func part2(in string) any {
	return leastHeatLoss(parse(in), 4, 10)
}

// parse reads the heat loss digit of each city block.
func parse(in string) grid.Grid[int] {
	return grid.Map(in, func(r rune) int { return int(r - '0') })
}

// Axes of the last straight run of the crucible.
const (
	horizontal = iota
	vertical
	none // none marks the start, before the crucible moves.
)

// state is a crucible that has just stopped at p after a straight run along
// axis.
//
// Each move in the search is one full straight run followed by a turn. So
// the state does not need the length of the current run: the next run
// always goes along the other axis. This keeps the number of states small.
type state struct {
	p    grid.Point
	axis int
}

// leastHeatLoss uses Dijkstra's algorithm to find the least heat loss to the
// bottom-right block. Each straight run is minRun to maxRun blocks long.
func leastHeatLoss(g grid.Grid[int], minRun, maxRun int) any {
	end := grid.P(g.W-1, g.H-1)

	next := func(s state) []search.Edge[state] {
		var edges []search.Edge[state]

		for _, d := range grid.Dirs4 {
			axis := horizontal
			if d.X == 0 {
				axis = vertical
			}

			if axis == s.axis {
				continue
			}

			// Walk one block at a time and add the heat loss of each block.
			// Only the stops that are minRun or more blocks away are moves.
			cost := 0
			p := s.p

			for run := 1; run <= maxRun; run++ {
				p = p.Add(d)
				if !g.InBounds(p) {
					break
				}

				cost += g.At(p)

				if run >= minRun {
					edges = append(edges, search.Edge[state]{To: state{p, axis}, Cost: cost})
				}
			}
		}

		return edges
	}

	goal := func(s state) bool { return s.p == end }

	loss, ok := search.Dijkstra(state{grid.P(0, 0), none}, next, goal)
	if !ok {
		return nil
	}

	return loss
}
