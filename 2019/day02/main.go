// Advent of Code 2019, day 2: 1202 Program Alarm.
// https://adventofcode.com/2019/day/2
package main

import (
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 2, part1, part2)
}

// alarmNoun and alarmVerb put the program back into its "1202 program alarm"
// state for part 1. target is the output that part 2 looks for.
//
// The worked example in the puzzle text runs its program unchanged, and it
// leaves 3500 at address 0. Both parts patch addresses 1 and 2 first, so they
// do not give the puzzle's answers for the example: part 1 gives 100, and no
// noun and verb make the example output 19690720, so part 2 has no answer.
const (
	alarmNoun = 12
	alarmVerb = 2
	target    = 19690720
)

// part1 runs the program in the "1202 program alarm" state and returns the
// value left at address 0.
func part1(in string) any {
	return run(intcode.Parse(in), alarmNoun, alarmVerb)
}

// part2 tries every noun and verb from 0 to 99 and returns 100*noun+verb for
// the pair that makes the program output the target value.
func part2(in string) any {
	program := intcode.Parse(in)

	for noun := range 100 {
		for verb := range 100 {
			if run(program, noun, verb) == target {
				return 100*noun + verb
			}
		}
	}

	return nil
}

// run puts noun at address 1 and verb at address 2, runs the program, and
// returns the value at address 0. intcode.New copies the program, so each run
// starts from clean memory and the caller's slice never changes.
func run(program []int, noun, verb int) int {
	m := intcode.New(program)
	m.Poke(1, noun)
	m.Poke(2, verb)

	m.Run()

	return m.Peek(0)
}
