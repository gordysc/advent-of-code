// Advent of Code 2024, day 17: Chronospatial Computer.
// https://adventofcode.com/2024/day/17
package main

import (
	"slices"
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 17, part1, part2)
}

// part1 runs the program and returns its output, joined with commas.
//
// The puzzle has two examples. The file example.txt holds the part 2 example.
// On it, part 1 gives 5,7,3,0 and part 2 gives 117440. On the part 1 example,
// part 1 gives 4,6,3,5,6,3,5,2,1,0. That program shifts A by 1 bit each loop,
// not 3, so part 2 finds no answer and gives nil.
func part1(in string) any {
	c := parse(in)
	out := c.run(c.a)

	// strconv.Itoa converts each number to its decimal text.
	parts := make([]string, len(out))
	for i, v := range out {
		parts[i] = strconv.Itoa(v)
	}

	return strings.Join(parts, ",")
}

// part2 returns the lowest positive start value of register A that makes the
// program output a copy of itself.
//
// This works for programs like the real inputs: each loop outputs one digit,
// shifts A right by 3 bits, and sets B and C only from A. Then the last
// output depends only on the top 3 bits of A, the output before it on the top
// 6 bits, and so on. So we build A 3 bits at a time, from the last output back
// to the first, and go back when no 3 bits fit.
func part2(in string) any {
	c := parse(in)

	a, ok := c.findQuine(len(c.program)-1, 0)
	if !ok {
		return nil
	}

	return a
}

// computer holds the start values of the registers and the program.
type computer struct {
	a, b, c int
	program []int
}

// parse reads the three registers and the program. input.Ints finds every
// number in the text, so the first three are the registers.
func parse(in string) computer {
	nums := input.Ints(in)

	return computer{a: nums[0], b: nums[1], c: nums[2], program: nums[3:]}
}

// run executes the program with a as the start value of register A. It
// returns the values that the out instruction writes.
func (c computer) run(a int) []int {
	b, cr := c.b, c.c
	var out []int

	for ip := 0; ip+1 < len(c.program); ip += 2 {
		op := c.program[ip+1]

		// The combo operand maps 0-3 to themselves and 4-6 to the registers.
		combo := op
		switch op {
		case 4:
			combo = a
		case 5:
			combo = b
		case 6:
			combo = cr
		}

		// A division by 2^combo is the same as a right shift by combo bits.
		switch c.program[ip] {
		case 0: // adv
			a >>= combo
		case 1: // bxl
			b ^= op
		case 2: // bst
			b = combo % 8
		case 3: // jnz
			// Subtract 2 because the loop adds 2 after each instruction.
			if a != 0 {
				ip = op - 2
			}
		case 4: // bxc
			b ^= cr
		case 5: // out
			out = append(out, combo%8)
		case 6: // bdv
			b = a >> combo
		case 7: // cdv
			cr = a >> combo
		}
	}

	return out
}

// findQuine adds 3 more bits to the partial value a, so that the program
// outputs program[i:]. It then continues with i-1. It tries the 3 bits in
// increasing order, so the first full answer it finds is the lowest one.
func (c computer) findQuine(i, a int) (int, bool) {
	if i < 0 {
		return a, true
	}

	for bits := range 8 {
		next := a<<3 | bits

		// A is 0 when the program stops, so a start value of 0 is not useful.
		if next == 0 {
			continue
		}

		// slices.Equal compares the lengths and then each value in order.
		if !slices.Equal(c.run(next), c.program[i:]) {
			continue
		}

		if found, ok := c.findQuine(i-1, next); ok {
			return found, true
		}
	}

	return 0, false
}
