// Advent of Code 2022, day 20: Grove Positioning System.
// https://adventofcode.com/2022/day/20
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 20, part1, part2)
}

// decryptionKey is the number that part 2 multiplies each value by.
const decryptionKey = 811589153

// part1 mixes the file one time and adds the grove coordinates.
func part1(in string) any {
	return decrypt(input.IntLines(in), 1, 1)
}

// part2 multiplies each value by the decryption key, mixes the file ten times
// and adds the grove coordinates.
func part2(in string) any {
	return decrypt(input.IntLines(in), decryptionKey, 10)
}

// decrypt mixes the values and returns the sum of the values 1000, 2000 and
// 3000 places after the value 0.
//
// The file can have the same value more than one time. So the mix does not
// move values; it moves the original index of each value. order holds these
// indexes in their current sequence.
func decrypt(values []int, key, rounds int) int {
	n := len(values)

	for i := range values {
		values[i] *= key
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	for range rounds {
		for i, v := range values {
			from := slices.Index(order, i)

			// When a value moves, it is first taken out of the circle. The
			// other n-1 values stay, so one full turn is n-1 steps.
			to := mathx.Mod(from+v, n-1)

			// Slide the values between the two places by one step, then put
			// index i in the gap. copy handles slices that overlap.
			if to > from {
				copy(order[from:to], order[from+1:to+1])
			} else {
				copy(order[to+1:from+1], order[to:from])
			}

			order[to] = i
		}
	}

	zero := slices.Index(order, slices.Index(values, 0))
	sum := 0

	for _, offset := range []int{1000, 2000, 3000} {
		sum += values[order[(zero+offset)%n]]
	}

	return sum
}
