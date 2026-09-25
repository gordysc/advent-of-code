// Advent of Code 2020, day 8: Handheld Halting.
// https://adventofcode.com/2020/day/8
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 8, part1, part2)
}

// instr is one line of boot code: an operation and its signed argument.
type instr struct {
	op  string
	arg int
}

// part1 gives the accumulator value just before any instruction runs a
// second time.
func part1(in string) any {
	acc, _ := run(parse(in))

	return acc
}

// part2 finds the one corrupt instruction. It swaps each jmp to nop, or each
// nop to jmp, one at a time, and runs the program. The first swap that lets
// the program end gives the answer. The program is short, so this brute
// force is fast.
func part2(in string) any {
	prog := parse(in)

	for i, ins := range prog {
		swap := map[string]string{"jmp": "nop", "nop": "jmp"}[ins.op]
		if swap == "" {
			continue
		}

		prog[i].op = swap
		acc, ok := run(prog)
		prog[i].op = ins.op

		if ok {
			return acc
		}
	}

	return nil
}

// run executes the program until it loops or ends. It returns the
// accumulator and true when the program ends by stepping just past the last
// instruction. It returns false when an instruction is about to run twice,
// because the program has no state other than its position and the
// accumulator, so a repeat position means an infinite loop.
func run(prog []instr) (int, bool) {
	seen := make([]bool, len(prog))
	acc, pc := 0, 0

	for pc >= 0 && pc < len(prog) {
		if seen[pc] {
			return acc, false
		}

		seen[pc] = true

		switch prog[pc].op {
		case "acc":
			acc += prog[pc].arg
			pc++
		case "jmp":
			pc += prog[pc].arg
		default:
			pc++
		}
	}

	return acc, pc == len(prog)
}

// parse reads each "op +n" line into an instruction.
func parse(in string) []instr {
	var prog []instr
	for _, line := range input.Lines(in) {
		op, arg, _ := strings.Cut(line, " ")
		prog = append(prog, instr{op, input.Int(arg)})
	}

	return prog
}
