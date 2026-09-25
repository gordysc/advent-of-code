// Advent of Code 2016, day 25: Clock Signal.
// https://adventofcode.com/2016/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
//
// The puzzle text has no example. example.txt holds a small handmade program
// that sends 0, 1, 0, 1, ... when a is 3 and a negative number for a smaller
// a, so its answer is 3.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 25, part1, nil)
}

// The assembunny opcodes.
const (
	cpy = iota
	inc
	dec
	jnz
	out
)

// opcodes maps each instruction name to its opcode.
var opcodes = map[string]int{"cpy": cpy, "inc": inc, "dec": dec, "jnz": jnz, "out": out}

// arg is one operand: a register (0 to 3 for a to d) or a number.
type arg struct {
	reg   bool
	value int
}

// instr is one instruction with up to two operands.
type instr struct {
	op   int
	x, y arg
}

// state is everything that decides what the program does next: the program
// counter, the registers, and the bit the clock signal needs next.
type state struct {
	pc   int
	regs [4]int
	want int
}

// part1 finds the lowest positive start value for register a that makes the
// program send the clock signal 0, 1, 0, 1, ... forever.
func part1(in string) any {
	prog := parse(in)

	for a := 1; ; a++ {
		if clock(prog, a) {
			return a
		}
	}
}

// parse reads one instruction from every line. Operands that are a letter
// become registers, and the rest become numbers.
func parse(in string) []instr {
	var prog []instr

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)
		ins := instr{op: opcodes[f[0]], x: parseArg(f[1])}
		if len(f) > 2 {
			ins.y = parseArg(f[2])
		}

		prog = append(prog, ins)
	}

	return prog
}

// parseArg reads one operand.
func parseArg(s string) arg {
	if s[0] >= 'a' && s[0] <= 'd' {
		return arg{reg: true, value: int(s[0] - 'a')}
	}

	return arg{value: input.Int(s)}
}

// clock runs the program with a in register a and reports whether it sends
// 0, 1, 0, 1, ... forever.
//
// The program has no end, so clock cannot watch the signal forever. Instead
// it notes the full state each time a bit goes out. When a state comes back,
// the program is in a loop: it sends the same bits from there on, and all of
// them were correct the first time round.
func clock(prog []instr, a int) bool {
	regs := [4]int{a}
	want := 0
	seen := map[state]bool{}

	// val returns the value of an operand.
	val := func(x arg) int {
		if x.reg {
			return regs[x.value]
		}

		return x.value
	}

	for pc := 0; pc >= 0 && pc < len(prog); {
		ins := prog[pc]

		switch ins.op {
		case cpy:
			regs[ins.y.value] = val(ins.x)
		case inc:
			regs[ins.x.value]++
		case dec:
			regs[ins.x.value]--
		case jnz:
			if val(ins.x) != 0 {
				pc += val(ins.y)
				continue
			}
		case out:
			if val(ins.x) != want {
				return false
			}

			s := state{pc, regs, want}
			if seen[s] {
				return true
			}

			seen[s] = true
			want = 1 - want
		}

		pc++
	}

	// The program stopped, so the signal does not go on forever.
	return false
}
