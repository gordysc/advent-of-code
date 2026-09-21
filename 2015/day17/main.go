// Advent of Code 2015, day 17: No Such Thing as Too Much.
// https://adventofcode.com/2015/day/17
package main

import (
	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// eggnogLiters is the amount of eggnog to store in the puzzle. The worked
// example in the puzzle text stores 25 liters instead, so the example answers
// from this program differ from the ones in the text.
const eggnogLiters = 150

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 17, part1, part2)
}

// part1 counts every combination of containers that holds exactly 150 liters.
func part1(in string) any {
	return slicesx.Sum(waysByCount(input.IntLines(in), eggnogLiters))
}

// part2 finds the smallest number of containers that can hold exactly 150
// liters, and counts the combinations that use that number of containers.
func part2(in string) any {
	// The index is the number of containers, so the first entry that is not
	// zero belongs to the smallest number of containers that works.
	for _, ways := range waysByCount(input.IntLines(in), eggnogLiters) {
		if ways > 0 {
			return ways
		}
	}

	return 0
}

// waysByCount returns a slice in which the entry at index k is the number of
// combinations of exactly k containers that hold exactly target liters. Two
// containers of the same size count as different containers.
func waysByCount(sizes []int, target int) []int {
	// ways[k][v] is the number of ways to fill exactly v liters with exactly k
	// of the containers seen so far. This is dynamic programming: each
	// container updates the table once, which replaces a search through all
	// 2^n subsets.
	ways := make([][]int, len(sizes)+1)

	for k := range ways {
		ways[k] = make([]int, target+1)
	}

	// There is one way to store 0 liters in 0 containers: use nothing.
	ways[0][0] = 1

	for _, size := range sizes {
		// A new combination of k containers is an old combination of k-1
		// containers plus this container. The loop goes down from the highest
		// k, so ways[k-1] still holds the numbers from before this container.
		// A loop that goes up would use the same container two times.
		for k := len(sizes); k >= 1; k-- {
			for v := size; v <= target; v++ {
				ways[k][v] += ways[k-1][v-size]
			}
		}
	}

	// Keep only the column for the full target, one entry for each count.
	counts := make([]int, len(ways))

	for k := range ways {
		counts[k] = ways[k][target]
	}

	return counts
}
