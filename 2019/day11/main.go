// Advent of Code 2019, day 11: Space Police.
// https://adventofcode.com/2019/day/11
//
// The puzzle text has no example program. example.txt holds a small handmade
// Intcode program that plays the robot's moves from the worked example. Seven
// times it reads a colour (3,50 stores it at address 50 and never looks at it
// again) and then outputs a fixed colour and turn with two immediate-mode
// output instructions (104,colour,104,turn). After the seventh move it halts.
// It paints 6 panels, as in the puzzle text, so part 1 gives 6. Part 2 draws
// the white panels it leaves behind.
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 11, part1, part2)
}

// black and white are the two colours a panel can have.
const (
	black = 0
	white = 1
)

// brain decides what the robot does. It gets the colour of the panel under
// the robot and returns the colour to paint it and the way to turn (0 left,
// 1 right). ok is false when the brain has stopped and the robot is done.
type brain func(colour int) (paint, turn int, ok bool)

// part1 counts the panels the robot paints at least once, starting on black.
func part1(in string) any {
	return len(run(black, program(in)))
}

// part2 starts the robot on a white panel and draws the white panels it
// leaves, cropped to the area they cover. Grid.String ends every row with a
// newline, so the last one is trimmed.
func part2(in string) any {
	hull := run(white, program(in))

	var whites []grid.Point
	for p, colour := range hull {
		if colour == white {
			whites = append(whites, p)
		}
	}

	if len(whites) == 0 {
		return nil
	}

	lo, hi := whites[0], whites[0]
	for _, p := range whites {
		lo = grid.P(min(lo.X, p.X), min(lo.Y, p.Y))
		hi = grid.P(max(hi.X, p.X), max(hi.Y, p.Y))
	}

	picture := grid.New[byte](hi.X-lo.X+1, hi.Y-lo.Y+1)
	for i := range picture.Cells {
		picture.Cells[i] = ' '
	}

	for _, p := range whites {
		picture.Set(p.Sub(lo), '#')
	}

	return strings.TrimRight(picture.String(), "\n")
}

// program returns a brain that runs the Intcode program. Each call sends the
// colour to the program and runs it until it wants the next colour. By then
// it has written the paint colour and the turn. When the program halts
// instead, it writes nothing and the brain reports that it is done.
//
// The brain is a closure: the function it returns keeps using the machine m
// from this call, so the program's memory carries over from one call to the
// next.
func program(in string) brain {
	m := intcode.New(intcode.Parse(in))

	return func(colour int) (int, int, bool) {
		m.Send(colour)
		out := m.Run()

		if len(out) < 2 {
			return 0, 0, false
		}

		return out[0], out[1], true
	}
}

// run moves the robot until its brain stops. The robot starts at (0, 0)
// facing up, on a panel of the start colour; every other panel starts
// black. It returns the colour of every panel the robot painted. A panel
// that was never painted is not in the map, even the start panel.
func run(start int, b brain) map[grid.Point]int {
	hull := map[grid.Point]int{}
	pos, dir := grid.P(0, 0), grid.Up

	for {
		colour, painted := hull[pos]
		if !painted {
			colour = black
			if pos == grid.P(0, 0) {
				colour = start
			}
		}

		paint, turn, ok := b(colour)
		if !ok {
			return hull
		}

		hull[pos] = paint

		if turn == 0 {
			dir = dir.TurnLeft()
		} else {
			dir = dir.TurnRight()
		}

		pos = pos.Add(dir)
	}
}
