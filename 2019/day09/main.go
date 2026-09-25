// Advent of Code 2019, day 9: Sensor Boost.
// https://adventofcode.com/2019/day/9
package main

import (
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 9, part1, part2)
}

// testMode and sensorMode are the inputs that start the BOOST program in its
// two modes.
//
// The worked examples in the puzzle text read no input. They only test the
// relative base and large numbers. example.txt is the one that prints
// 1125899906842624. Both parts ignore their input value for it and return
// the last output, so both give 1125899906842624.
const (
	testMode   = 1
	sensorMode = 2
)

// part1 runs BOOST in test mode. If every instruction works, the only output
// is the BOOST keycode.
func part1(in string) any {
	return run(in, testMode)
}

// part2 runs BOOST in sensor boost mode and returns the distress signal's
// coordinates.
func part2(in string) any {
	return run(in, sensorMode)
}

// run sends one input value to the program, runs it to the end and returns
// its last output. In test mode BOOST also prints the opcodes that failed
// before the keycode, so the last output is the one to report.
func run(in string, mode int) any {
	m := intcode.New(intcode.Parse(in))
	m.Send(mode)

	out := m.Run()
	if len(out) == 0 {
		return nil
	}

	return out[len(out)-1]
}
