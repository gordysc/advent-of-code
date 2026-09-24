// Advent of Code 2016, day 1: No Time for a Taxicab.
// https://adventofcode.com/2016/day/1
package main

import (
	"strconv"
	"strings"

	"aoc/lib/grid"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 1, part1, part2)
}

// part1 follows every instruction and returns the taxicab distance from the
// start to the final position.
func part1(in string) any {
	pos := grid.Point{}

	for _, p := range walk(in) {
		pos = p
	}

	return pos.Manhattan(grid.Point{})
}

// part2 returns the taxicab distance to the first position visited twice.
// Every block along each leg counts as visited, not only the corners.
func part2(in string) any {
	seen := set.Of(grid.Point{})

	for _, p := range walk(in) {
		if seen.Has(p) {
			return p.Manhattan(grid.Point{})
		}

		seen.Add(p)
	}

	return nil
}

// walk follows the comma-separated instructions from the origin, facing Up,
// and returns every block stepped on in order. The origin is not included.
func walk(in string) []grid.Point {
	pos, dir := grid.Point{}, grid.Up
	var path []grid.Point

	for _, ins := range strings.Split(strings.TrimSpace(in), ", ") {
		if ins[0] == 'L' {
			dir = dir.TurnLeft()
		} else {
			dir = dir.TurnRight()
		}

		blocks, _ := strconv.Atoi(ins[1:])

		for i := 0; i < blocks; i++ {
			pos = pos.Add(dir)
			path = append(path, pos)
		}
	}

	return path
}
