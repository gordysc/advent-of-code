// Advent of Code 2019, day 1: The Tyranny of the Rocket Equation.
// https://adventofcode.com/2019/day/1
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 1, part1, part2)
}

// part1 adds up the fuel that each module needs for its own mass.
func part1(in string) any {
	total := 0

	for _, mass := range input.IntLines(in) {
		total += fuel(mass)
	}

	return total
}

// part2 adds up the fuel for each module, and also the fuel that the fuel
// itself needs.
func part2(in string) any {
	total := 0

	for _, mass := range input.IntLines(in) {
		total += totalFuel(mass)
	}

	return total
}

// fuel returns the fuel for a mass: divide by three, round down, subtract 2.
// Go's integer division rounds down for positive numbers, so mass/3 does the
// rounding. Small masses give zero or a negative number.
func fuel(mass int) int {
	return mass/3 - 2
}

// totalFuel returns the fuel for a mass plus the fuel for that fuel, and so
// on. Each new amount of fuel is smaller than the last, so the loop stops when
// an amount needs no more fuel (zero or less).
func totalFuel(mass int) int {
	total := 0

	for f := fuel(mass); f > 0; f = fuel(f) {
		total += f
	}

	return total
}
