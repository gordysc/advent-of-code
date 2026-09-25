// Advent of Code 2019, day 5: Sunny with a Chance of Asteroids.
// https://adventofcode.com/2019/day/5
package main

import (
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 5, part1, part2)
}

// airConditioner and thermalRadiator are the system IDs that parts 1 and 2
// send to the diagnostic program.
//
// The example program from the puzzle text is not a diagnostic program. It
// reads one number and outputs 999 if it is below 8, 1000 if it is 8, and
// 1001 if it is above 8. Both IDs are below 8, so both parts give 999.
const (
	airConditioner  = 1
	thermalRadiator = 5
)

// part1 runs the diagnostic for the air conditioner and returns the
// diagnostic code. Each test before it outputs 0 when it passes, and the
// code is the last output.
func part1(in string) any {
	return diagnose(in, airConditioner)
}

// part2 runs the diagnostic for the thermal radiator controller. This time
// the program outputs only the diagnostic code.
func part2(in string) any {
	return diagnose(in, thermalRadiator)
}

// diagnose runs the program with one system ID as input and returns the last
// output, or nil when the program outputs nothing.
func diagnose(in string, id int) any {
	m := intcode.New(intcode.Parse(in))
	m.Send(id)

	out := m.Run()
	if len(out) == 0 {
		return nil
	}

	return out[len(out)-1]
}
