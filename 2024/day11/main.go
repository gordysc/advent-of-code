// Advent of Code 2024, day 11: Plutonian Pebbles.
// https://adventofcode.com/2024/day/11
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 11, part1, part2)
}

// part1 counts the stones after 25 blinks.
func part1(in string) any {
	return blink(in, 25)
}

// part2 counts the stones after 75 blinks.
//
// The puzzle text gives no example answer for part 2. On the example, this
// part gives 65601038650482.
func part2(in string) any {
	return blink(in, 75)
}

// blink counts the stones after the given number of blinks.
//
// The order of the stones has no effect on how each stone changes, and the
// answer is only a count. So we keep a count for each engraved number, not a
// list of stones. After a few blinks, most stones share a small set of
// numbers, so the map stays small while the stone count grows very fast.
func blink(in string, times int) int {
	counts := map[int]int{}
	for _, n := range input.UInts(in) {
		counts[n]++
	}

	for range times {
		next := make(map[int]int, len(counts))

		for n, c := range counts {
			for _, m := range change(n) {
				next[m] += c
			}
		}

		counts = next
	}

	total := 0
	for _, c := range counts {
		total += c
	}

	return total
}

// change applies the first rule that fits to one stone, and returns the
// stones that replace it.
func change(n int) []int {
	if n == 0 {
		return []int{1}
	}

	digits := mathx.NumDigits(n)
	if digits%2 == 0 {
		// Split the digits in two halves. The divisor has one zero for each
		// digit in the right half, so the quotient is the left half and the
		// remainder is the right half. Leading zeros go away by themselves.
		div := mathx.Pow(10, digits/2)

		return []int{n / div, n % div}
	}

	return []int{n * 2024}
}
