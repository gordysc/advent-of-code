// Advent of Code 2015, day 5: Doesn't He Have Intern-Elves For This?
// https://adventofcode.com/2015/day/5
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 5, part1, part2)
}

// part1 counts the strings that pass the first set of rules.
func part1(in string) any {
	return slicesx.Count(input.Lines(in), nice1)
}

// part2 counts the strings that pass the second set of rules.
func part2(in string) any {
	return slicesx.Count(input.Lines(in), nice2)
}

// nice1 reports whether s has at least three vowels, a letter that appears
// twice in a row, and none of the forbidden pairs ab, cd, pq or xy.
func nice1(s string) bool {
	vowels := 0
	double := false

	for i := 0; i < len(s); i++ {
		if strings.IndexByte("aeiou", s[i]) >= 0 {
			vowels++
		}

		if i == 0 {
			continue
		}

		pair := s[i-1 : i+1]
		switch pair {
		case "ab", "cd", "pq", "xy":
			return false
		}

		if s[i-1] == s[i] {
			double = true
		}
	}

	return vowels >= 3 && double
}

// nice2 reports whether s has a pair of letters that appears twice without
// overlapping, and a letter that repeats with exactly one letter in between.
func nice2(s string) bool {
	return hasRepeatedPair(s) && hasSandwich(s)
}

// hasRepeatedPair looks for a two-letter pair that appears at least twice with
// no overlap. Overlap only matters for a triple such as "aaa", so the check
// searches for each pair again starting two positions later.
func hasRepeatedPair(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if strings.Contains(s[i+2:], s[i:i+2]) {
			return true
		}
	}

	return false
}

// hasSandwich looks for a pattern like "aba" or "xxx": the same letter at
// positions i and i+2.
func hasSandwich(s string) bool {
	for i := 0; i+2 < len(s); i++ {
		if s[i] == s[i+2] {
			return true
		}
	}

	return false
}
