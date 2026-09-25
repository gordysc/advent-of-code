// Advent of Code 2016, day 23: Safe Cracking.
// https://adventofcode.com/2016/day/23
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 23, part1, part2)
}

// eggs1 and eggs2 are the values register a starts with in each part. The
// worked example in the puzzle text sets a to 2 on its first line, so the
// start value does not matter there and both parts give 3.
const (
	eggs1 = 7
	eggs2 = 12
)

// The assembunny opcodes.
const (
	cpy = iota
	inc
	dec
	jnz
	tgl
)

// opcodes maps each instruction name to its opcode.
var opcodes = map[string]int{"cpy": cpy, "inc": inc, "dec": dec, "jnz": jnz, "tgl": tgl}

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

// part1 runs the program with the first egg count in register a.
func part1(in string) any {
	return run(parse(in), eggs1)
}

// part2 runs the program with the second egg count in register a.
func part2(in string) any {
	return run(parse(in), eggs2)
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

// run executes the program with a in register a and returns register a when
// the program ends. tgl rewrites the program, so run works on its own copy.
//
// The programs compute a factorial with nested loops that add one at a time,
// which takes billions of steps for part 2. run spots the multiply loop and
// does it in one go instead.
func run(prog []instr, a int) int {
	prog = slices.Clone(prog)
	regs := [4]int{a}

	// val returns the value of an operand.
	val := func(x arg) int {
		if x.reg {
			return regs[x.value]
		}

		return x.value
	}

	for pc := 0; pc >= 0 && pc < len(prog); {
		if dst, x, outer, inner, ok := multiply(prog[pc:]); ok {
			regs[dst] += val(x) * regs[outer]
			regs[outer], regs[inner] = 0, 0
			pc += 6
			continue
		}

		ins := prog[pc]

		switch ins.op {
		case cpy:
			// A toggle can make a cpy with a number as its target. It does nothing.
			if ins.y.reg {
				regs[ins.y.value] = val(ins.x)
			}
		case inc:
			if ins.x.reg {
				regs[ins.x.value]++
			}
		case dec:
			if ins.x.reg {
				regs[ins.x.value]--
			}
		case jnz:
			if val(ins.x) != 0 {
				pc += val(ins.y)
				continue
			}
		case tgl:
			if t := pc + val(ins.x); t >= 0 && t < len(prog) {
				prog[t].op = toggle(prog[t].op)
			}
		}

		pc++
	}

	return regs[0]
}

// toggle returns the opcode that tgl turns op into. One-operand instructions
// become dec if they were inc and inc otherwise. Two-operand instructions
// become cpy if they were jnz and jnz otherwise.
func toggle(op int) int {
	switch op {
	case inc:
		return dec
	case dec, tgl:
		return inc
	case jnz:
		return cpy
	default:
		return jnz
	}
}

// multiply checks whether code starts with this multiply loop:
//
//	cpy x c
//	inc a
//	dec c
//	jnz c -2
//	dec d
//	jnz d -5
//
// The inner three lines add c to a, and the outer loop repeats that d times,
// so together they add x*d to a and leave c and d at zero. The register names
// can differ. It returns register a, operand x, register d and register c.
// The loop is checked again every time, because tgl can change the lines in it.
func multiply(code []instr) (dst int, x arg, outer, inner int, ok bool) {
	if len(code) < 6 {
		return
	}

	c, a, d := code[0].y, code[1].x, code[4].x
	if code[0].op != cpy || !c.reg || !a.reg || !d.reg {
		return
	}

	loop := code[1].op == inc &&
		code[2].op == dec && code[2].x == c &&
		code[3].op == jnz && code[3].x == c && code[3].y == arg{value: -2} &&
		code[4].op == dec &&
		code[5].op == jnz && code[5].x == d && code[5].y == arg{value: -5}

	x = code[0].x
	distinct := a != c && a != d && c != d && x != a && x != c && x != d
	if !loop || !distinct {
		return
	}

	return a.value, x, d.value, c.value, true
}
