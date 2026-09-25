// Advent of Code 2022, day 24: Blizzard Basin.
// https://adventofcode.com/2022/day/24
//
// example.txt holds the second, larger example from the puzzle text. It
// gives 18 for part 1 and 54 for part 2.
package main

import (
	"aoc/lib/grid"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 24, part1, part2)
}

// part1 finds the fewest minutes to go from the start to the goal.
func part1(in string) any {
	v := parseValley(in)

	return v.travel(v.start, v.goal, 0)
}

// part2 finds the fewest minutes to go to the goal, back to the start, and
// to the goal again.
//
// It is always best to arrive at each end as early as possible: an early
// arrival can wait at the end tile, where no blizzard goes. So the three
// fastest trips, one after the other, give the fastest full journey.
func part2(in string) any {
	v := parseValley(in)

	t := v.travel(v.start, v.goal, 0)
	t = v.travel(v.goal, v.start, t)
	t = v.travel(v.start, v.goal, t)

	return t
}

// valley holds the map with the blizzards at minute 0.
type valley struct {
	g           grid.Grid[byte]
	start, goal grid.Point

	// w and h are the size of the inside of the valley, without the walls.
	w, h int
}

// parseValley reads the map. The start is the gap in the top wall, and the
// goal is the gap in the bottom wall.
func parseValley(in string) valley {
	g := grid.Parse(in)
	v := valley{g: g, w: g.W - 2, h: g.H - 2}

	for x := range g.W {
		if g.At(grid.P(x, 0)) == '.' {
			v.start = grid.P(x, 0)
		}
		if g.At(grid.P(x, g.H-1)) == '.' {
			v.goal = grid.P(x, g.H-1)
		}
	}

	return v
}

// open reports if the tile p is free of walls and blizzards at minute t.
//
// A blizzard moves in a straight line and wraps around, so it is not
// necessary to move all the blizzards. To find if a blizzard is at p at
// minute t, look back along each of the four lines through p. For example,
// a '>' blizzard is at p when there was a '>' t tiles to the left of p at
// minute 0.
func (v valley) open(p grid.Point, t int) bool {
	c, ok := v.g.Get(p)
	if !ok || c == '#' {
		return false
	}

	// No blizzard goes into the top or bottom row, so the start and goal
	// are always free.
	if p.Y == 0 || p.Y == v.g.H-1 {
		return true
	}

	// x and y are the position inside the walls, from 0.
	x, y := p.X-1, p.Y-1

	switch {
	case v.g.At(grid.P(mathx.Mod(x-t, v.w)+1, p.Y)) == '>':
		return false
	case v.g.At(grid.P(mathx.Mod(x+t, v.w)+1, p.Y)) == '<':
		return false
	case v.g.At(grid.P(p.X, mathx.Mod(y-t, v.h)+1)) == 'v':
		return false
	case v.g.At(grid.P(p.X, mathx.Mod(y+t, v.h)+1)) == '^':
		return false
	}

	return true
}

// travel finds the first minute at which the expedition can be at to, when
// it starts at from at minute t.
//
// It keeps the set of all tiles the expedition can be on at each minute.
// Each minute, a tile is in the next set if it is open and the expedition
// can get to it (or stay on it) from a tile in the current set. This is a
// breadth-first search where time is the distance.
func (v valley) travel(from, to grid.Point, t int) int {
	here := grid.New[bool](v.g.W, v.g.H)
	here.Set(from, true)

	for !here.At(to) {
		t++
		next := grid.New[bool](v.g.W, v.g.H)

		for p, ok := range here.All() {
			if !ok {
				continue
			}

			for _, q := range append(p.Neighbors4(), p) {
				if v.open(q, t) {
					next.Set(q, true)
				}
			}
		}

		here = next
	}

	return t
}
