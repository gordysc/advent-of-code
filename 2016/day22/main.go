// Advent of Code 2016, day 22: Grid Computing.
// https://adventofcode.com/2016/day/22
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 22, part1, part2)
}

// node is the disk use of one storage node, in terabytes.
type node struct {
	size, used int
}

// avail returns the free space on the node.
func (n node) avail() int {
	return n.size - n.used
}

// part1 counts the viable pairs: ordered pairs of different nodes where the
// first holds some data and all of it fits in the free space of the second.
func part1(in string) any {
	nodes := parse(in).Cells
	count := 0

	for i, a := range nodes {
		if a.used == 0 {
			continue
		}

		for j, b := range nodes {
			if i != j && a.used <= b.avail() {
				count++
			}
		}
	}

	return count
}

// part2 finds the fewest moves to bring the data from the top-right node to
// the top-left node.
//
// The puzzle's grids have one empty node, a few huge nodes whose data cannot
// move anywhere (walls), and many ordinary nodes that all fit into the empty
// one. So every move slides data into the empty node, which is the same as
// moving the empty node one step. First the empty node walks to the cell left
// of the goal data, and a final move swaps them. After that, each step left
// for the goal takes 5 moves: 4 to take the empty node around the goal to its
// left side again, and 1 to swap.
func part2(in string) any {
	g := parse(in)

	empty, ok := g.Find(func(n node) bool { return n.used == 0 })
	if !ok {
		return nil
	}

	space := g.At(empty).size
	goal := grid.P(g.W-1, 0)
	target := grid.P(g.W-2, 0)

	// next moves the empty node to a neighbour whose data fits into it. The
	// goal data must not move while the empty node gets into place.
	next := func(p grid.Point) []grid.Point {
		var out []grid.Point
		for _, q := range g.Neighbors4(p) {
			if q != goal && g.At(q).used <= space {
				out = append(out, q)
			}
		}

		return out
	}

	steps, ok := search.BFS(empty, next, func(p grid.Point) bool { return p == target })
	if !ok {
		return nil
	}

	return steps + 1 + 5*target.X
}

// parse reads the df listing into a grid of nodes. The first two lines are
// the command and the column headers, so only lines that name a node count.
// Each node line holds x, y, size, used, avail and use%, in that order.
func parse(in string) grid.Grid[node] {
	var rows [][]int
	w, h := 0, 0

	for _, line := range input.Lines(in) {
		if !strings.HasPrefix(line, "/dev/grid/") {
			continue
		}

		n := input.UInts(line)
		rows = append(rows, n)
		w, h = max(w, n[0]+1), max(h, n[1]+1)
	}

	g := grid.New[node](w, h)
	for _, n := range rows {
		g.Set(grid.P(n[0], n[1]), node{size: n[2], used: n[3]})
	}

	return g
}
