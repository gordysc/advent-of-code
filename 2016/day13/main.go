// Advent of Code 2016, day 13: A Maze of Twisty Little Cubicles.
// https://adventofcode.com/2016/day/13
package main

import (
	"math/bits"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 13, part1, part2)
}

// start is where the walk begins.
var start = grid.P(1, 1)

// target is the location part 1 asks about. The worked example in the puzzle
// text asks about 7,4 instead, so part 1 on the example does not give the 11
// steps from the puzzle text. With the example's favorite number, walls close
// in the area around the start and 31,39 is outside it, so part 1 has no
// answer for the example.
var target = grid.P(31, 39)

// maxSteps is the step limit for part 2.
const maxSteps = 50

// part1 finds the fewest steps from the start to the target.
func part1(in string) any {
	fav := input.Int(in)

	steps, ok := search.BFS(start, moves(fav), func(p grid.Point) bool { return p == target })
	if !ok {
		return nil
	}

	return steps
}

// part2 counts the locations that can be reached in at most 50 steps.
//
// The maze has no end, so a plain flood fill would never stop. A location
// that is more than 50 steps away in a straight line is also more than 50
// steps away through the maze. So the flood never goes past that distance,
// and the distances it finds inside the limit are still exact.
func part2(in string) any {
	fav := input.Int(in)
	step := moves(fav)

	near := func(p grid.Point) []grid.Point {
		var out []grid.Point
		for _, n := range step(p) {
			if n.Manhattan(start) <= maxSteps {
				out = append(out, n)
			}
		}

		return out
	}

	count := 0
	for _, d := range search.Flood(start, near) {
		if d <= maxSteps {
			count++
		}
	}

	return count
}

// moves returns the function that lists the open locations next to a
// location, for the maze made from the favorite number.
func moves(fav int) func(grid.Point) []grid.Point {
	return func(p grid.Point) []grid.Point {
		var out []grid.Point
		for _, n := range p.Neighbors4() {
			if open(n, fav) {
				out = append(out, n)
			}
		}

		return out
	}
}

// open reports whether a location is open space. Negative coordinates are
// outside the building. For the others, the puzzle's formula gives a number,
// and the location is open when that number has an even count of 1 bits.
func open(p grid.Point, fav int) bool {
	if p.X < 0 || p.Y < 0 {
		return false
	}

	x, y := p.X, p.Y
	n := x*x + 3*x + 2*x*y + y + y*y + fav

	return bits.OnesCount(uint(n))%2 == 0
}
