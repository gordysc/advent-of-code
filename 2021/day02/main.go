// Advent of Code 2021, day 2: Dive!
// https://adventofcode.com/2021/day/2
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 2, part1, part2)
}

// command is one step of the planned course, for example "forward 5".
type command struct {
	dir string
	n   int
}

// part1 follows the course, where "down" and "up" change the depth directly,
// and multiplies the final horizontal position by the final depth.
func part1(in string) any {
	pos, depth := 0, 0

	for _, c := range parse(in) {
		switch c.dir {
		case "forward":
			pos += c.n
		case "down":
			depth += c.n
		case "up":
			depth -= c.n
		}
	}

	return pos * depth
}

// part2 follows the course again, but "down" and "up" now change the aim.
// Each "forward" moves the submarine ahead and also changes the depth by the
// aim multiplied by the distance.
func part2(in string) any {
	pos, depth, aim := 0, 0, 0

	for _, c := range parse(in) {
		switch c.dir {
		case "forward":
			pos += c.n
			depth += aim * c.n
		case "down":
			aim += c.n
		case "up":
			aim -= c.n
		}
	}

	return pos * depth
}

// parse reads one command from each line.
func parse(in string) []command {
	var cmds []command

	for _, line := range input.Lines(in) {
		dir, n, _ := strings.Cut(line, " ")
		cmds = append(cmds, command{dir, input.Int(n)})
	}

	return cmds
}
