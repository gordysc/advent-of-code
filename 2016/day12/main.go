// Advent of Code 2016, day 12: Leonardo's Monorail.
// https://adventofcode.com/2016/day/12
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 12, part1, part2)
}

// registers holds the values of registers a, b, c and d, in that order.
type registers [4]int

// opcode is the kind of an assembunny instruction.
type opcode int

// The four assembunny instructions.
const (
	cpy opcode = iota // cpy x y copies x into register y
	inc               // inc x adds one to register x
	dec               // dec x takes one from register x
	jnz               // jnz x y jumps y instructions away when x is not zero
)

// opcodes maps the name of each instruction to its opcode.
var opcodes = map[string]opcode{"cpy": cpy, "inc": inc, "dec": dec, "jnz": jnz}

// operand is one argument of an instruction: a register or a number.
type operand struct {
	reg   int // register index, 0 to 3 for a to d, or -1 for a number
	value int // the number, when reg is -1
}

// get returns the value of the operand: the register's value or the number.
// It takes a pointer so the registers are not copied on every call.
func (o operand) get(regs *registers) int {
	if o.reg < 0 {
		return o.value
	}

	return regs[o.reg]
}

// instruction is one parsed line of the program.
type instruction struct {
	op   opcode
	x, y operand
}

// part1 runs the program with every register at zero and returns register a.
func part1(in string) any {
	regs := run(parse(in), registers{})

	return regs[0]
}

// part2 runs the program with register c set to 1 and returns register a.
func part2(in string) any {
	regs := run(parse(in), registers{2: 1})

	return regs[0]
}

// parse turns each line into an instruction. Parsing once up front keeps the
// hot loop in run free of string work.
func parse(in string) []instruction {
	var prog []instruction

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)

		ins := instruction{op: opcodes[f[0]], x: parseOperand(f[1])}
		if len(f) > 2 {
			ins.y = parseOperand(f[2])
		}

		prog = append(prog, ins)
	}

	return prog
}

// parseOperand reads a register name (a to d) or a number.
func parseOperand(s string) operand {
	if s[0] >= 'a' && s[0] <= 'd' {
		return operand{reg: int(s[0] - 'a')}
	}

	return operand{reg: -1, value: input.Int(s)}
}

// run executes the program from the first instruction until the program
// counter leaves the program, and returns the final registers.
func run(prog []instruction, regs registers) registers {
	for pc := 0; pc >= 0 && pc < len(prog); {
		ins := prog[pc]

		switch ins.op {
		case cpy:
			regs[ins.y.reg] = ins.x.get(&regs)
		case inc:
			regs[ins.x.reg]++
		case dec:
			regs[ins.x.reg]--
		case jnz:
			if ins.x.get(&regs) != 0 {
				pc += ins.y.get(&regs)
				continue
			}
		}

		pc++
	}

	return regs
}
