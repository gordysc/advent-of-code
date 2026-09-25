// Advent of Code 2022, day 18: Boiling Boulders.
// https://adventofcode.com/2022/day/18
package main

import (
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 18, part1, part2)
}

// cube is the position of one 1x1x1 cube in 3D space.
type cube struct {
	x, y, z int
}

// sides are the six steps from a cube to the cubes that share a face with it.
var sides = []cube{
	{1, 0, 0}, {-1, 0, 0},
	{0, 1, 0}, {0, -1, 0},
	{0, 0, 1}, {0, 0, -1},
}

// add returns the cube that is one step d away from c.
func (c cube) add(d cube) cube {
	return cube{c.x + d.x, c.y + d.y, c.z + d.z}
}

// part1 counts the faces of the lava cubes that do not touch another lava cube.
func part1(in string) any {
	lava := parse(in)
	area := 0

	for c := range lava.All() {
		for _, d := range sides {
			if !lava.Has(c.add(d)) {
				area++
			}
		}
	}

	return area
}

// part2 counts only the faces that the outside air (steam) can touch.
//
// A box one unit larger than the droplet on each side surrounds all of the
// outside air. A flood fill from one corner of the box finds all the air cubes
// that are outside. Air pockets inside the droplet are not reached. Each face
// of a lava cube that is next to a reached air cube is on the outside.
func part2(in string) any {
	lava := parse(in)
	lo, hi := bounds(lava)

	// inBox tells if a cube is in the box around the droplet.
	inBox := func(c cube) bool {
		return c.x >= lo.x && c.x <= hi.x &&
			c.y >= lo.y && c.y <= hi.y &&
			c.z >= lo.z && c.z <= hi.z
	}

	outside := search.Flood(lo, func(c cube) []cube {
		var next []cube

		for _, d := range sides {
			n := c.add(d)
			if inBox(n) && !lava.Has(n) {
				next = append(next, n)
			}
		}

		return next
	})

	area := 0

	for c := range lava.All() {
		for _, d := range sides {
			// Flood returns a map from cube to distance. We only need to
			// know if the cube is a key, so we ignore the distance.
			if _, ok := outside[c.add(d)]; ok {
				area++
			}
		}
	}

	return area
}

// bounds returns the two corners of a box that holds all the lava cubes, with
// one unit of free space on each side.
func bounds(lava set.Set[cube]) (lo, hi cube) {
	first := true

	for c := range lava.All() {
		if first {
			lo, hi = c, c
			first = false
		}

		lo = cube{min(lo.x, c.x), min(lo.y, c.y), min(lo.z, c.z)}
		hi = cube{max(hi.x, c.x), max(hi.y, c.y), max(hi.z, c.z)}
	}

	lo = lo.add(cube{-1, -1, -1})
	hi = hi.add(cube{1, 1, 1})

	return lo, hi
}

// parse reads one "x,y,z" cube on each line into a set.
func parse(in string) set.Set[cube] {
	lava := set.New[cube]()

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		lava.Add(cube{n[0], n[1], n[2]})
	}

	return lava
}
