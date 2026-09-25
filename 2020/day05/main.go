// Advent of Code 2020, day 5: Binary Boarding.
// https://adventofcode.com/2020/day/5
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 5, part1, part2)
}

// part1 finds the highest seat ID on a boarding pass.
func part1(in string) any {
	return slices.Max(seatIDs(in))
}

// part2 finds our seat. It is the only missing ID where the IDs just below
// and just above are both on a boarding pass. After a sort, this shows as
// two neighbouring IDs with a difference of 2.
//
// The example has only four unrelated boarding passes and no such gap, so
// part 2 has no answer for it.
func part2(in string) any {
	ids := seatIDs(in)
	slices.Sort(ids)

	for i := 1; i < len(ids); i++ {
		if ids[i]-ids[i-1] == 2 {
			return ids[i] - 1
		}
	}

	return nil
}

// seatIDs decodes each boarding pass into its seat ID. The row is the first
// 7 characters and the column is the last 3, both in binary, and the ID is
// row*8 + column. That is the same as reading all 10 characters as one
// binary number, with 'B' and 'R' as 1 and 'F' and 'L' as 0.
func seatIDs(in string) []int {
	var ids []int

	for _, pass := range input.Lines(in) {
		id := 0
		for _, c := range pass {
			id <<= 1
			if c == 'B' || c == 'R' {
				id |= 1
			}
		}

		ids = append(ids, id)
	}

	return ids
}
