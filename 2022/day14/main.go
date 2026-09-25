// Advent of Code 2022, day 14: Regolith Reservoir.
// https://adventofcode.com/2022/day/14
package main

import (
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 14, part1, part2)
}

// part1 counts the units of sand that come to rest before sand starts to fall
// into the abyss below the lowest rock.
func part1(in string) any {
	return pour(parse(in), false)
}

// part2 counts the units of sand that come to rest before the source is
// blocked. A floor two rows below the lowest rock stops all sand.
func part2(in string) any {
	return pour(parse(in), true)
}

// source is where the sand comes into the cave.
var source = grid.P(500, 0)

// cave holds the blocked cells. Sand moves at most one column to the side per
// row, so it always stays in a triangle below the source. The grid covers only
// that triangle, down to the row above the floor. Column 0 of the grid is x =
// 500 - floor.
type cave struct {
	blocked grid.Grid[bool]
	offset  int // offset is the x value of grid column 0.
	lowest  int // lowest is the y value of the lowest rock.
	floor   int // floor is the y value of the floor.
}

// parse reads the rock paths and draws them into a new cave.
func parse(in string) cave {
	var paths [][]grid.Point
	lowest := 0

	for _, line := range input.Lines(in) {
		nums := input.Ints(line)

		var path []grid.Point
		for i := 0; i+1 < len(nums); i += 2 {
			path = append(path, grid.P(nums[i], nums[i+1]))
			lowest = max(lowest, nums[i+1])
		}

		paths = append(paths, path)
	}

	floor := lowest + 2
	c := cave{
		blocked: grid.New[bool](2*floor+1, floor),
		offset:  source.X - floor,
		lowest:  lowest,
		floor:   floor,
	}

	for _, path := range paths {
		for i := 1; i < len(path); i++ {
			c.drawLine(path[i-1], path[i])
		}
	}

	return c
}

// drawLine marks each cell on the straight line from a to b as rock. Rock
// outside the triangle of the grid cannot touch sand, so it is ignored.
func (c cave) drawLine(a, b grid.Point) {
	step := grid.P(mathx.Sign(b.X-a.X), mathx.Sign(b.Y-a.Y))

	for p := a; ; p = p.Add(step) {
		if q := c.local(p); c.blocked.InBounds(q) {
			c.blocked.Set(q, true)
		}

		if p == b {
			break
		}
	}
}

// local changes a cave position into a grid position.
func (c cave) local(p grid.Point) grid.Point {
	return grid.P(p.X-c.offset, p.Y)
}

// isBlocked reports whether sand cannot move into p. The floor is always
// blocked.
func (c cave) isBlocked(p grid.Point) bool {
	if p.Y >= c.floor {
		return true
	}

	return c.blocked.At(c.local(p))
}

// pour drops sand until it falls into the abyss (withFloor false) or until
// the source is blocked (withFloor true), and returns the number of units at
// rest.
//
// A new unit follows the same path as the unit before, until the cell where
// that unit came to rest. So we keep the path in a stack. When a unit comes
// to rest, we pop its cell, and the next unit starts from the cell before it.
// This way each cell is visited only a small number of times.
func pour(c cave, withFloor bool) int {
	moves := []grid.Point{grid.Down, grid.DownLeft, grid.DownRight}
	path := []grid.Point{source}
	count := 0

	for len(path) > 0 {
		p := path[len(path)-1]

		// Without the floor, sand below the lowest rock falls forever.
		if !withFloor && p.Y >= c.lowest {
			return count
		}

		moved := false
		for _, m := range moves {
			if q := p.Add(m); !c.isBlocked(q) {
				path = append(path, q)
				moved = true

				break
			}
		}

		if moved {
			continue
		}

		// The unit cannot move, so it comes to rest here.
		c.blocked.Set(c.local(p), true)
		path = path[:len(path)-1]
		count++
	}

	// The stack is empty, so the unit at the source came to rest.
	return count
}
