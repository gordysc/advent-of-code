// Advent of Code 2017, day 11: Hex Ed.
// https://adventofcode.com/2017/day/11
package main

import (
	"strings"

	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 11, part1, part2)
}

// hex is a position on the hex grid in cube coordinates. The three axes always
// sum to zero, and each step moves one unit along two of them.
type hex struct {
	x, y, z int
}

// steps maps each direction in the input to its move in cube coordinates.
var steps = map[string]hex{
	"n":  {0, 1, -1},
	"s":  {0, -1, 1},
	"ne": {1, 0, -1},
	"sw": {-1, 0, 1},
	"se": {1, -1, 0},
	"nw": {-1, 1, 0},
}

// part1 finds how many steps the child process ends up from the start.
func part1(in string) any {
	final, _ := walk(in)

	return final
}

// part2 finds the furthest the child process ever got from the start.
func part2(in string) any {
	_, furthest := walk(in)

	return furthest
}

// walk follows the path and returns the final distance from the start and the
// largest distance seen along the way.
func walk(in string) (final, furthest int) {
	var pos hex

	for _, dir := range strings.Split(strings.TrimSpace(in), ",") {
		step := steps[dir]
		pos = hex{pos.x + step.x, pos.y + step.y, pos.z + step.z}
		furthest = max(furthest, distance(pos))
	}

	return distance(pos), furthest
}

// distance is the smallest number of steps from the origin to p. Each step
// changes two axes by one, so it is half the sum of the absolute axes.
func distance(p hex) int {
	return (mathx.Abs(p.x) + mathx.Abs(p.y) + mathx.Abs(p.z)) / 2
}
