// Advent of Code 2015, day 23: Opening the Turing Lock.
// https://adventofcode.com/2015/day/23
//
// The example in the puzzle text only touches register a: it ends with a = 2
// in part 1 and a = 7 in part 2. Both parts report register b, which stays 0
// for the example. The real input computes a Collatz sequence in b.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 23, part1, part2)
}

// part1 runs the program with both registers at 0 and reports register b.
func part1(in string) any {
	return run(parse(in), 0)
}

// part2 runs the program with register a set to 1 and reports register b.
func part2(in string) any {
	return run(parse(in), 1)
}

// instruction is one line of the program. Not every field is used by every
// operation: hlf, tpl and inc use only reg; jmp uses only offset; jie and jio
// use both.
type instruction struct {
	op     string
	reg    byte
	offset int
}

// parse turns each line into an instruction. The register is the first byte of
// the second field, so "a," and "a" both give 'a'. The offset is the last
// integer on the line, if there is one.
func parse(in string) []instruction {
	var program []instruction

	for _, line := range input.Lines(in) {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			panic("bad instruction: " + line)
		}

		ins := instruction{op: fields[0]}

		// jmp takes no register, only an offset.
		if ins.op != "jmp" {
			ins.reg = fields[1][0]
		}

		if nums := input.Ints(line); len(nums) > 0 {
			ins.offset = nums[len(nums)-1]
		}

		program = append(program, ins)
	}

	return program
}

// run executes the program until the instruction pointer leaves it and
// returns the final value of register b.
func run(program []instruction, a int) int {
	// The two registers, indexed by 'a' - 'a' and 'b' - 'a'.
	regs := [2]int{a, 0}

	for pc := 0; pc >= 0 && pc < len(program); {
		ins := program[pc]

		// jmp has no register, so only look one up when the line names one.
		var r *int
		if ins.reg != 0 {
			r = &regs[ins.reg-'a']
		}

		// Every instruction but a taken jump moves to the next line, so start
		// with that and let the jumps overwrite it.
		next := pc + 1

		switch ins.op {
		case "hlf":
			*r /= 2
		case "tpl":
			*r *= 3
		case "inc":
			*r++
		case "jmp":
			next = pc + ins.offset
		case "jie":
			if *r%2 == 0 {
				next = pc + ins.offset
			}
		case "jio":
			// "jio" jumps if the register is one, not odd.
			if *r == 1 {
				next = pc + ins.offset
			}
		default:
			panic("unknown instruction: " + ins.op)
		}

		pc = next
	}

	return regs[1]
}
