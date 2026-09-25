// Advent of Code 2018, day 2: Inventory Management System.
// https://adventofcode.com/2018/day/2
package main

import (
	"aoc/lib/input"
	"aoc/lib/strx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 2, part1, part2)
}

// part1 multiplies the number of IDs that have some letter exactly twice by
// the number of IDs that have some letter exactly three times.
func part1(in string) any {
	twos, threes := 0, 0

	for _, id := range input.Lines(in) {
		var counts [26]int
		for i := range len(id) {
			counts[id[i]-'a']++
		}

		hasTwo, hasThree := false, false
		for _, n := range counts {
			hasTwo = hasTwo || n == 2
			hasThree = hasThree || n == 3
		}

		if hasTwo {
			twos++
		}
		if hasThree {
			threes++
		}
	}

	return twos * threes
}

// part2 finds the two IDs that differ by one letter in the same position, and
// returns the letters they have in common. The puzzle uses a different example
// for this part. On the part 1 example it gives "abcde", from abcdef and
// abcdee.
func part2(in string) any {
	ids := input.Lines(in)

	for i, a := range ids {
		for _, b := range ids[i+1:] {
			if strx.Hamming(a, b) != 1 {
				continue
			}

			common := make([]byte, 0, len(a)-1)
			for k := range len(a) {
				if a[k] == b[k] {
					common = append(common, a[k])
				}
			}

			return string(common)
		}
	}

	return nil
}
