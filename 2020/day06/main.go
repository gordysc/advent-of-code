// Advent of Code 2020, day 6: Custom Customs.
// https://adventofcode.com/2020/day/6
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 6, part1, part2)
}

// part1 adds up, for each group, the number of questions that anyone in the
// group answered "yes" to.
func part1(in string) any {
	total := 0
	for _, group := range input.Blocks(in) {
		for _, n := range counts(group) {
			if n > 0 {
				total++
			}
		}
	}

	return total
}

// part2 adds up, for each group, the number of questions that everyone in
// the group answered "yes" to. A question qualifies when its count equals
// the number of people in the group.
func part2(in string) any {
	total := 0
	for _, group := range input.Blocks(in) {
		for _, n := range counts(group) {
			if n == len(group) {
				total++
			}
		}
	}

	return total
}

// counts tells how many people in a group answered "yes" to each question.
// Each line is one person, and each letter a to z is one question.
func counts(group []string) [26]int {
	var c [26]int
	for _, person := range group {
		for _, q := range person {
			c[q-'a']++
		}
	}

	return c
}
