// Advent of Code 2016, day 6: Signals and Noise.
// https://adventofcode.com/2016/day/6
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 6, part1, part2)
}

// part1 takes the most common letter in each column.
func part1(in string) any {
	return decode(in, func(count, best int) bool { return count > best })
}

// part2 takes the least common letter in each column.
func part2(in string) any {
	return decode(in, func(count, best int) bool { return count < best })
}

// decode counts the letters in each column of the messages and keeps one
// letter per column. better reports whether a letter's count beats the best
// count so far. Letters that never appear in a column are skipped.
func decode(in string, better func(count, best int) bool) string {
	lines := input.Lines(in)
	message := make([]byte, len(lines[0]))

	for col := range message {
		var counts [26]int

		for _, line := range lines {
			counts[line[col]-'a']++
		}

		best := -1

		for i, count := range counts {
			if count == 0 {
				continue
			}

			if best == -1 || better(count, counts[best]) {
				best = i
			}
		}

		message[col] = 'a' + byte(best)
	}

	return string(message)
}
