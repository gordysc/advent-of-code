// Advent of Code 2023, day 21: Step Counter.
// https://adventofcode.com/2023/day/21
package main

import (
	"aoc/lib/grid"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 21, part1, part2)
}

// part2Steps is the number of steps the elf takes in part 2.
const part2Steps = 26501365

// part1 counts the garden plots the elf can be on after exactly 64 steps.
//
// The example asks for 6 steps (answer 16) and is only 11 plots wide, so a
// grid less than 20 plots wide selects 6 steps.
//
// The elf can step back and forth between two plots. So a plot that the elf
// reaches in d steps is also reachable in d+2, d+4, ... steps. The plot
// counts when d is not more than the step count and has the same parity.
//
// Part 1 uses one copy of the map, but the search uses the repeated map of
// part 2. In the real input the start is 65 plots from each edge, so 64
// steps cannot leave the map and the answer is the same.
func part1(in string) any {
	g := grid.Parse(in)

	steps := 64
	if g.W < 20 {
		steps = 6
	}

	return reachable(g, steps)
}

// part2 counts the plots the elf can reach in 26501365 steps on a map that
// repeats forever in all directions.
//
// This relies on the structure of the real input. The map is 131 plots
// square, the start is at the centre, and the row and the column through the
// start have no rocks. 26501365 is 65 + 202300*131, so the elf walks to the
// edge of the start map and then exactly 202300 more maps in each direction.
//
// Let f(n) be the plot count after 65 + 131*n steps. Each time n grows by 1,
// the diamond of reached maps grows by one more ring of maps, and the number
// of maps in a ring grows linearly. So f is a quadratic in n. We count f(0),
// f(1) and f(2) with a search on the repeated map and extend the quadratic
// to n = 202300.
//
// The example does not have a clear row and column through the start, so
// the result is nil for it.
func part2(in string) any {
	g := grid.Parse(in)
	start, _ := g.FindByte('S')

	if !hasCentreCross(g, start) || part2Steps%g.W != g.W/2 {
		return nil
	}

	return extrapolate(g, part2Steps)
}

// hasCentreCross tells if the map is square, the start is at the centre, and
// the row and column through the start have no rocks.
func hasCentreCross(g grid.Grid[byte], start grid.Point) bool {
	if g.W != g.H || g.W%2 == 0 || start != grid.P(g.W/2, g.H/2) {
		return false
	}

	for i := range g.W {
		if g.At(grid.P(i, start.Y)) == '#' || g.At(grid.P(start.X, i)) == '#' {
			return false
		}
	}

	return true
}

// extrapolate counts the plots reachable in steps steps on the repeated map
// from three sample counts and the quadratic through them. steps must be
// g.W/2 plus a multiple of g.W.
func extrapolate(g grid.Grid[byte], steps int) int {
	size := g.W
	rest := steps % size
	n := steps / size

	f0 := reachable(g, rest)
	f1 := reachable(g, rest+size)
	f2 := reachable(g, rest+2*size)

	// Newton's forward differences: f(n) = f0 + n*d1 + n(n-1)/2 * d2.
	d1 := f1 - f0
	d2 := f2 - 2*f1 + f0

	return f0 + n*d1 + n*(n-1)/2*d2
}

// reachable counts the plots the elf can be on after exactly steps steps.
// The map repeats forever, so the search uses as many copies of the map as
// steps can reach.
func reachable(g grid.Grid[byte], steps int) int {
	start, _ := g.FindByte('S')

	// The search can go at most steps plots from the start. Pick enough
	// copies of the map on each side so it cannot leave the tiled area.
	copies := mathx.DivCeil(steps, min(g.W, g.H))
	w, h := g.W*(2*copies+1), g.H*(2*copies+1)
	origin := grid.P(start.X+copies*g.W, start.Y+copies*g.H)

	// dist holds the step count to each plot of the tiled area, or -1 when
	// the search has not reached it. A flat slice is much faster than a map.
	dist := make([]int32, w*h)

	for i := range dist {
		dist[i] = -1
	}

	dist[origin.Y*w+origin.X] = 0
	queue := []grid.Point{origin}
	count := 0

	// A breadth-first search reaches each plot by its shortest path. A
	// slice with a moving head is a first-in, first-out queue.
	for head := 0; head < len(queue); head++ {
		p := queue[head]
		d := int(dist[p.Y*w+p.X])

		if d%2 == steps%2 {
			count++
		}

		if d == steps {
			continue
		}

		for _, q := range p.Neighbors4() {
			if q.X < 0 || q.Y < 0 || q.X >= w || q.Y >= h || dist[q.Y*w+q.X] >= 0 {
				continue
			}

			if g.At(grid.P(q.X%g.W, q.Y%g.H)) == '#' {
				continue
			}

			dist[q.Y*w+q.X] = int32(d + 1)
			queue = append(queue, q)
		}
	}

	return count
}
