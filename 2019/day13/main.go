// Advent of Code 2019, day 13: Care Package.
// https://adventofcode.com/2019/day/13
//
// The puzzle text has no example program. example.txt holds a small handmade
// Intcode program. It starts with a harmless add (1,0,0,50 writes to address
// 50), so that address 0 holds an opcode like the real game does. Part 2's
// free-play poke turns it into a multiply, which is just as harmless. Then it
// draws five tiles with immediate-mode output instructions (104,x,104,y,104,id):
// three blocks in a row, the ball below the middle block, and the paddle below
// the ball. Last it outputs a score of 100 and halts. So part 1 gives 3, and
// part 2 gives 100 without ever moving the joystick.
package main

import (
	"aoc/lib/intcode"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 13, part1, part2)
}

// The tile ids the game draws. Walls (1) and empty tiles (0) do not matter.
const (
	block  = 2
	paddle = 3
	ball   = 4
)

// part1 counts the block tiles on the first screen the game draws.
func part1(in string) any {
	out := intcode.New(intcode.Parse(in)).Run()
	blocks := 0

	// The output is a flat list of (x, y, tile id) triples, so step by 3.
	for i := 0; i+2 < len(out); i += 3 {
		if out[i+2] == block {
			blocks++
		}
	}

	return blocks
}

// part2 plays the game to the end and returns the final score. Each turn the
// game runs until it wants the joystick position, and its output tells where
// the ball and the paddle are now. Moving the paddle towards the ball means it
// is always under the ball in time, so the game never ends early and every
// block gets broken.
func part2(in string) any {
	m := intcode.New(intcode.Parse(in))
	m.Poke(0, 2) // address 0 holds the number of quarters: 2 means free play

	score, ballX, paddleX := 0, 0, 0

	for {
		out := m.Run()

		for i := 0; i+2 < len(out); i += 3 {
			x, y, id := out[i], out[i+1], out[i+2]

			switch {
			case x == -1 && y == 0:
				score = id // this triple is not a tile: its third value is the score
			case id == ball:
				ballX = x
			case id == paddle:
				paddleX = x
			}
		}

		if m.Halted() {
			return score
		}

		// Sign gives -1, 0 or 1, which are exactly the joystick positions
		// left, neutral and right.
		m.Send(mathx.Sign(ballX - paddleX))
	}
}
