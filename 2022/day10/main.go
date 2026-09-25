// Advent of Code 2022, day 10: Cathode-Ray Tube.
// https://adventofcode.com/2022/day/10
package main

import (
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 10, part1, part2)
}

// part1 adds the signal strengths during cycles 20, 60, 100, 140, 180 and 220.
//
// The signal strength is the cycle number multiplied by the value of X
// during that cycle.
func part1(in string) any {
	total := 0

	for i, x := range registerValues(in) {
		cycle := i + 1
		if cycle%40 == 20 && cycle <= 220 {
			total += cycle * x
		}
	}

	return total
}

// part2 draws the CRT image, one character for each pixel.
//
// The CRT draws one pixel per cycle, 40 pixels on each row. It draws a lit
// pixel ('#') when the sprite covers the pixel. The sprite is 3 pixels wide
// and X is the position of its middle pixel.
func part2(in string) any {
	var screen strings.Builder

	for i, x := range registerValues(in) {
		col := i % 40
		if i > 0 && col == 0 {
			screen.WriteByte('\n')
		}

		if mathx.Abs(col-x) <= 1 {
			screen.WriteByte('#')
		} else {
			screen.WriteByte('.')
		}
	}

	return screen.String()
}

// registerValues runs the program and returns the value of X during each
// cycle. Index 0 holds the value during cycle 1.
//
// A noop takes one cycle. An addx takes two cycles, and X changes only after
// the second cycle. So the code records the old value once for noop and twice
// for addx, and then applies the change.
func registerValues(in string) []int {
	x := 1

	var values []int

	for _, line := range input.Lines(in) {
		values = append(values, x)

		if line == "noop" {
			continue
		}

		values = append(values, x)

		n, _ := strconv.Atoi(strings.TrimPrefix(line, "addx "))
		x += n
	}

	return values
}
