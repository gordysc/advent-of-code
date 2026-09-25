// Advent of Code 2017, day 9: Stream Processing.
// https://adventofcode.com/2017/day/9
package main

import (
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 9, part1, part2)
}

// part1 adds up the score of every group. A group scores one more than the
// group around it.
func part1(in string) any {
	score, _ := scan(in)

	return score
}

// part2 counts the characters inside garbage, without the "<" and ">" around
// it and without cancelled characters or the "!" that cancels them.
func part2(in string) any {
	_, garbage := scan(in)

	return garbage
}

// scan reads the stream one character at a time. It returns the total group
// score and the number of garbage characters.
func scan(in string) (int, int) {
	score, depth, garbage := 0, 0, 0
	inGarbage := false

	for i := 0; i < len(in); i++ {
		c := in[i]

		if inGarbage {
			switch c {
			case '!':
				i++ // skip the cancelled character
			case '>':
				inGarbage = false
			default:
				garbage++
			}

			continue
		}

		switch c {
		case '{':
			depth++
			score += depth
		case '}':
			depth--
		case '<':
			inGarbage = true
		}
	}

	return score, garbage
}
