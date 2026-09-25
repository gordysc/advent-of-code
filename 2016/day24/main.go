// Advent of Code 2016, day 24: Air Duct Spelunking.
// https://adventofcode.com/2016/day/24
package main

import (
	"math"

	"aoc/lib/grid"
	"aoc/lib/search"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 24, part1, part2)
}

// wall is the character for a spot the robot cannot enter.
const wall = '#'

// part1 finds the fewest steps to start at 0 and visit every number.
func part1(in string) any {
	return shortest(distances(in), false)
}

// part2 finds the fewest steps to visit every number and then return to 0.
func part2(in string) any {
	return shortest(distances(in), true)
}

// distances returns the fewest steps between every pair of numbers in the
// maze. dist[i][j] is the distance from number i to number j. It floods the
// maze once from each number.
func distances(in string) [][]int {
	g := grid.Parse(in)

	// The numbers run from 0 up without gaps, so 0 to len(spots)-1 are all of them.
	spots := map[int]grid.Point{}
	for p, c := range g.All() {
		if c >= '0' && c <= '9' {
			spots[int(c-'0')] = p
		}
	}

	// next returns the open spots next to p. The maze has walls all around
	// it, so every neighbour is inside the grid.
	next := func(p grid.Point) []grid.Point {
		return slicesx.Filter(p.Neighbors4(), func(q grid.Point) bool { return g.At(q) != wall })
	}

	dist := make([][]int, len(spots))
	for i := range dist {
		reach := search.Flood(spots[i], next)

		dist[i] = make([]int, len(spots))
		for j := range dist[i] {
			dist[i][j] = reach[spots[j]]
		}
	}

	return dist
}

// shortest tries every order to visit the numbers after 0 and returns the
// length of the best route. With back set, the route also returns to 0.
// The mazes have at most eight numbers, so there are few enough orders.
func shortest(dist [][]int, back bool) int {
	rest := make([]int, len(dist)-1)
	for i := range rest {
		rest[i] = i + 1
	}

	best := math.MaxInt
	for order := range slicesx.Permutations(rest) {
		steps, at := 0, 0
		for _, n := range order {
			steps += dist[at][n]
			at = n
		}

		if back {
			steps += dist[at][0]
		}

		best = min(best, steps)
	}

	return best
}
