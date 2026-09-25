// Advent of Code 2016, day 3: Squares With Three Sides.
// https://adventofcode.com/2016/day/3
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 3, part1, part2)
}

// part1 counts the rows that are valid triangles.
func part1(in string) any {
	rows := slicesx.Map(input.Lines(in), input.Ints)

	return slicesx.Count(rows, isTriangle)
}

// part2 counts the valid triangles read down the columns. Each column holds
// its triangles as runs of three numbers, one after the other.
func part2(in string) any {
	rows := slicesx.Map(input.Lines(in), input.Ints)
	columns := slices.Concat(slicesx.Transpose(rows)...)

	return slicesx.Count(slicesx.Chunk(columns, 3), isTriangle)
}

// isTriangle reports whether the sum of any two sides is more than the third.
// It is enough to check that the two short sides add up to more than the long one.
func isTriangle(sides []int) bool {
	longest := slices.Max(sides)

	return slicesx.Sum(sides)-longest > longest
}
