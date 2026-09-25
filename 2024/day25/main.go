// Advent of Code 2024, day 25: Code Chronicle.
// https://adventofcode.com/2024/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 25, part1, nil)
}

// part1 counts the pairs of a lock and a key that fit together, where no
// column of the lock overlaps the same column of the key.
func part1(in string) any {
	var locks, keys [][]int
	space := 0

	for _, block := range input.Blocks(in) {
		if len(block) < 2 {
			continue
		}

		// The top and bottom rows are the solid base of a lock or a key. The
		// rows between them are the space that pins and cuts share.
		space = len(block) - 2

		// A lock has its base at the top. A key has its base at the bottom.
		if block[0][0] == '#' {
			locks = append(locks, heights(block))
		} else {
			keys = append(keys, heights(block))
		}
	}

	count := 0

	for _, lock := range locks {
		for _, key := range keys {
			if fits(lock, key, space) {
				count++
			}
		}
	}

	return count
}

// heights returns the height of each column of a lock or a key. The base row
// is not part of the height, so a column with only its base has height 0.
func heights(block []string) []int {
	h := make([]int, len(block[0]))

	for _, row := range block {
		for x := range len(row) {
			if row[x] == '#' {
				h[x]++
			}
		}
	}

	for x := range h {
		h[x]--
	}

	return h
}

// fits reports whether the lock and the key have room in every column.
func fits(lock, key []int, space int) bool {
	for x := range lock {
		if lock[x]+key[x] > space {
			return false
		}
	}

	return true
}
