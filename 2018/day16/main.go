// Advent of Code 2018, day 16: Chronal Classification.
// https://adventofcode.com/2018/day/16
package main

import (
	"math/bits"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 16, part1, part2)
}

// registers is the device's set of four registers.
type registers [4]int

// op computes the value an opcode writes to register C, given inputs A and B.
type op func(r registers, a, b int) int

// ops holds the sixteen opcodes. Their order here is only a fixed index; the
// numbers the device really uses are worked out in part 2.
var ops = [16]op{
	func(r registers, a, b int) int { return r[a] + r[b] },           // addr
	func(r registers, a, b int) int { return r[a] + b },              // addi
	func(r registers, a, b int) int { return r[a] * r[b] },           // mulr
	func(r registers, a, b int) int { return r[a] * b },              // muli
	func(r registers, a, b int) int { return r[a] & r[b] },           // banr
	func(r registers, a, b int) int { return r[a] & b },              // bani
	func(r registers, a, b int) int { return r[a] | r[b] },           // borr
	func(r registers, a, b int) int { return r[a] | b },              // bori
	func(r registers, a, b int) int { return r[a] },                  // setr
	func(r registers, a, b int) int { return a },                     // seti
	func(r registers, a, b int) int { return boolInt(a > r[b]) },     // gtir
	func(r registers, a, b int) int { return boolInt(r[a] > b) },     // gtri
	func(r registers, a, b int) int { return boolInt(r[a] > r[b]) },  // gtrr
	func(r registers, a, b int) int { return boolInt(a == r[b]) },    // eqir
	func(r registers, a, b int) int { return boolInt(r[a] == b) },    // eqri
	func(r registers, a, b int) int { return boolInt(r[a] == r[b]) }, // eqrr
}

// sample is one recorded instruction with the registers before and after it.
type sample struct {
	before, after registers
	ins           [4]int
}

// part1 counts the samples that behave like three or more opcodes.
func part1(in string) any {
	samples, _ := parse(in)

	count := 0

	for _, s := range samples {
		// OnesCount16 counts the set bits, which is the number of opcodes.
		if bits.OnesCount16(matches(s)) >= 3 {
			count++
		}
	}

	return count
}

// part2 works out which number belongs to which opcode, then runs the test
// program and returns register 0. It has no answer when there is no program,
// as in the example, or when the samples do not pin down every number.
func part2(in string) any {
	samples, program := parse(in)
	if len(program) == 0 {
		return nil
	}

	table, ok := resolve(samples)
	if !ok {
		return nil
	}

	var r registers

	for _, ins := range program {
		r[ins[3]] = ops[table[ins[0]]](r, ins[1], ins[2])
	}

	return r[0]
}

// parse splits the input into the samples and the test program. Three blank
// lines separate the two sections; the example has only samples.
func parse(in string) ([]sample, [][4]int) {
	// Cut splits at the first separator. When there is none, tail is empty.
	head, tail, _ := strings.Cut(in, "\n\n\n\n")

	var samples []sample

	for _, block := range input.Blocks(head) {
		var s sample
		copy(s.before[:], input.Ints(block[0]))
		copy(s.ins[:], input.Ints(block[1]))
		copy(s.after[:], input.Ints(block[2]))
		samples = append(samples, s)
	}

	var program [][4]int

	for _, line := range input.Lines(tail) {
		var ins [4]int
		copy(ins[:], input.Ints(line))
		program = append(program, ins)
	}

	return samples, program
}

// matches returns a bit set of the opcodes that turn the sample's before
// registers into its after registers.
func matches(s sample) uint16 {
	var set uint16

	for i, f := range ops {
		// Arrays are values in Go, so this copy leaves s.before unchanged.
		r := s.before
		r[s.ins[3]] = f(r, s.ins[1], s.ins[2])

		if r == s.after {
			set |= 1 << i
		}
	}

	return set
}

// resolve maps each opcode number to an index in ops. Every sample narrows
// the choices for its number. Then, again and again, a number with only one
// choice left claims that opcode, and the other numbers lose it.
func resolve(samples []sample) ([16]int, bool) {
	var options [16]uint16

	for i := range options {
		options[i] = 0xffff
	}

	for _, s := range samples {
		options[s.ins[0]] &= matches(s)
	}

	var table [16]int

	for range 16 {
		n := -1

		for i, set := range options {
			if bits.OnesCount16(set) == 1 {
				n = i
				break
			}
		}

		if n < 0 {
			return table, false
		}

		// The set has one bit, and its position is the opcode's index.
		claimed := options[n]
		table[n] = bits.TrailingZeros16(claimed)

		// &^ is "and not": it clears the claimed bit in every set.
		for i := range options {
			options[i] &^= claimed
		}
	}

	return table, true
}

// boolInt turns a comparison result into the 1 or 0 the device stores.
func boolInt(b bool) int {
	if b {
		return 1
	}

	return 0
}
