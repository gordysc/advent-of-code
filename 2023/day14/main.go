// Advent of Code 2023, day 14: Parabolic Reflector Dish.
// https://adventofcode.com/2023/day/14
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 14, part1, part2)
}

// part1 tilts the platform north one time and gives the load on the north
// support beams.
func part1(in string) any {
	g := grid.Parse(in)
	tilt(g, grid.Up)

	return load(g)
}

// part2 gives the load after 1,000,000,000 spin cycles.
//
// It is not possible to do that many cycles. But after some cycles, the
// platform gets back to a state that it had before, and from then on the
// states repeat in a loop. When we find the first repeat, we know the length
// of the loop. Then we can calculate which earlier state is the same as the
// state after the last cycle.
func part2(in string) any {
	const cycles = 1_000_000_000

	g := grid.Parse(in)

	// seen maps each state (the grid cells as a string) to the number of
	// cycles done when the state first occurred. loads[k] is the load after
	// k cycles.
	seen := map[string]int{string(g.Cells): 0}
	loads := []int{load(g)}

	for done := 1; done <= cycles; done++ {
		spin(g)

		// string(g.Cells) copies the bytes, so later changes to the grid do
		// not change the key that is kept in the map.
		key := string(g.Cells)
		if first, ok := seen[key]; ok {
			period := done - first

			return loads[first+(cycles-first)%period]
		}

		seen[key] = done
		loads = append(loads, load(g))
	}

	return load(g)
}

// spin does one spin cycle: it tilts the platform north, west, south and then
// east.
func spin(g grid.Grid[byte]) {
	for _, d := range []grid.Point{grid.Up, grid.Left, grid.Down, grid.Right} {
		tilt(g, d)
	}
}

// tilt rolls all round rocks ('O') as far as they can go in direction d. Cube
// rocks ('#') do not move.
//
// Each lane is a row or column parallel to d. We walk each lane from the edge
// that d points to, and we keep the nearest free cell. A round rock moves to
// that free cell. A cube rock makes the cell after it the new free cell.
func tilt(g grid.Grid[byte], d grid.Point) {
	step := d.Reverse()

	for _, start := range laneStarts(g, d) {
		free := start

		for p := start; g.InBounds(p); p = p.Add(step) {
			switch g.At(p) {
			case '#':
				free = p.Add(step)
			case 'O':
				// Clear the old cell first, because free can be the same cell.
				g.Set(p, '.')
				g.Set(free, 'O')
				free = free.Add(step)
			}
		}
	}
}

// laneStarts gives the first cell of each lane for a tilt in direction d.
// These are the cells on the edge of the grid that d points to.
func laneStarts(g grid.Grid[byte], d grid.Point) []grid.Point {
	var starts []grid.Point

	switch d {
	case grid.Up, grid.Down:
		y := 0
		if d == grid.Down {
			y = g.H - 1
		}

		for x := range g.W {
			starts = append(starts, grid.P(x, y))
		}
	case grid.Left, grid.Right:
		x := 0
		if d == grid.Right {
			x = g.W - 1
		}

		for y := range g.H {
			starts = append(starts, grid.P(x, y))
		}
	}

	return starts
}

// load adds the load of all round rocks. The load of one rock is its number
// of rows from the south edge, where the bottom row counts as 1.
func load(g grid.Grid[byte]) int {
	total := 0

	for p, v := range g.All() {
		if v == 'O' {
			total += g.H - p.Y
		}
	}

	return total
}
