// Advent of Code 2020, day 2: Password Philosophy.
// https://adventofcode.com/2020/day/2
package main

import (
	"fmt"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// entry is one line of the database: a policy and the password it applies to.
// The meaning of lo and hi changes between the two parts.
type entry struct {
	lo, hi   int
	letter   byte
	password string
}

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 2, part1, part2)
}

// part1 counts the passwords where the letter occurs between lo and hi
// times, inclusive.
func part1(in string) any {
	valid := 0
	for _, e := range parse(in) {
		n := strings.Count(e.password, string(e.letter))
		if n >= e.lo && n <= e.hi {
			valid++
		}
	}

	return valid
}

// part2 counts the passwords where the letter is at exactly one of the
// positions lo and hi. The positions start at 1, not 0.
func part2(in string) any {
	valid := 0
	for _, e := range parse(in) {
		first := e.password[e.lo-1] == e.letter
		second := e.password[e.hi-1] == e.letter

		if first != second {
			valid++
		}
	}

	return valid
}

// parse reads the "1-3 a: abcde" lines into entries.
func parse(in string) []entry {
	var entries []entry

	for _, line := range input.Lines(in) {
		var e entry

		_, err := fmt.Sscanf(line, "%d-%d %c: %s", &e.lo, &e.hi, &e.letter, &e.password)
		if err != nil {
			panic(fmt.Sprintf("bad line %q: %v", line, err))
		}

		entries = append(entries, e)
	}

	return entries
}
