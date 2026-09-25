// Advent of Code 2017, day 10: Knot Hash.
// https://adventofcode.com/2017/day/10
package main

import (
	"encoding/hex"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 10, part1, part2)
}

// listSize is the number of marks on the circular list. The worked example in
// the puzzle text uses a list of 5 instead, so part 1 on the example does not
// give the 12 from the puzzle text.
const listSize = 256

// rounds is how many times part 2 applies the full length sequence.
const rounds = 64

// suffix is the fixed set of lengths that part 2 adds after the input bytes.
var suffix = []int{17, 31, 73, 47, 23}

// part1 does one round with the input as a list of numbers, then multiplies
// the first two marks.
func part1(in string) any {
	list := knot(input.Ints(in), 1)

	return list[0] * list[1]
}

// part2 computes the full knot hash of the input text.
func part2(in string) any {
	return hash(strings.TrimSpace(in))
}

// hash reads each byte of the text as a length, adds the fixed suffix, and
// does 64 rounds. It then XORs each block of 16 marks into one byte to get
// the dense hash, and writes it as lowercase hex.
func hash(text string) string {
	var lengths []int
	for _, b := range []byte(text) {
		lengths = append(lengths, int(b))
	}

	lengths = append(lengths, suffix...)
	list := knot(lengths, rounds)

	// XOR with 0 changes nothing, so each byte of dense starts at 0 and
	// collects its 16 marks.
	dense := make([]byte, len(list)/16)
	for i, mark := range list {
		dense[i/16] ^= byte(mark)
	}

	// hex.EncodeToString writes two lowercase hex digits for each byte.
	return hex.EncodeToString(dense)
}

// knot ties the knots for the given number of rounds (times). For each
// length it reverses that many marks from the current position, wrapping
// around, then moves forward by the length plus the skip size. The position
// and the skip size carry over from one round to the next.
func knot(lengths []int, times int) []int {
	list := make([]int, listSize)
	for i := range list {
		list[i] = i
	}

	pos, skip := 0, 0

	for range times {
		for _, length := range lengths {
			// i and j walk toward each other and can go past the end of the
			// list, so the modulo wraps them back to the start.
			for i, j := pos, pos+length-1; i < j; i, j = i+1, j-1 {
				a, b := i%listSize, j%listSize
				list[a], list[b] = list[b], list[a]
			}

			pos = (pos + length + skip) % listSize
			skip++
		}
	}

	return list
}
