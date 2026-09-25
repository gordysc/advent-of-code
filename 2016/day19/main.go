// Advent of Code 2016, day 19: An Elephant Named Joseph.
// https://adventofcode.com/2016/day/19
package main

import (
	"math/bits"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 19, part1, part2)
}

// part1 finds the elf that gets all the presents when each elf takes from the
// elf on its left.
//
// This is the Josephus problem with every second elf out. Let p be the
// largest power of 2 that is not more than n. When n is exactly p, each lap
// removes half of the elves and elf 1 always goes first, so elf 1 wins. With
// n = p + k, after k elves are out it is the turn of elf 2k+1, and p elves
// are left. So elf 2k+1 wins.
func part1(in string) any {
	n := input.Int(strings.TrimSpace(in))
	p := 1 << (bits.Len(uint(n)) - 1)

	return 2*(n-p) + 1
}

// part2 finds the winner when each elf takes from the elf across the circle.
//
// A simulation for small n shows a pattern based on powers of 3. Let p be the
// largest power of 3 that is less than n (for n = 1, the winner is elf 1).
// For n up to 2p the winner goes up by one for each extra elf, so it is
// n - p. After that it goes up by two for each extra elf, so it is 2n - 3p.
// At n = 3p the winner is n, and the next power of 3 starts the count again.
func part2(in string) any {
	n := input.Int(strings.TrimSpace(in))
	if n == 1 {
		return 1
	}

	p := 1
	for p*3 < n {
		p *= 3
	}

	if n <= 2*p {
		return n - p
	}

	return 2*n - 3*p
}
