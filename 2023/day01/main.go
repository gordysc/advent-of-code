// Advent of Code 2023, day 1: Trebuchet?!.
// https://adventofcode.com/2023/day/1
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 1, part1, part2)
}

// part1 adds the calibration values, where each value uses only the digit
// characters in a line.
//
// The puzzle has two examples. The file example.txt holds the part 2 example,
// because the part 1 example gives no useful answer with the part 2 rules.
// Some lines of the part 2 example have no digit character, so part 1 skips
// those lines. On this example, part 1 gives 209 and part 2 gives 281.
func part1(in string) any {
	return calibrationSum(in, false)
}

// part2 adds the calibration values, where each value also accepts digits
// spelled out as words, such as "one" or "seven".
func part2(in string) any {
	return calibrationSum(in, true)
}

// calibrationSum adds the calibration value of every line. The value is the
// first digit times ten plus the last digit. A line with no digit adds
// nothing.
func calibrationSum(in string, words bool) int {
	total := 0

	for _, line := range input.Lines(in) {
		digits := lineDigits(line, words)
		if len(digits) == 0 {
			continue
		}

		total += digits[0]*10 + digits[len(digits)-1]
	}

	return total
}

// digitWords holds the spelled-out digits. The value of each word is its
// index plus one. The puzzle does not count "zero" as a digit word.
var digitWords = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

// lineDigits finds the digits in a line from left to right. When words is
// true, it also finds spelled-out digits.
//
// Words can overlap, as in "eightwo", which holds 8 and then 2. A simple
// search-and-replace would destroy one of the two words. So the scan tests
// every position on its own and never skips the letters of a word it found.
func lineDigits(line string, words bool) []int {
	var digits []int

	for i := 0; i < len(line); i++ {
		if line[i] >= '0' && line[i] <= '9' {
			digits = append(digits, int(line[i]-'0'))
			continue
		}

		if !words {
			continue
		}

		for index, word := range digitWords {
			if strings.HasPrefix(line[i:], word) {
				digits = append(digits, index+1)
				break
			}
		}
	}

	return digits
}
