// Advent of Code 2015, day 10: Elves Look, Elves Say.
// https://adventofcode.com/2015/day/10
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 10, part1, part2)
}

// part1 applies the look-and-say step 40 times and reports the length.
func part1(in string) any {
	return len(repeat(strings.TrimSpace(in), 40))
}

// part2 does the same for 50 steps. The sequence grows by roughly 30% each
// step, so this is about four times as much work as part 1.
func part2(in string) any {
	return len(repeat(strings.TrimSpace(in), 50))
}

// repeat applies lookAndSay to s the given number of times.
func repeat(s string, times int) string {
	for range times {
		s = lookAndSay(s)
	}

	return s
}

// lookAndSay reads the digits aloud: each run of the same digit becomes the
// run length followed by the digit, so "111221" becomes "312211".
func lookAndSay(s string) string {
	// The result is at most twice the input, so reserving that up front
	// avoids repeated growth as the builder is filled.
	out := make([]byte, 0, 2*len(s))

	for i := 0; i < len(s); {
		digit := s[i]
		j := i

		for j < len(s) && s[j] == digit {
			j++
		}

		// After the first step a run is never longer than three, and the
		// puzzle inputs have no long runs either, so one digit holds the count.
		out = append(out, byte('0'+j-i), digit)
		i = j
	}

	return string(out)
}
