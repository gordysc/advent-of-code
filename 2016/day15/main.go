// Advent of Code 2016, day 15: Timing is Everything.
// https://adventofcode.com/2016/day/15
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 15, part1, part2)
}

// disc is one spinning disc: its number of positions and its position at
// time 0. The slot is at position 0.
type disc struct {
	positions, start int
}

// part1 finds the first time to press the button so the capsule falls
// through every disc.
func part1(in string) any {
	return firstTime(parse(in))
}

// part2 adds a new disc below the others, with 11 positions, at position 0
// at time 0.
func part2(in string) any {
	return firstTime(append(parse(in), disc{positions: 11, start: 0}))
}

// parse reads one disc from each line. The numbers on a line are the disc
// number, its number of positions, the time 0, and its start position.
func parse(in string) []disc {
	var discs []disc

	for _, line := range input.Lines(in) {
		n := input.UInts(line)
		discs = append(discs, disc{positions: n[1], start: n[3]})
	}

	return discs
}

// firstTime finds the first button time t at which the capsule passes every
// disc. The capsule reaches disc k (counting from 1) at time t+k, so disc k
// must be at position 0 then: (start + t + k) mod positions == 0.
//
// It works one disc at a time, like a sieve. When a time works for the discs
// so far, adding the LCM of their sizes gives the next time that also works
// for them. So step by that LCM until the next disc also lines up.
func firstTime(discs []disc) int {
	t, step := 0, 1

	for i, d := range discs {
		k := i + 1

		for (d.start+t+k)%d.positions != 0 {
			t += step
		}

		step = mathx.LCM(step, d.positions)
	}

	return t
}
