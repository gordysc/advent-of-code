// Advent of Code 2018, day 19: Go With The Flow.
// https://adventofcode.com/2018/day/19
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 19, part1, part2)
}

// registers is the device's set of six registers.
type registers [6]int

// op computes the value an opcode writes to register C, given inputs A and B.
type op func(r *registers, a, b int) int

// ops maps each opcode name to what it computes.
var ops = map[string]op{
	"addr": func(r *registers, a, b int) int { return r[a] + r[b] },
	"addi": func(r *registers, a, b int) int { return r[a] + b },
	"mulr": func(r *registers, a, b int) int { return r[a] * r[b] },
	"muli": func(r *registers, a, b int) int { return r[a] * b },
	"banr": func(r *registers, a, b int) int { return r[a] & r[b] },
	"bani": func(r *registers, a, b int) int { return r[a] & b },
	"borr": func(r *registers, a, b int) int { return r[a] | r[b] },
	"bori": func(r *registers, a, b int) int { return r[a] | b },
	"setr": func(r *registers, a, b int) int { return r[a] },
	"seti": func(r *registers, a, b int) int { return a },
	"gtir": func(r *registers, a, b int) int { return boolInt(a > r[b]) },
	"gtri": func(r *registers, a, b int) int { return boolInt(r[a] > b) },
	"gtrr": func(r *registers, a, b int) int { return boolInt(r[a] > r[b]) },
	"eqir": func(r *registers, a, b int) int { return boolInt(a == r[b]) },
	"eqri": func(r *registers, a, b int) int { return boolInt(r[a] == b) },
	"eqrr": func(r *registers, a, b int) int { return boolInt(r[a] == r[b]) },
}

// instruction is one line of the program, with its opcode already looked up.
type instruction struct {
	f       op
	a, b, c int
}

// program is the parsed input: the register bound to the instruction pointer
// and the list of instructions.
type program struct {
	ipReg int
	code  []instruction
}

// part1 runs the program from all-zero registers until it halts and returns
// register 0. It runs every instruction; the real input takes a few million
// steps, which is quick.
func part1(in string) any {
	p := parse(in)

	var r registers

	p.run(&r, 0, func(int, int) bool { return false })

	return r[0]
}

// part2 returns register 0 when the program starts with register 0 set to 1.
//
// The real program is slow on purpose. It first builds a large number in the
// setup code at the end of the program, then jumps back to instruction 1. The
// loop from there adds up every divisor of that number with two nested loops,
// which would take about 10^14 rounds. So this runs only the setup, stops at
// the first backward jump, and sums the divisors of the largest register.
//
// The example program never jumps backward. It halts after a few steps, and
// its register 0 is the answer. The example binds the ip to register 0, so
// the ip overwrites the starting 1 at once, and the answer is 6, as in part 1.
func part2(in string) any {
	p := parse(in)

	// A short array literal fills the rest with zeros: only register 0 is 1.
	r := registers{1}
	looped := p.run(&r, 0, func(from, to int) bool { return to < from })

	if !looped {
		return r[0]
	}

	return divisorSum(max(r[0], r[1], r[2], r[3], r[4], r[5]))
}

// parse reads the "#ip N" line and the instructions after it.
func parse(in string) program {
	lines := input.Lines(in)
	p := program{ipReg: input.Int(strings.TrimPrefix(lines[0], "#ip "))}

	for _, line := range lines[1:] {
		f := strings.Fields(line)

		p.code = append(p.code, instruction{
			f: ops[f[0]],
			a: input.Int(f[1]),
			b: input.Int(f[2]),
			c: input.Int(f[3]),
		})
	}

	return p
}

// run executes the program from instruction ip until the ip leaves the
// program, and reports false. Before each step the ip is written to its bound
// register, and after the step it is read back and moved on by one.
//
// stop sees each jump from one ip to the next. If it returns true, run stops
// before the next instruction and reports true.
func (p program) run(r *registers, ip int, stop func(from, to int) bool) bool {
	for ip >= 0 && ip < len(p.code) {
		ins := p.code[ip]

		r[p.ipReg] = ip
		r[ins.c] = ins.f(r, ins.a, ins.b)

		next := r[p.ipReg] + 1
		if stop(ip, next) {
			return true
		}

		ip = next
	}

	return false
}

// divisorSum adds up every number that divides n. Divisors come in pairs d
// and n/d, so it is enough to try d up to the square root of n.
func divisorSum(n int) int {
	sum := 0

	for d := 1; d*d <= n; d++ {
		if n%d != 0 {
			continue
		}

		sum += d
		if d*d != n {
			sum += n / d
		}
	}

	return sum
}

// boolInt turns a comparison result into the 1 or 0 the device stores.
func boolInt(b bool) int {
	if b {
		return 1
	}

	return 0
}
