// Advent of Code 2023, day 18: Lavaduct Lagoon.
// https://adventofcode.com/2023/day/18
package main

import (
	"strconv"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 18, part1, part2)
}

// step is one dig instruction: a direction and a number of meters.
type step struct {
	dir grid.Point
	n   int
}

// part1 gives the volume of the lagoon when the direction and the number
// come from the first two fields of each line.
func part1(in string) any {
	var steps []step

	for _, line := range input.Lines(in) {
		f := input.Fields(line)
		steps = append(steps, step{grid.DirFromRune[rune(f[0][0])], input.Int(f[1])})
	}

	return lagoonSize(steps)
}

// hexDirs maps the last hex digit of a color code to a direction.
var hexDirs = []grid.Point{grid.Right, grid.Down, grid.Left, grid.Up}

// part2 gives the volume of the lagoon when each instruction comes from the
// color code. The first five hex digits are the number of meters, and the
// last digit is the direction. The numbers are too large to dig a grid, so
// lagoonSize calculates the size from the corners only.
func part2(in string) any {
	var steps []step

	for _, line := range input.Lines(in) {
		code := input.Fields(line)[2] // code has the form "(#70c710)".
		hex := code[2:8]

		n, err := strconv.ParseInt(hex[:5], 16, 64)
		if err != nil {
			panic(err)
		}

		steps = append(steps, step{hexDirs[hex[5]-'0'], int(n)})
	}

	return lagoonSize(steps)
}

// lagoonSize gives the number of cubic meters in the dug-out lagoon: the
// trench plus all cells inside it.
//
// The shoelace formula gives the area of the polygon through the centers of
// the trench cells. Pick's theorem says that area = inside + boundary/2 - 1,
// where inside and boundary count the cells. So inside = area - boundary/2
// + 1, and the total is inside + boundary = area + boundary/2 + 1.
func lagoonSize(steps []step) int {
	var p grid.Point
	twiceArea := 0
	boundary := 0

	for _, s := range steps {
		q := p.Add(s.dir.Scale(s.n))

		// Each term of the shoelace sum is the cross product of two corners
		// that follow each other. The sum is twice the signed area.
		twiceArea += p.X*q.Y - q.X*p.Y
		boundary += s.n
		p = q
	}

	return mathx.Abs(twiceArea)/2 + boundary/2 + 1
}
