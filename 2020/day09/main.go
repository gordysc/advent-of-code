// Advent of Code 2020, day 9: Encoding Error.
// https://adventofcode.com/2020/day/9
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 9, part1, part2)
}

// part1 finds the first number that is not the sum of two different numbers
// in the preamble just before it.
func part1(in string) any {
	return invalid(input.IntLines(in))
}

// part2 finds a contiguous range of at least two numbers that adds up to the
// invalid number from part 1, and returns the sum of the smallest and largest
// numbers in that range. All inputs are positive, so a sliding window works:
// grow the window on the right while the sum is too small, and shrink it on
// the left while the sum is too large.
func part2(in string) any {
	nums := input.IntLines(in)
	want := invalid(nums)

	lo, sum := 0, 0
	for hi, n := range nums {
		sum += n
		for sum > want && lo < hi {
			sum -= nums[lo]
			lo++
		}

		if sum == want && hi > lo {
			small, large := nums[lo], nums[lo]
			for _, m := range nums[lo : hi+1] {
				small = min(small, m)
				large = max(large, m)
			}

			return small + large
		}
	}

	return nil
}

// invalid returns the first number that is not the sum of two different
// numbers in the preamble before it. The real input uses a preamble of 25.
// The example uses a preamble of 5 and has only 20 numbers, so a short input
// selects the example preamble.
func invalid(nums []int) int {
	preamble := 25
	if len(nums) <= 25 {
		preamble = 5
	}

	for i := preamble; i < len(nums); i++ {
		if !pairSum(nums[i-preamble:i], nums[i]) {
			return nums[i]
		}
	}

	return 0
}

// pairSum reports whether two numbers with different values in window add up
// to want. The window is small, so a check of every pair is fast.
func pairSum(window []int, want int) bool {
	for i, a := range window {
		for _, b := range window[i+1:] {
			if a != b && a+b == want {
				return true
			}
		}
	}

	return false
}
