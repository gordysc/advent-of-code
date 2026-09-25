// Advent of Code 2017, day 18: Duet.
// https://adventofcode.com/2017/day/18
package main

import (
	"strings"

	"aoc/lib/ds"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 18, part1, part2)
}

// opcode is the kind of a Duet instruction.
type opcode int

// The seven Duet instructions.
const (
	snd opcode = iota // snd x plays a sound (part 1) or sends x (part 2)
	set               // set x y copies y into register x
	add               // add x y adds y to register x
	mul               // mul x y multiplies register x by y
	mod               // mod x y sets register x to the remainder of x / y
	rcv               // rcv x recovers the last sound (part 1) or receives into x (part 2)
	jgz               // jgz x y jumps y instructions away when x is greater than zero
)

// opcodes maps the name of each instruction to its opcode.
var opcodes = map[string]opcode{
	"snd": snd, "set": set, "add": add, "mul": mul, "mod": mod, "rcv": rcv, "jgz": jgz,
}

// registers holds the values of registers a to z, in that order.
type registers [26]int

// operand is one argument of an instruction: a register or a number.
type operand struct {
	reg   int // register index, 0 to 25 for a to z, or -1 for a number
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

// machine is one running copy of the program: its registers, its program
// counter and, for part 2, how many values it has sent.
type machine struct {
	regs registers
	pc   int
	sent int
}

// part1 plays sounds until the first rcv whose value is not zero, and returns
// the last sound played before it. It returns nil if the program ends first.
func part1(in string) any {
	prog := parse(in)
	var m machine
	sound := 0

	for m.inside(prog) {
		ins := prog[m.pc]

		switch ins.op {
		case snd:
			sound = ins.x.get(&m.regs)
			m.pc++
		case rcv:
			if ins.x.get(&m.regs) != 0 {
				return sound
			}

			m.pc++
		default:
			m.exec(ins)
		}
	}

	return nil
}

// part2 runs two copies of the program, with register p set to 0 and 1. Each
// one sends values to the other's queue. They take turns, and each turn runs a
// program until it waits on an empty queue or ends. When neither program can
// run one more instruction, they are stuck, and part2 returns how many values
// program 1 sent.
//
// The example.txt file holds the part 1 example from the puzzle text. Part 2
// gives 1 on it: each program sends one value, receives the other's value,
// and then waits forever on its second rcv. The puzzle text has a separate
// part 2 example (snd 1, snd 2, snd p, rcv a, rcv b, rcv c, rcv d), which
// gives 3.
func part2(in string) any {
	prog := parse(in)
	p := int('p' - 'a')

	var m0, m1 machine
	m1.regs[p] = 1
	q0, q1 := ds.NewQueue[int](), ds.NewQueue[int]()

	for {
		ran := m0.run(prog, q0, q1)
		ran += m1.run(prog, q1, q0)

		if ran == 0 {
			return m1.sent
		}
	}
}

// parse turns each line into an instruction. Parsing once up front keeps the
// hot loop free of string work.
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

// parseOperand reads a register name (a to z) or a number.
func parseOperand(s string) operand {
	if s[0] >= 'a' && s[0] <= 'z' {
		return operand{reg: int(s[0] - 'a')}
	}

	return operand{reg: -1, value: input.Int(s)}
}

// inside reports whether the program counter still points at an instruction.
// The program ends when it jumps outside.
func (m *machine) inside(prog []instruction) bool {
	return m.pc >= 0 && m.pc < len(prog)
}

// exec runs one instruction that the two parts treat the same way: set, add,
// mul, mod or jgz. The machine is a pointer receiver, so the changes to its
// registers and program counter stay after the call.
func (m *machine) exec(ins instruction) {
	regs := &m.regs

	switch ins.op {
	case set:
		regs[ins.x.reg] = ins.y.get(regs)
	case add:
		regs[ins.x.reg] += ins.y.get(regs)
	case mul:
		regs[ins.x.reg] *= ins.y.get(regs)
	case mod:
		regs[ins.x.reg] %= ins.y.get(regs)
	case jgz:
		if ins.x.get(regs) > 0 {
			m.pc += ins.y.get(regs)
			return
		}
	}

	m.pc++
}

// run executes instructions until the program ends or waits on an empty
// inbox. snd pushes a value onto outbox, and rcv pops one from inbox. It
// returns how many instructions ran, so the caller can see when both
// programs are stuck.
func (m *machine) run(prog []instruction, inbox, outbox *ds.Queue[int]) int {
	ran := 0

	for m.inside(prog) {
		ins := prog[m.pc]

		switch ins.op {
		case snd:
			outbox.Push(ins.x.get(&m.regs))
			m.sent++
			m.pc++
		case rcv:
			if inbox.Len() == 0 {
				return ran
			}

			m.regs[ins.x.reg] = inbox.Pop()
			m.pc++
		default:
			m.exec(ins)
		}

		ran++
	}

	return ran
}
