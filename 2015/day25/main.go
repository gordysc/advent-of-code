// Advent of Code 2015, day 25: Let It Snow.
// https://adventofcode.com/2015/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 25, part1, nil)
}

// The code generator from the puzzle: it starts at firstCode and every later
// code is the previous one times multiplier, modulo modulus.
const (
	firstCode  = 20151125
	multiplier = 252533
	modulus    = 33554393
)

// part1 reads the row and column from the input and returns the code that
// belongs at that position in the machine's grid.
func part1(in string) any {
	nums := input.Ints(in)
	if len(nums) != 2 {
		panic("expected a row and a column in the input")
	}

	row, col := nums[0], nums[1]

	return codeAt(row, col)
}

// codeAt returns the code at a 1-based row and column. The grid is filled
// along diagonals that run up and to the right, so the position's index in the
// code sequence follows from which diagonal it sits on and how far along it is.
// The code itself is firstCode advanced index-1 steps, computed in one go with
// modular exponentiation instead of stepping through every earlier code.
func codeAt(row, col int) int {
	// Diagonal d holds the cells whose row+col-1 equals d, and the diagonals
	// before it hold 1+2+...+(d-1) cells in total.
	diagonal := row + col - 1
	index := diagonal*(diagonal-1)/2 + col

	return firstCode * powMod(multiplier, index-1, modulus) % modulus
}

// powMod returns base raised to exp, modulo m, by squaring. Every product
// stays below m*m, which fits comfortably in an int for this puzzle's modulus.
func powMod(base, exp, m int) int {
	result := 1
	base %= m

	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % m
		}

		base = base * base % m
		exp >>= 1
	}

	return result
}
