// Advent of Code 2015, day 8: Matchsticks.
// https://adventofcode.com/2015/day/8
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 8, part1, part2)
}

// part1 sums, over every line, how many more characters the string literal
// takes up in the file than the string it represents takes up in memory.
func part1(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		total += len(line) - decodedLen(line)
	}

	return total
}

// part2 goes the other way: it sums how many more characters each line would
// take up if it were itself encoded as a string literal.
func part2(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		total += encodedLen(line) - len(line)
	}

	return total
}

// decodedLen counts the characters a literal represents in memory. The
// surrounding quotes count for nothing, and each escape sequence counts as one
// character: \\ and \" are two characters of literal, \xHH is four.
func decodedLen(literal string) int {
	body := literal[1 : len(literal)-1]
	n := 0

	// i jumps past each escape sequence in one step, so the loop cannot use
	// range, which would visit every byte.
	for i := 0; i < len(body); i++ {
		if body[i] == '\\' {
			if body[i+1] == 'x' {
				i += 3
			} else {
				i++
			}
		}

		n++
	}

	return n
}

// encodedLen counts the characters needed to write the line as a literal.
// Every character stays as it is except " and \, which each gain a backslash,
// and the whole thing is wrapped in a new pair of quotes.
func encodedLen(line string) int {
	return len(line) + strings.Count(line, `"`) + strings.Count(line, `\`) + 2
}
