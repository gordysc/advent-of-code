// Advent of Code 2024, day 3: Mull It Over.
// https://adventofcode.com/2024/day/3
package main

import (
	"regexp"
	"strconv"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 3, part1, part2)
}

// part1 adds the results of all the valid mul instructions in the corrupted
// memory.
//
// The puzzle has two examples. The file example.txt holds the part 2 example,
// because it also works for part 1. It is the part 1 example with a do() and
// a don't() added. On this example, part 1 gives 161 and part 2 gives 48.
// The part 1 example gives 161 for both parts, because it has no do() or
// don't().
func part1(in string) any {
	return run(in, false)
}

// part2 adds the results of the mul instructions that are enabled. A don't()
// disables the mul instructions after it, and a do() enables them again.
func part2(in string) any {
	return run(in, true)
}

// instruction matches the three instructions. A mul takes two numbers of one
// to three digits, with no spaces. The parentheses capture the two numbers,
// so a match of do() or don't() has empty captures.
var instruction = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)

// run adds the products of the mul instructions from left to right. When
// conditionals is true, it obeys do() and don't(). Memory starts enabled.
func run(in string, conditionals bool) int {
	total := 0
	enabled := true

	// FindAllStringSubmatch returns each match as a slice: the full text of
	// the match first, then one string per capture group. The -1 means "no
	// limit on the number of matches".
	for _, m := range instruction.FindAllStringSubmatch(in, -1) {
		switch m[0] {
		case "do()":
			enabled = true
			continue
		case "don't()":
			enabled = false
			continue
		}

		if conditionals && !enabled {
			continue
		}

		// The pattern allows only digits, so Atoi cannot fail here.
		a, _ := strconv.Atoi(m[1])
		b, _ := strconv.Atoi(m[2])
		total += a * b
	}

	return total
}
