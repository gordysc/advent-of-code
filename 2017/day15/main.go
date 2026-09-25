// Advent of Code 2017, day 15: Dueling Generators.
// https://adventofcode.com/2017/day/15
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 15, part1, part2)
}

// The generators multiply their last value by a factor and keep the remainder
// after dividing by modulus.
const (
	factorA = 16807
	factorB = 48271
	modulus = 2147483647
)

// The number of pairs the judge compares in each part.
const (
	pairs1 = 40_000_000
	pairs2 = 5_000_000
)

// part1 counts the pairs where the lowest 16 bits of both values match.
func part1(in string) any {
	a, b := parse(in)
	count := 0

	for range pairs1 {
		a = a * factorA % modulus
		b = b * factorB % modulus

		if a&0xffff == b&0xffff {
			count++
		}
	}

	return count
}

// part2 counts matching pairs again, but generator A only hands over multiples
// of 4 and generator B only hands over multiples of 8.
func part2(in string) any {
	a, b := parse(in)
	count := 0

	for range pairs2 {
		a = next(a, factorA, 4)
		b = next(b, factorB, 8)

		if a&0xffff == b&0xffff {
			count++
		}
	}

	return count
}

// next generates values from prev until it finds a multiple of mult.
func next(prev, factor, mult int) int {
	for {
		prev = prev * factor % modulus
		if prev%mult == 0 {
			return prev
		}
	}
}

// parse reads the two starting values.
func parse(in string) (a, b int) {
	nums := input.Ints(in)

	return nums[0], nums[1]
}
