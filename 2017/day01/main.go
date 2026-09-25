// Advent of Code 2017, day 1: Inverse Captcha.
// https://adventofcode.com/2017/day/1
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 1, part1, part2)
}

// part1 sums every digit that matches the next digit. The list is circular,
// so the last digit is compared with the first.
func part1(in string) any {
	return captcha(input.Digits(in), 1)
}

// part2 sums every digit that matches the digit halfway around the list.
func part2(in string) any {
	digits := input.Digits(in)

	return captcha(digits, len(digits)/2)
}

// captcha sums every digit that matches the digit offset places ahead of it,
// wrapping around the end of the list.
func captcha(digits []int, offset int) int {
	sum := 0

	for i, d := range digits {
		if d == digits[(i+offset)%len(digits)] {
			sum += d
		}
	}

	return sum
}
