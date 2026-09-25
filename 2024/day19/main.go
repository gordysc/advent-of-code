// Advent of Code 2024, day 19: Linen Layout.
// https://adventofcode.com/2024/day/19
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 19, part1, part2)
}

// part1 counts the designs that the towel patterns can make.
func part1(in string) any {
	patterns, designs := parse(in)
	possible := 0

	for _, d := range designs {
		if ways(d, patterns) > 0 {
			possible++
		}
	}

	return possible
}

// part2 adds up the number of different ways to make each design.
func part2(in string) any {
	patterns, designs := parse(in)
	total := 0

	for _, d := range designs {
		total += ways(d, patterns)
	}

	return total
}

// parse reads the comma-separated towel patterns from the first block and the
// designs, one per line, from the second block.
func parse(in string) ([]string, []string) {
	blocks := input.Blocks(in)
	if len(blocks) < 2 || len(blocks[0]) == 0 {
		return nil, nil
	}

	patterns := strings.Split(blocks[0][0], ", ")

	return patterns, blocks[1]
}

// ways counts the different ways to make design from the towel patterns.
//
// count[i] is the number of ways to make design[i:]. We fill it from the end
// back to the start. For each pattern that design[i:] starts with, add the
// ways to make the rest after that pattern. The empty rest has one way.
func ways(design string, patterns []string) int {
	count := make([]int, len(design)+1)
	count[len(design)] = 1

	for i := len(design) - 1; i >= 0; i-- {
		for _, p := range patterns {
			if strings.HasPrefix(design[i:], p) {
				count[i] += count[i+len(p)]
			}
		}
	}

	return count[0]
}
