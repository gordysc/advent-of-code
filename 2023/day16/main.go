// Advent of Code 2023, day 16: The Floor Will Be Lava.
// https://adventofcode.com/2023/day/16
package main

import (
	"aoc/lib/ds"
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 16, part1, part2)
}

// Direction indexes into grid.Dirs4, which lists the directions clockwise.
const (
	up = iota
	right
	down
	left
)

// part1 counts the energized tiles when the beam enters at the top-left
// corner and moves right.
func part1(in string) any {
	c := newContraption(grid.Parse(in))

	return c.energize(beam{grid.P(0, 0), right})
}

// part2 tries each tile on the edge as the entry, with the beam moving away
// from that edge, and gives the largest number of energized tiles.
func part2(in string) any {
	c := newContraption(grid.Parse(in))
	best := 0

	for x := range c.g.W {
		best = max(best, c.energize(beam{grid.P(x, 0), down}))
		best = max(best, c.energize(beam{grid.P(x, c.g.H-1), up}))
	}

	for y := range c.g.H {
		best = max(best, c.energize(beam{grid.P(0, y), right}))
		best = max(best, c.energize(beam{grid.P(c.g.W-1, y), left}))
	}

	return best
}

// beam is the front of a light beam: the tile it is on and the direction it
// moves in.
type beam struct {
	p   grid.Point
	dir int
}

// contraption is the grid of mirrors and splitters. It also keeps a visited
// buffer, so that part 2 can use the same memory for each entry.
type contraption struct {
	g grid.Grid[byte]

	// visited has one flag for each tile and direction, at index
	// (y*W + x)*4 + dir. A flat bool slice is much faster than a map.
	visited []bool
}

// newContraption makes a contraption for the grid g.
func newContraption(g grid.Grid[byte]) *contraption {
	return &contraption{g: g, visited: make([]bool, g.W*g.H*4)}
}

// energize follows the beam that starts at start and counts the tiles that
// the light goes through.
//
// A beam that comes back to a tile in a direction it had before only repeats
// earlier work, so we stop it. This also stops beams that go in loops.
func (c *contraption) energize(start beam) int {
	clear(c.visited)

	var stack ds.Stack[beam]
	stack.Push(start)

	for !stack.Empty() {
		b := stack.Pop()

		if !c.g.InBounds(b.p) {
			continue
		}

		i := (b.p.Y*c.g.W+b.p.X)*4 + b.dir
		if c.visited[i] {
			continue
		}
		c.visited[i] = true

		for _, dir := range next(c.g.At(b.p), b.dir) {
			stack.Push(beam{b.p.Add(grid.Dirs4[dir]), dir})
		}
	}

	count := 0

	for tile := 0; tile < len(c.visited); tile += 4 {
		v := c.visited[tile : tile+4]
		if v[up] || v[right] || v[down] || v[left] {
			count++
		}
	}

	return count
}

// next gives the directions in which the light leaves a tile, when it comes
// in with direction dir.
func next(tile byte, dir int) []int {
	switch tile {
	case '/':
		// Up and right swap, and down and left swap. With the clockwise
		// numbering 0..3, that is the same as flipping the lowest bit.
		return []int{dir ^ 1}
	case '\\':
		// Up and left swap, and right and down swap: 0<->3 and 1<->2.
		return []int{3 - dir}
	case '|':
		if dir == left || dir == right {
			return []int{up, down}
		}
	case '-':
		if dir == up || dir == down {
			return []int{left, right}
		}
	}

	return []int{dir}
}
