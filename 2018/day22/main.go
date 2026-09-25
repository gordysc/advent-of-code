// Advent of Code 2018, day 22: Mode Maze.
// https://adventofcode.com/2018/day/22
package main

import (
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 22, part1, part2)
}

// The numbers from the puzzle that turn a position into an erosion level.
const (
	xFactor = 16807
	yFactor = 48271
	modulus = 20183
)

// The region types. A region's type is its erosion level modulo 3.
const (
	rocky = iota
	wet
	narrow
)

// The tools. Their numbers are chosen so that each region type forbids the
// tool with the same number: rocky forbids neither, wet forbids the torch,
// and narrow forbids the climbing gear.
const (
	neither = iota
	torch
	gear
)

// The costs in minutes of one step and of a change of tool.
const (
	moveCost   = 1
	switchCost = 7
)

// cave works out erosion levels as the search asks for them and keeps each
// one, because every level depends on the levels to its left and above it.
type cave struct {
	depth   int
	target  grid.Point
	erosion map[grid.Point]int
}

// state is a position in the cave together with the tool in hand.
type state struct {
	pos  grid.Point
	tool int
}

// part1 adds up the risk levels (the region types) of the rectangle from the
// mouth to the target.
func part1(in string) any {
	c := parse(in)

	risk := 0
	for y := 0; y <= c.target.Y; y++ {
		for x := 0; x <= c.target.X; x++ {
			risk += c.kind(grid.P(x, y))
		}
	}

	return risk
}

// part2 finds the fewest minutes to reach the target with the torch in hand.
// Dijkstra's algorithm searches over (position, tool). The cave is worked out
// on demand, so the search can go as far past the target as it needs.
func part2(in string) any {
	c := parse(in)
	start := state{grid.P(0, 0), torch}
	goal := state{c.target, torch}

	next := func(s state) []search.Edge[state] {
		kind := c.kind(s.pos)

		// The other tool that this region allows.
		edges := []search.Edge[state]{{To: state{s.pos, 3 - kind - s.tool}, Cost: switchCost}}

		for _, n := range s.pos.Neighbors4() {
			if n.X < 0 || n.Y < 0 || c.kind(n) == s.tool {
				continue
			}

			edges = append(edges, search.Edge[state]{To: state{n, s.tool}, Cost: moveCost})
		}

		return edges
	}

	minutes, _ := search.Dijkstra(start, next, func(s state) bool { return s == goal })

	return minutes
}

// parse reads the depth and the target position.
func parse(in string) *cave {
	n := input.Ints(in)

	return &cave{depth: n[0], target: grid.P(n[1], n[2]), erosion: map[grid.Point]int{}}
}

// kind returns the region type at p.
func (c *cave) kind(p grid.Point) int {
	return c.erosionAt(p) % 3
}

// erosionAt returns the erosion level at p, which is its geologic index plus
// the depth, modulo 20183.
func (c *cave) erosionAt(p grid.Point) int {
	if e, ok := c.erosion[p]; ok {
		return e
	}

	var geo int
	switch {
	case p == grid.P(0, 0) || p == c.target:
		geo = 0
	case p.Y == 0:
		geo = p.X * xFactor
	case p.X == 0:
		geo = p.Y * yFactor
	default:
		geo = c.erosionAt(grid.P(p.X-1, p.Y)) * c.erosionAt(grid.P(p.X, p.Y-1))
	}

	e := (geo + c.depth) % modulus
	c.erosion[p] = e

	return e
}
