// Advent of Code 2021, day 24: Arithmetic Logic Unit.
// https://adventofcode.com/2021/day/24
package main

import (
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 24, part1, part2)
}

// part1 finds the largest model number that MONAD accepts.
//
// solve pairs the blocks of the program and finds the digits. The ALU then
// runs the program on the number to make sure that z is 0.
func part1(in string) any {
	return solve(in, true)
}

// part2 finds the smallest model number that MONAD accepts.
//
// It uses the same method as part 1, but picks the smallest digits.
func part2(in string) any {
	return solve(in, false)
}

// block holds the three constants that make one of the 14 MONAD blocks
// different from the others.
type block struct {
	div    int // the operand of "div z": 1 or 26
	checkX int // the operand of "add x" after "div z"
	addY   int // the operand of "add y" after "add y w"
}

// solve returns the largest (or smallest) accepted model number, or nil if
// the ALU does not accept it.
//
// MONAD uses z as a stack of base-26 digits. A block with "div z 1" always
// pushes w + addY, because its checkX is 10 or more and x can never equal w.
// A block with "div z 26" pops the top value v. It pushes again unless
// w = v + checkX. For z to be 0 at the end, no pop block can push. So each
// pop block j pairs with the push block i that pushed its value, and the
// digits must obey w[j] = w[i] + addY[i] + checkX[j].
func solve(in string, largest bool) any {
	blocks := parseBlocks(in)
	digits := make([]int, len(blocks))

	type pushed struct {
		index, addY int
	}

	var stack []pushed

	for j, b := range blocks {
		if b.div == 1 {
			stack = append(stack, pushed{j, b.addY})
			continue
		}

		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		i := top.index
		delta := top.addY + b.checkX

		// Pick w[i] so that both w[i] and w[j] = w[i] + delta are digits
		// from 1 to 9.
		if largest {
			digits[i] = min(9, 9-delta)
		} else {
			digits[i] = max(1, 1-delta)
		}

		digits[j] = digits[i] + delta
	}

	// Run the program on the digits to make sure that the ALU accepts them.
	if run(in, digits)['z'-'w'] != 0 {
		return nil
	}

	var sb strings.Builder
	for _, d := range digits {
		sb.WriteByte(byte('0' + d))
	}

	return sb.String()
}

// parseBlocks reads the constants of each block. A block starts at each
// "inp w" line. The constants are on lines 4, 5 and 15 of the block.
func parseBlocks(in string) []block {
	lines := input.Lines(in)

	var blocks []block
	for start := 0; start < len(lines); start += 18 {
		arg := func(offset int) int {
			f := strings.Fields(lines[start+offset])
			return input.Int(f[2])
		}

		blocks = append(blocks, block{div: arg(4), checkX: arg(5), addY: arg(15)})
	}

	return blocks
}

// run executes the ALU program with the given input digits. It returns the
// registers w, x, y and z, in that order.
func run(in string, inputs []int) [4]int {
	var reg [4]int

	// value reads an operand: a register name or a number.
	value := func(s string) int {
		if s >= "w" && s <= "z" && len(s) == 1 {
			return reg[s[0]-'w']
		}

		n, _ := strconv.Atoi(s)
		return n
	}

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)
		a := f[1][0] - 'w'

		if f[0] == "inp" {
			reg[a] = inputs[0]
			inputs = inputs[1:]
			continue
		}

		b := value(f[2])

		switch f[0] {
		case "add":
			reg[a] += b
		case "mul":
			reg[a] *= b
		case "div":
			reg[a] /= b
		case "mod":
			reg[a] %= b
		case "eql":
			if reg[a] == b {
				reg[a] = 1
			} else {
				reg[a] = 0
			}
		}
	}

	return reg
}
