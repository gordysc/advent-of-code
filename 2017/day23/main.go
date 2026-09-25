// Advent of Code 2017, day 23: Coprocessor Conflagration.
// https://adventofcode.com/2017/day/23
//
// The puzzle text has no example. Every real input is the same program with a
// different start value for b, so example.txt holds that program with b set to
// 57. Its answers are 3025 for part 1 and 915 for part 2.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 23, part1, part2)
}

// The coprocessor opcodes.
const (
	set = iota
	sub
	mul
	jnz
)

// opcodes maps each instruction name to its opcode.
var opcodes = map[string]int{"set": set, "sub": sub, "mul": mul, "jnz": jnz}

// Register numbers for the registers that part 2 reads (a is 0, h is 7).
const (
	regA = 0
	regB = 1
	regC = 2
)

// arg is one operand: a register (0 to 7 for a to h) or a number.
type arg struct {
	reg   bool
	value int
}

// instr is one instruction with its two operands.
type instr struct {
	op   int
	x, y arg
}

// cpu is the coprocessor: its program, program counter and registers.
type cpu struct {
	prog []instr
	pc   int
	regs [8]int
	muls int // how many mul instructions have run
}

// part1 runs the program with every register at 0 and counts how many times
// a mul instruction runs.
func part1(in string) any {
	c := &cpu{prog: parse(in)}
	for c.running() {
		c.step()
	}

	return c.muls
}

// part2 finds the value that register h holds when the program ends with
// register a set to 1.
//
// Run as written, that takes far too long, so this reads what the program
// does instead. After a short setup it walks b from its start value up to c
// in fixed steps. For each b, the two inner loops try every pair d, e in
// [2, b) and clear flag f when d*e == b. When f is clear, h goes up by 1. So
// h ends up as the number of values of b that are not prime.
//
// The start values of b and c come from the setup, so part2 runs the program
// with a = 1 until it gets to the top of the main loop. That is the target of
// the last instruction, which jumps back there. The step comes from the
// instruction just before it, "sub b -17", which adds 17 to b.
func part2(in string) any {
	prog := parse(in)
	last := len(prog) - 1
	loopStart := last + prog[last].y.value
	step := -prog[last-1].y.value

	c := &cpu{prog: prog}
	c.regs[regA] = 1
	for c.pc != loopStart {
		c.step()
	}

	count := 0
	for b := c.regs[regB]; b <= c.regs[regC]; b += step {
		if !prime(b) {
			count++
		}
	}

	return count
}

// parse reads one instruction from every line. Operands that are a letter
// become registers, and the rest become numbers.
func parse(in string) []instr {
	var prog []instr

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)
		prog = append(prog, instr{op: opcodes[f[0]], x: parseArg(f[1]), y: parseArg(f[2])})
	}

	return prog
}

// parseArg reads one operand.
func parseArg(s string) arg {
	if s[0] >= 'a' && s[0] <= 'h' {
		return arg{reg: true, value: int(s[0] - 'a')}
	}

	return arg{value: input.Int(s)}
}

// running reports whether the program counter is still inside the program.
func (c *cpu) running() bool {
	return c.pc >= 0 && c.pc < len(c.prog)
}

// val returns the value of an operand.
func (c *cpu) val(x arg) int {
	if x.reg {
		return c.regs[x.value]
	}

	return x.value
}

// step runs the instruction at the program counter.
func (c *cpu) step() {
	ins := c.prog[c.pc]

	switch ins.op {
	case set:
		c.regs[ins.x.value] = c.val(ins.y)
	case sub:
		c.regs[ins.x.value] -= c.val(ins.y)
	case mul:
		c.regs[ins.x.value] *= c.val(ins.y)
		c.muls++
	case jnz:
		if c.val(ins.x) != 0 {
			c.pc += c.val(ins.y)
			return
		}
	}

	c.pc++
}

// prime reports whether n is a prime number, by trial division.
func prime(n int) bool {
	if n < 2 {
		return false
	}

	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}

	return true
}
