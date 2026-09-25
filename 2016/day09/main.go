// Advent of Code 2016, day 9: Explosives in Cyberspace.
// https://adventofcode.com/2016/day/9
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 9, part1, part2)
}

// part1 finds the decompressed length when markers inside repeated data are
// copied as plain text.
func part1(in string) any {
	return length(compact(in), false)
}

// part2 finds the decompressed length when markers inside repeated data are
// expanded too.
func part2(in string) any {
	return length(compact(in), true)
}

// compact removes the whitespace, which the puzzle says to ignore.
func compact(in string) string {
	return strings.Join(strings.Fields(in), "")
}

// length returns the decompressed length of s without building the text.
// A marker (AxB) repeats the next A characters B times. When recurse is true,
// the markers in those A characters are expanded as well.
func length(s string, recurse bool) int {
	total := 0

	for i := 0; i < len(s); {
		if s[i] != '(' {
			total++
			i++
			continue
		}

		end := i + strings.IndexByte(s[i:], ')')
		nums := input.Ints(s[i+1 : end])
		span, times := nums[0], nums[1]

		data := s[end+1 : end+1+span]
		size := len(data)
		if recurse {
			size = length(data, true)
		}

		total += size * times
		i = end + 1 + span
	}

	return total
}
