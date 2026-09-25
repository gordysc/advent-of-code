// Advent of Code 2017, day 22: Sporifica Virus.
// https://adventofcode.com/2017/day/22
package main

import (
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 22, part1, part2)
}

// part1Bursts and part2Bursts are how many bursts of activity each part runs.
// The puzzle text also shows shorter runs of the example: 70 bursts give 41
// infections in part 1, and 100 bursts give 26 in part 2.
const (
	part1Bursts = 10_000
	part2Bursts = 10_000_000
)

// state is the condition of one node. The states are in the order that the
// evolved virus of part 2 moves a node through them.
type state uint8

const (
	clean state = iota
	weakened
	infected
	flagged
)

// cluster is the infinite grid of nodes. It keeps the nodes in one flat slice
// with the point (0, 0) at index offset in both directions, and it grows when
// the carrier walks off the edge. A flat slice is much faster than a map for
// the ten million bursts of part 2.
type cluster struct {
	cells  []state
	size   int // width and height of the stored square
	offset int // added to x and y to find the stored row and column
}

// part1 counts the bursts that infect a node when the virus only toggles
// nodes between clean and infected.
func part1(in string) any {
	return run(in, part1Bursts, 2)
}

// part2 counts the bursts that infect a node when the evolved virus moves
// each node through all four states.
func part2(in string) any {
	return run(in, part2Bursts, 1)
}

// run lets the carrier work for the given number of bursts and counts the
// bursts that leave a node infected.
//
// At each burst the node moves step places forward through the states, so a
// step of 2 swaps clean and infected (part 1) and a step of 1 goes clean,
// weakened, infected, flagged (part 2). The turn rules are the same in both
// parts, because the part 1 virus never makes a weakened or flagged node.
func run(in string, bursts int, step state) int {
	c, pos := parse(in)
	dir := grid.Up
	infections := 0

	for range bursts {
		node := c.at(pos)

		// A weakened node does not change the direction.
		switch *node {
		case clean:
			dir = dir.TurnLeft()
		case infected:
			dir = dir.TurnRight()
		case flagged:
			dir = dir.Reverse()
		}

		*node = (*node + step) % 4
		if *node == infected {
			infections++
		}

		pos = pos.Add(dir)
	}

	return infections
}

// parse reads the map of infected nodes and returns the cluster with the
// carrier's start point, which is the middle of the map.
func parse(in string) (*cluster, grid.Point) {
	lines := input.Lines(in)
	c := &cluster{cells: []state{clean}, size: 1}

	for y, line := range lines {
		for x, ch := range line {
			if ch == '#' {
				*c.at(grid.P(x, y)) = infected
			}
		}
	}

	return c, grid.P(len(lines[0])/2, len(lines)/2)
}

// at returns a pointer to the node at p, so the caller can read and change it
// in one lookup. It grows the cluster first when p is outside it.
func (c *cluster) at(p grid.Point) *state {
	for !c.inside(p) {
		c.grow()
	}

	return &c.cells[(p.Y+c.offset)*c.size+p.X+c.offset]
}

// inside reports whether p is in the stored square.
func (c *cluster) inside(p grid.Point) bool {
	x, y := p.X+c.offset, p.Y+c.offset

	return x >= 0 && x < c.size && y >= 0 && y < c.size
}

// grow triples the width and height of the stored square and puts the old
// square in the middle, so the cluster has room on every side.
func (c *cluster) grow() {
	size := c.size * 3
	cells := make([]state, size*size)

	for y := range c.size {
		copy(cells[(y+c.size)*size+c.size:], c.cells[y*c.size:(y+1)*c.size])
	}

	c.cells = cells
	c.offset += c.size
	c.size = size
}
