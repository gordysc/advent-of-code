// Advent of Code 2022, day 6: Tuning Trouble.
// https://adventofcode.com/2022/day/6
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 6, part1, part2)
}

// part1 finds the end of the first start-of-packet marker: four characters
// that are all different.
//
// The puzzle gives five examples. example.txt has only the first one, which
// gives 7 for part 1 and 19 for part 2.
func part1(in string) any {
	return findMarker(strings.TrimSpace(in), 4)
}

// part2 finds the end of the first start-of-message marker: fourteen
// characters that are all different.
func part2(in string) any {
	return findMarker(strings.TrimSpace(in), 14)
}

// findMarker gives the number of characters read when the last size
// characters are all different for the first time.
//
// A sliding window keeps a count for each letter in the window, and the
// number of letters that occur more than once. Each step adds one letter and
// removes one letter, so the full search is linear in the length of the
// signal.
func findMarker(signal string, size int) int {
	var counts [256]int
	repeats := 0

	for i := 0; i < len(signal); i++ {
		// Add the new letter to the window.
		counts[signal[i]]++

		if counts[signal[i]] == 2 {
			repeats++
		}

		// Remove the letter that falls out of the window.
		if i >= size {
			old := signal[i-size]
			counts[old]--

			if counts[old] == 1 {
				repeats--
			}
		}

		if i >= size-1 && repeats == 0 {
			return i + 1
		}
	}

	return -1
}
