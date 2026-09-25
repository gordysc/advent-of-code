// Advent of Code 2017, day 19: A Series of Tubes.
// https://adventofcode.com/2017/day/19
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 19, part1, part2)
}

// diagram is the routing diagram as text rows. The rows can have different
// lengths, so it is not a grid.Grid.
type diagram []string

// at returns the character at p, or a space when p is off the diagram.
func (d diagram) at(p grid.Point) byte {
	if p.Y < 0 || p.Y >= len(d) || p.X < 0 || p.X >= len(d[p.Y]) {
		return ' '
	}

	return d[p.Y][p.X]
}

// part1 returns the letters the packet passes, in order.
func part1(in string) any {
	letters, _ := follow(in)

	return letters
}

// part2 returns how many steps the packet takes, the first square included.
func part2(in string) any {
	_, steps := follow(in)

	return steps
}

// follow walks the packet along the path. It starts at the | in the top row
// and goes down. It goes straight until it reaches a +, where it turns to the
// side that continues the path. The walk ends when the next square is empty.
// It returns the letters it passed and the number of squares it visited.
func follow(in string) (string, int) {
	d := diagram(input.Lines(in))
	pos := grid.P(strings.IndexByte(d[0], '|'), 0)
	dir := grid.Down

	var letters strings.Builder
	steps := 0

	for {
		c := d.at(pos)
		if c == ' ' {
			return letters.String(), steps
		}

		steps++

		switch {
		case c >= 'A' && c <= 'Z':
			letters.WriteByte(c)
		case c == '+':
			// The path never goes back, so try the two turns. Keep the one
			// whose next square is not empty.
			if left := dir.TurnLeft(); d.at(pos.Add(left)) != ' ' {
				dir = left
			} else {
				dir = dir.TurnRight()
			}
		}

		pos = pos.Add(dir)
	}
}
