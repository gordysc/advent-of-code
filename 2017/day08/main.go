// Advent of Code 2017, day 8: I Heard You Like Registers.
// https://adventofcode.com/2017/day/8
package main

import (
	"math"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 8, part1, part2)
}

// part1 finds the largest register value after the last instruction.
func part1(in string) any {
	final, _ := run(in)

	return final
}

// part2 finds the largest value any register held at any time.
func part2(in string) any {
	_, highest := run(in)

	return highest
}

// run does each instruction, such as "b inc 5 if a > 1", in order. Registers
// start at 0. It returns the largest register at the end and the largest
// value seen during the run.
func run(in string) (int, int) {
	regs := map[string]int{}
	highest := 0

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)

		// A register that is named but never changed still holds 0, so it
		// must count when part 1 looks for the largest register.
		for _, name := range []string{f[0], f[4]} {
			if _, ok := regs[name]; !ok {
				regs[name] = 0
			}
		}

		if !compare(regs[f[4]], f[5], input.Int(f[6])) {
			continue
		}

		delta := input.Int(f[2])
		if f[1] == "dec" {
			delta = -delta
		}

		regs[f[0]] += delta
		highest = max(highest, regs[f[0]])
	}

	final := math.MinInt
	for _, v := range regs {
		final = max(final, v)
	}

	return final, highest
}

// compare checks the condition "a op b".
func compare(a int, op string, b int) bool {
	switch op {
	case "<":
		return a < b
	case "<=":
		return a <= b
	case ">":
		return a > b
	case ">=":
		return a >= b
	case "==":
		return a == b
	case "!=":
		return a != b
	}

	panic("unknown operator " + op)
}
