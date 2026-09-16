// Advent of Code 2015, day 6: Probably a Fire Hazard.
// https://adventofcode.com/2015/day/6
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 6, part1, part2)
}

// size is the width and height of the light grid.
const size = 1000

// action is what an instruction does to one light. The three kinds are
// distinguished by the words at the start of the line.
type action int

const (
	turnOn action = iota
	turnOff
	toggle
)

// instruction is one parsed line: an action and the inclusive corners of the
// rectangle it applies to.
type instruction struct {
	act    action
	x1, y1 int
	x2, y2 int
}

// part1 treats every light as on or off and counts how many end up on.
func part1(in string) any {
	lights := apply(in, func(v int, act action) int {
		switch act {
		case turnOn:
			return 1
		case turnOff:
			return 0
		default:
			return 1 - v
		}
	})

	return slicesx.Sum(lights.Cells)
}

// part2 treats every light as a brightness level and sums the total.
// Turn on adds one, turn off removes one down to zero, toggle adds two.
func part2(in string) any {
	lights := apply(in, func(v int, act action) int {
		switch act {
		case turnOn:
			return v + 1
		case turnOff:
			return max(v-1, 0)
		default:
			return v + 2
		}
	})

	return slicesx.Sum(lights.Cells)
}

// apply runs every instruction over a fresh grid, using fn to compute the new
// value of each light in the rectangle from its old value and the action.
func apply(in string, fn func(v int, act action) int) grid.Grid[int] {
	g := grid.New[int](size, size)

	for _, ins := range parse(in) {
		for y := ins.y1; y <= ins.y2; y++ {
			// Slicing one row of the flat cell array avoids a bounds check and an
			// index calculation per light, which matters over a million cells.
			row := g.Cells[y*g.W+ins.x1 : y*g.W+ins.x2+1]

			for i, v := range row {
				row[i] = fn(v, ins.act)
			}
		}
	}

	return g
}

// parse turns each line into an instruction. The action comes from the leading
// words and the four coordinates from the numbers in the rest of the line.
func parse(in string) []instruction {
	var out []instruction

	for _, line := range input.Lines(in) {
		var act action
		switch {
		case strings.HasPrefix(line, "turn on"):
			act = turnOn
		case strings.HasPrefix(line, "turn off"):
			act = turnOff
		default:
			act = toggle
		}

		n := input.Ints(line)
		out = append(out, instruction{act, n[0], n[1], n[2], n[3]})
	}

	return out
}
