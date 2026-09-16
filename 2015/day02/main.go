// Advent of Code 2015, day 2: I Was Told There Would Be No Math.
// https://adventofcode.com/2015/day/2
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 2, part1, part2)
}

// part1 sums the wrapping paper for every present: the surface area of the box
// plus the area of its smallest face as slack.
func part1(in string) any {
	total := 0

	for _, dims := range boxes(in) {
		l, w, h := dims[0], dims[1], dims[2]

		// The dimensions are sorted, so l*w is the smallest face.
		total += 2*(l*w+w*h+h*l) + l*w
	}

	return total
}

// part2 sums the ribbon for every present: the perimeter of the smallest face
// to wrap around the box, plus the volume as the length of the bow.
func part2(in string) any {
	total := 0

	for _, dims := range boxes(in) {
		l, w, h := dims[0], dims[1], dims[2]

		total += 2*(l+w) + l*w*h
	}

	return total
}

// boxes parses every "LxWxH" line into its three dimensions, sorted from
// smallest to largest so both parts can pick the smallest face directly.
func boxes(in string) [][]int {
	var out [][]int

	for _, line := range input.Lines(in) {
		dims := input.Ints(line)
		slices.Sort(dims)

		out = append(out, dims)
	}

	return out
}
