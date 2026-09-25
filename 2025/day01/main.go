// Advent of Code 2025, day 1: Secret Entrance.
// https://adventofcode.com/2025/day/1
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// dialSize is the number of positions on the dial, 0 to 99.
const dialSize = 100

// startPos is the position of the dial before the first rotation.
const startPos = 50

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 1, part1, part2)
}

// part1 counts the rotations that leave the dial at 0.
func part1(in string) any {
	pos := startPos
	count := 0

	for _, step := range parse(in) {
		// mathx.Mod always gives a result from 0 to dialSize-1. The % operator
		// in Go keeps the sign of the left side, so -5 % 100 is -5.
		pos = mathx.Mod(pos+step, dialSize)

		if pos == 0 {
			count++
		}
	}

	return count
}

// part2 counts every click that puts the dial at 0, also the clicks in the
// middle of a rotation.
func part2(in string) any {
	pos := startPos
	count := 0

	for _, step := range parse(in) {
		count += zeroHits(pos, step)
		pos = mathx.Mod(pos+step, dialSize)
	}

	return count
}

// zeroHits returns how many clicks of one rotation stop at 0. The dial starts
// at pos and moves step clicks. A positive step turns right and a negative
// step turns left.
//
// A right turn passes 0 at pos+step = 100, 200, and so on, so it hits 0
// (pos+step)/100 times. A left turn is the same as a right turn on a mirrored
// dial, where position p becomes (100-p) % 100. The % keeps a start at 0 at 0,
// so a left turn from 0 does not count the start as a hit.
func zeroHits(pos, step int) int {
	if step >= 0 {
		return (pos + step) / dialSize
	}

	mirrored := (dialSize - pos) % dialSize

	return (mirrored - step) / dialSize
}

// parse reads the rotations as signed click counts: "L68" gives -68 and "R48"
// gives 48.
func parse(in string) []int {
	var steps []int

	for _, line := range input.Lines(in) {
		if line == "" {
			continue
		}

		n := input.Int(line[1:])
		if line[0] == 'L' {
			n = -n
		}

		steps = append(steps, n)
	}

	return steps
}
