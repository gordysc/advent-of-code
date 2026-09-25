// Advent of Code 2021, day 8: Seven Segment Search.
// https://adventofcode.com/2021/day/8
package main

import (
	"math/bits"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 8, part1, part2)
}

// entry is one line of notes: the ten unique patterns and the four output
// digits. Each pattern is a bit mask with one bit for each segment a to g.
type entry struct {
	patterns []uint8
	outputs  []uint8
}

// part1 counts the output digits that are 1, 4, 7 or 8.
//
// These four digits each use a number of segments that no other digit uses
// (2, 4, 3 and 7), so the segment count alone identifies them.
func part1(in string) any {
	count := 0

	for _, e := range parse(in) {
		for _, out := range e.outputs {
			switch bits.OnesCount8(out) {
			case 2, 3, 4, 7:
				count++
			}
		}
	}

	return count
}

// part2 decodes each four-digit output value and adds them.
//
// The code first finds 1 and 4 by their segment counts. Then each output
// digit is identified by its segment count and by how many segments it
// shares with 1 and with 4. These three numbers are different for each of
// the ten digits.
func part2(in string) any {
	total := 0

	for _, e := range parse(in) {
		var one, four uint8

		for _, p := range e.patterns {
			switch bits.OnesCount8(p) {
			case 2:
				one = p
			case 4:
				four = p
			}
		}

		value := 0
		for _, out := range e.outputs {
			value = value*10 + decode(out, one, four)
		}

		total += value
	}

	return total
}

// decode returns the digit shown by pattern p, given the patterns for 1 and 4.
func decode(p, one, four uint8) int {
	withOne := bits.OnesCount8(p & one)
	withFour := bits.OnesCount8(p & four)

	switch bits.OnesCount8(p) {
	case 2:
		return 1
	case 3:
		return 7
	case 4:
		return 4
	case 7:
		return 8
	case 5:
		// 2, 3 and 5 all use five segments.
		switch {
		case withOne == 2:
			return 3
		case withFour == 3:
			return 5
		default:
			return 2
		}
	default:
		// 0, 6 and 9 all use six segments.
		switch {
		case withOne == 1:
			return 6
		case withFour == 4:
			return 9
		default:
			return 0
		}
	}
}

// parse reads each line into an entry.
func parse(in string) []entry {
	var entries []entry

	for _, line := range input.Lines(in) {
		left, right, _ := strings.Cut(line, " | ")

		entries = append(entries, entry{
			patterns: masks(left),
			outputs:  masks(right),
		})
	}

	return entries
}

// masks turns space-separated segment words into bit masks.
func masks(s string) []uint8 {
	var out []uint8

	for _, word := range strings.Fields(s) {
		var m uint8
		for _, c := range word {
			m |= 1 << (c - 'a')
		}

		out = append(out, m)
	}

	return out
}
