// Advent of Code 2022, day 2: Rock Paper Scissors.
// https://adventofcode.com/2022/day/2
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 2, part1, part2)
}

// part1 reads the second column as the shape that we play: X is rock, Y is
// paper and Z is scissors. It gives the total score for all rounds.
func part1(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		them, us := shapes(line)
		total += score(them, us)
	}

	return total
}

// part2 reads the second column as the result that we must get: X is a
// loss, Y is a draw and Z is a win. It gives the total score for all rounds.
//
// Number the shapes 0 (rock), 1 (paper) and 2 (scissors). Each shape beats
// the shape one step before it, in a circle. So a draw plays the same shape,
// a win plays the next shape and a loss plays the previous shape. That is an
// offset of -1, 0 or +1 from their shape, modulo 3.
func part2(in string) any {
	total := 0

	for _, line := range input.Lines(in) {
		them, result := shapes(line)
		us := mathx.Mod(them+result-1, 3)
		total += score(them, us)
	}

	return total
}

// shapes turns a line such as "A Y" into two numbers from 0 to 2: one for
// the letter A, B or C and one for the letter X, Y or Z.
func shapes(line string) (int, int) {
	return int(line[0] - 'A'), int(line[2] - 'X')
}

// score gives our score for one round. The shape gives 1, 2 or 3 points. The
// result gives 0 for a loss, 3 for a draw and 6 for a win.
//
// The difference us-them, modulo 3, is 0 for a draw, 1 for a win and 2 for a
// loss. Adding 1 and taking modulo 3 again maps these to 1, 2 and 0, which
// are the result points divided by 3.
func score(them, us int) int {
	outcome := mathx.Mod(us-them+1, 3)

	return us + 1 + outcome*3
}
