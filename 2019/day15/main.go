// Advent of Code 2019, day 15: Oxygen System.
// https://adventofcode.com/2019/day/15
//
// The puzzle text has no example program, only pictures of the droid's map.
// example.txt holds a small handmade Intcode droid for the map from the part
// 2 example:
//
//	 ##
//	#..##
//	#.#..#
//	#.O.#
//	 ###
//
// The program keeps the droid's position and a table of the map cells in its
// memory. For each command it looks up the next cell, replies with the
// cell's status code, and moves when the cell is open. The droid starts in
// the top-left open cell, so part 1 gives 3 and part 2 gives 4.
package main

import (
	"aoc/lib/grid"
	"aoc/lib/intcode"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 15, part1, part2)
}

// The status codes the droid replies with after each movement command.
const (
	hitWall   = 0
	moved     = 1
	gotOxygen = 2
)

// commands pairs each movement command with the direction it moves the
// droid: 1 north, 2 south, 3 west and 4 east.
var commands = []struct {
	code int
	dir  grid.Point
}{
	{1, grid.Up},
	{2, grid.Down},
	{3, grid.Left},
	{4, grid.Right},
}

// area is the map the droid finds. The droid starts at (0, 0).
type area struct {
	open   map[grid.Point]bool // every cell the droid can stand on
	oxygen grid.Point          // where the oxygen system is
	found  bool                // whether the droid found the oxygen system
}

// part1 finds the fewest movement commands that take the droid from its
// start to the oxygen system.
func part1(in string) any {
	a := explore(intcode.Parse(in))
	if !a.found {
		return nil
	}

	steps, ok := search.BFS(grid.Point{}, a.next, func(p grid.Point) bool {
		return p == a.oxygen
	})
	if !ok {
		return nil
	}

	return steps
}

// part2 finds how many minutes the oxygen takes to fill every open cell.
func part2(in string) any {
	a := explore(intcode.Parse(in))
	if !a.found {
		return nil
	}

	return fill(a)
}

// fill returns the minutes the oxygen takes to spread from the oxygen system
// to every open cell. Each minute it moves one cell further, so the answer
// is the distance to the cell that is furthest away.
func fill(a area) int {
	minutes := 0
	for _, d := range search.Flood(a.oxygen, a.next) {
		minutes = max(minutes, d)
	}

	return minutes
}

// next returns the open cells next to p, for the searches.
func (a area) next(p grid.Point) []grid.Point {
	var out []grid.Point
	for _, n := range p.Neighbors4() {
		if a.open[n] {
			out = append(out, n)
		}
	}

	return out
}

// explore drives the droid to every cell it can reach and returns the map.
//
// It is a breadth-first search in which each queued cell keeps its own copy
// of the droid, as it is after it walked there. To look at a neighbour, the
// search clones that droid and sends one command. So the droid never has to
// walk back, and the other copies are not changed.
func explore(program []int) area {
	type state struct {
		pos   grid.Point
		droid *intcode.Machine
	}

	start := grid.Point{}
	a := area{open: map[grid.Point]bool{start: true}}

	// seen holds the walls as well as the open cells, so that no cell is
	// tested twice.
	seen := map[grid.Point]bool{start: true}
	queue := []state{{start, intcode.New(program)}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, c := range commands {
			pos := cur.pos.Add(c.dir)
			if seen[pos] {
				continue
			}
			seen[pos] = true

			droid := cur.droid.Clone()
			droid.Send(c.code)
			out := droid.Run()
			if len(out) == 0 || out[0] == hitWall {
				continue
			}

			a.open[pos] = true
			if out[0] == gotOxygen {
				a.oxygen, a.found = pos, true
			}

			queue = append(queue, state{pos, droid})
		}
	}

	return a
}
