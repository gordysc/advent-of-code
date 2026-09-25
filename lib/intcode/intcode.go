// Package intcode runs the Intcode programs from Advent of Code 2019.
//
// A [Machine] holds one program's memory, its instruction pointer and its
// relative base. [Machine.Run] executes instructions until the program halts,
// or until it needs input that has not been sent yet. The caller can then send
// more input with [Machine.Send] and call Run again. This lets one solution
// drive a game loop, a robot, or a network of machines step by step.
//
// Like package input, the machine panics on bad data (an unknown opcode or
// parameter mode) instead of returning an error. Puzzle programs are trusted.
package intcode

import (
	"fmt"
	"strings"

	"aoc/lib/input"
)

// Opcodes of the Intcode instruction set.
const (
	opAdd       = 1
	opMul       = 2
	opIn        = 3
	opOut       = 4
	opJumpTrue  = 5
	opJumpFalse = 6
	opLess      = 7
	opEqual     = 8
	opBase      = 9
	opHalt      = 99
)

// Parameter modes. A position parameter is an address, an immediate parameter
// is the value itself, and a relative parameter is an address from the
// relative base.
const (
	modePosition  = 0
	modeImmediate = 1
	modeRelative  = 2
)

// Machine is one Intcode computer. Its methods have pointer receivers, so
// they change this machine and not a copy of it. mem is the program's memory,
// ip is the address of the next instruction, base is the relative base, and
// input holds the values sent but not read yet.
type Machine struct {
	mem    []int
	ip     int
	base   int
	input  []int
	halted bool
}

// Parse reads a comma-separated Intcode program.
func Parse(s string) []int {
	return input.Ints(strings.TrimSpace(s))
}

// New returns a machine that runs a copy of program, so the caller's slice
// never changes.
func New(program []int) *Machine {
	return &Machine{mem: append([]int(nil), program...)}
}

// Clone returns an independent copy of the machine, with the same memory,
// position and pending input.
func (m *Machine) Clone() *Machine {
	// Copying the struct copies the slice headers, not the arrays under them,
	// so the copy would still share memory with m. Appending to a nil slice
	// makes a new array for each one.
	c := *m
	c.mem = append([]int(nil), m.mem...)
	c.input = append([]int(nil), m.input...)

	return &c
}

// Peek returns the value at an address. Memory past the program reads as 0.
func (m *Machine) Peek(addr int) int {
	if addr < 0 {
		panic(fmt.Sprintf("intcode: negative address %d", addr))
	}

	if addr >= len(m.mem) {
		return 0
	}

	return m.mem[addr]
}

// Poke writes a value to an address, and grows memory when needed.
func (m *Machine) Poke(addr, value int) {
	if addr < 0 {
		panic(fmt.Sprintf("intcode: negative address %d", addr))
	}

	if addr >= len(m.mem) {
		m.mem = append(m.mem, make([]int, addr+1-len(m.mem))...)
	}

	m.mem[addr] = value
}

// Send queues values for the program's input instructions to read.
func (m *Machine) Send(values ...int) {
	m.input = append(m.input, values...)
}

// SendLine queues a line of ASCII text followed by a newline, for the
// programs that take text commands.
func (m *Machine) SendLine(s string) {
	for _, r := range s {
		m.input = append(m.input, int(r))
	}

	m.input = append(m.input, '\n')
}

// Halted reports whether the program has run its halt instruction.
func (m *Machine) Halted() bool {
	return m.halted
}

// Waiting reports whether the program stopped because it needs input that
// has not been sent yet.
func (m *Machine) Waiting() bool {
	return !m.halted && len(m.input) == 0 && m.Peek(m.ip)%100 == opIn
}

// Run executes instructions until the program halts or needs input it does
// not have. It returns the values the program wrote in that time, oldest
// first. Calling Run on a halted machine returns nothing.
func (m *Machine) Run() []int {
	var out []int

	for !m.halted {
		instr := m.Peek(m.ip)
		op := instr % 100

		switch op {
		case opAdd:
			m.Poke(m.addr(instr, 3), m.param(instr, 1)+m.param(instr, 2))
			m.ip += 4

		case opMul:
			m.Poke(m.addr(instr, 3), m.param(instr, 1)*m.param(instr, 2))
			m.ip += 4

		case opIn:
			// With no input, stop here without moving ip. The next call to
			// Run starts on this same instruction and reads the new input.
			if len(m.input) == 0 {
				return out
			}

			m.Poke(m.addr(instr, 1), m.input[0])
			m.input = m.input[1:]
			m.ip += 2

		case opOut:
			out = append(out, m.param(instr, 1))
			m.ip += 2

		case opJumpTrue:
			if m.param(instr, 1) != 0 {
				m.ip = m.param(instr, 2)
			} else {
				m.ip += 3
			}

		case opJumpFalse:
			if m.param(instr, 1) == 0 {
				m.ip = m.param(instr, 2)
			} else {
				m.ip += 3
			}

		case opLess:
			m.Poke(m.addr(instr, 3), boolInt(m.param(instr, 1) < m.param(instr, 2)))
			m.ip += 4

		case opEqual:
			m.Poke(m.addr(instr, 3), boolInt(m.param(instr, 1) == m.param(instr, 2)))
			m.ip += 4

		case opBase:
			m.base += m.param(instr, 1)
			m.ip += 2

		case opHalt:
			m.halted = true

		default:
			panic(fmt.Sprintf("intcode: unknown opcode %d at %d", op, m.ip))
		}
	}

	return out
}

// RunString runs the program like [Machine.Run] and returns its output as
// ASCII text. It is for the programs that print a picture or a message.
func (m *Machine) RunString() string {
	var sb strings.Builder
	for _, v := range m.Run() {
		sb.WriteRune(rune(v))
	}

	return sb.String()
}

// mode returns the mode of parameter n (counted from 1) of an instruction.
// The modes are the digits above the two opcode digits, read right to left.
func mode(instr, n int) int {
	div := 100
	for range n - 1 {
		div *= 10
	}

	return instr / div % 10
}

// param returns the value of parameter n of the current instruction.
func (m *Machine) param(instr, n int) int {
	raw := m.Peek(m.ip + n)

	switch mode(instr, n) {
	case modePosition:
		return m.Peek(raw)
	case modeImmediate:
		return raw
	case modeRelative:
		return m.Peek(m.base + raw)
	}

	panic(fmt.Sprintf("intcode: bad mode in %d at %d", instr, m.ip))
}

// addr returns the address that parameter n of the current instruction
// writes to. A write parameter is never in immediate mode.
func (m *Machine) addr(instr, n int) int {
	raw := m.Peek(m.ip + n)

	switch mode(instr, n) {
	case modePosition:
		return raw
	case modeRelative:
		return m.base + raw
	}

	panic(fmt.Sprintf("intcode: bad write mode in %d at %d", instr, m.ip))
}

// boolInt turns true into 1 and false into 0.
func boolInt(b bool) int {
	if b {
		return 1
	}

	return 0
}
