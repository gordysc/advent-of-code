// Advent of Code 2024, day 12: Garden Groups.
// https://adventofcode.com/2024/day/12
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 12, part1, part2)
}

// part1 adds the fence price of every region, where the price is the area
// times the perimeter.
//
// The puzzle has several examples. The file example.txt holds the largest
// one. On it, part 1 gives 1930 and part 2 gives 1206.
func part1(in string) any {
	total := 0
	for _, r := range regions(grid.Parse(in)) {
		total += r.area * r.perimeter
	}

	return total
}

// part2 adds the fence price of every region, where the price is the area
// times the number of sides.
func part2(in string) any {
	total := 0
	for _, r := range regions(grid.Parse(in)) {
		total += r.area * r.corners
	}

	return total
}

// region holds the measurements of one connected group of the same plant.
type region struct {
	area      int
	perimeter int
	// corners is equal to the number of sides. A closed outline turns once at
	// the end of each straight side, so it has as many corners as sides.
	corners int
}

// regions finds every region in the garden with a flood fill, and measures
// each one while it fills.
func regions(g grid.Grid[byte]) []region {
	seen := grid.New[bool](g.W, g.H)
	var out []region

	for start := range g.Points() {
		if seen.At(start) {
			continue
		}

		plant := g.At(start)
		seen.Set(start, true)
		stack := []grid.Point{start}
		var r region

		for len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			r.area++
			r.corners += corners(g, p)

			for _, d := range grid.Dirs4 {
				q := p.Add(d)

				// A neighbor with a different plant, or no neighbor at the
				// edge of the map, needs one unit of fence.
				if !same(g, q, plant) {
					r.perimeter++
					continue
				}

				if !seen.At(q) {
					seen.Set(q, true)
					stack = append(stack, q)
				}
			}
		}

		out = append(out, r)
	}

	return out
}

// corners counts the outline corners at the four corners of the plot p.
//
// Look at two orthogonal neighbors a and b, and the diagonal c between them.
// There is an outside corner when a and b are both in another region. There
// is an inside corner when a and b are both in the region, but c is not.
func corners(g grid.Grid[byte], p grid.Point) int {
	plant := g.At(p)
	n := 0

	for _, d := range grid.Dirs4 {
		e := d.TurnRight()

		a := same(g, p.Add(d), plant)
		b := same(g, p.Add(e), plant)
		c := same(g, p.Add(d).Add(e), plant)

		if !a && !b {
			n++
		}

		if a && b && !c {
			n++
		}
	}

	return n
}

// same reports whether p is inside the map and has the given plant.
func same(g grid.Grid[byte], p grid.Point, plant byte) bool {
	v, ok := g.Get(p)

	return ok && v == plant
}
