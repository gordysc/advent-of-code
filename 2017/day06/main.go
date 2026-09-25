// Advent of Code 2017, day 6: Memory Reallocation.
// https://adventofcode.com/2017/day/6
package main

import (
	"fmt"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 6, part1, part2)
}

// part1 counts the redistribution cycles before a layout of blocks repeats.
func part1(in string) any {
	steps, _ := findLoop(in)

	return steps
}

// part2 counts the cycles in the loop, from the first time the repeated
// layout was seen to the time it came back.
func part2(in string) any {
	steps, first := findLoop(in)

	return steps - first
}

// findLoop redistributes the banks until a layout repeats. It returns the
// number of cycles done and the cycle at which that layout first appeared.
func findLoop(in string) (int, int) {
	banks := input.Ints(in)
	seen := map[string]int{}

	for steps := 0; ; steps++ {
		// A slice cannot be a map key, so the layout is stored as its text,
		// such as "[0 2 7 0]".
		key := fmt.Sprint(banks)
		if first, ok := seen[key]; ok {
			return steps, first
		}

		seen[key] = steps
		redistribute(banks)
	}
}

// redistribute empties the fullest bank (the first one on a tie) and deals
// its blocks one at a time to the banks after it, wrapping around.
func redistribute(banks []int) {
	top := 0
	for i, b := range banks {
		if b > banks[top] {
			top = i
		}
	}

	blocks := banks[top]
	banks[top] = 0

	for i := top + 1; blocks > 0; i++ {
		banks[i%len(banks)]++
		blocks--
	}
}
