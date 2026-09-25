// Advent of Code 2019, day 4: Secure Container.
// https://adventofcode.com/2019/day/4
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 4, part1, part2)
}

// The puzzle text gives no example range, only single passwords to check, so
// example.txt holds the small range 111110-111125, checked by hand:
//
//   - 111110, 111120 and 111121 have a digit that goes down, so they fail.
//   - 111111 to 111119 and 111122 to 111125 pass part 1. That is 13.
//   - For part 2, the 1s in all of these form a group of four or more.
//     Only 111122 also has a group of exactly two (the 22), so part 2 gives 1.

// part1 counts the passwords in the range that never decrease and have at
// least two equal digits next to each other.
func part1(in string) any {
	return count(in, func(run int) bool { return run >= 2 })
}

// part2 counts the passwords that never decrease and have a group of exactly
// two equal digits. A longer group, such as 444, does not count as a pair.
func part2(in string) any {
	return count(in, func(run int) bool { return run == 2 })
}

// count checks every number in the input range "lo-hi". input.UInts reads the
// numbers without a minus sign, because the dash here is a separator.
// pairOK decides which group lengths satisfy the pair rule.
func count(in string, pairOK func(run int) bool) int {
	bounds := input.UInts(in)
	n := 0

	for pw := bounds[0]; pw <= bounds[1]; pw++ {
		if valid(pw, pairOK) {
			n++
		}
	}

	return n
}

// valid reports whether the digits of pw never decrease from left to right
// and at least one group of equal digits has a length that pairOK accepts.
//
// It reads the digits from the right with pw%10 and pw/10, so from this side
// each digit must be less than or equal to the one before it. Equal digits
// next to each other form a group, and run counts the length of the group so
// far. A group ends when the digit changes, or after the last digit.
func valid(pw int, pairOK func(run int) bool) bool {
	prev := pw % 10
	pw /= 10
	run := 1
	pair := false

	for pw > 0 {
		d := pw % 10
		pw /= 10

		if d > prev {
			return false
		}

		if d == prev {
			run++
			continue
		}

		pair = pair || pairOK(run)
		prev = d
		run = 1
	}

	return pair || pairOK(run)
}
