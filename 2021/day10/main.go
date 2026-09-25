// Advent of Code 2021, day 10: Syntax Scoring.
// https://adventofcode.com/2021/day/10
package main

import (
	"slices"

	"aoc/lib/ds"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 10, part1, part2)
}

// closer maps each opening character to its closing character.
var closer = map[rune]rune{'(': ')', '[': ']', '{': '}', '<': '>'}

// part1 adds the syntax error scores of the corrupted lines.
//
// A line is corrupted when a closing character does not match the last
// open chunk. The score depends on that first bad character.
func part1(in string) any {
	points := map[rune]int{')': 3, ']': 57, '}': 1197, '>': 25137}

	total := 0
	for _, line := range input.Lines(in) {
		bad, _ := check(line)
		total += points[bad]
	}

	return total
}

// part2 finds the middle completion score of the incomplete lines.
//
// For an incomplete line, the stack of open chunks tells which closing
// characters are missing. The code closes them from the top of the stack
// down. Each character multiplies the score by 5 and then adds its value.
func part2(in string) any {
	points := map[rune]int{')': 1, ']': 2, '}': 3, '>': 4}

	var scores []int

	for _, line := range input.Lines(in) {
		bad, open := check(line)
		if bad != 0 {
			continue
		}

		score := 0
		for !open.Empty() {
			score = score*5 + points[closer[open.Pop()]]
		}

		scores = append(scores, score)
	}

	slices.Sort(scores)

	return scores[len(scores)/2]
}

// check reads a line and returns the first bad closing character, or 0 if
// there is none. It also returns the stack of chunks that are still open.
func check(line string) (rune, ds.Stack[rune]) {
	var open ds.Stack[rune]

	for _, c := range line {
		if _, isOpener := closer[c]; isOpener {
			open.Push(c)
			continue
		}

		if open.Empty() || closer[open.Pop()] != c {
			return c, open
		}
	}

	return 0, open
}
