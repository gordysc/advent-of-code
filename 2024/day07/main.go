// Advent of Code 2024, day 7: Bridge Repair.
// https://adventofcode.com/2024/day/7
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 7, part1, part2)
}

// part1 adds the test values of the equations that add (+) and multiply (*)
// can make true.
func part1(in string) any {
	return calibration(in, false)
}

// part2 adds the test values of the equations that add, multiply, and
// concatenate (||) can make true.
func part2(in string) any {
	return calibration(in, true)
}

// calibration adds the test value of every equation that can be true. concat
// allows the || operator.
func calibration(in string, concat bool) int {
	total := 0

	for _, line := range input.Lines(in) {
		// UInts also finds the test value before the colon, so it comes first.
		nums := input.UInts(line)
		if len(nums) < 2 {
			continue
		}

		if solvable(nums[0], nums[1:], concat) {
			total += nums[0]
		}
	}

	return total
}

// solvable reports whether the operators can join nums, from left to right,
// to give target.
//
// It works backward from the last number. The last operator must undo
// cleanly: a sum must leave a remainder that is not negative, a product must
// divide exactly, and a concatenation must end in the digits of the last
// number. Most branches fail these tests at once, so this is much faster than
// trying every operator combination from the front.
func solvable(target int, nums []int, concat bool) bool {
	last := nums[len(nums)-1]
	rest := nums[:len(nums)-1]

	if len(rest) == 0 {
		return target == last
	}

	if target >= last && solvable(target-last, rest, concat) {
		return true
	}

	if last != 0 && target%last == 0 && solvable(target/last, rest, concat) {
		return true
	}

	if concat {
		// For last = 45, shift is 100, and 12345 splits into 123 and 45.
		shift := mathx.Pow(10, mathx.NumDigits(last))
		if target%shift == last && solvable(target/shift, rest, concat) {
			return true
		}
	}

	return false
}
