// Advent of Code 2022, day 5: Supply Stacks.
// https://adventofcode.com/2022/day/5
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 5, part1, part2)
}

// part1 moves the crates one at a time, so each move reverses the order of
// the crates that it lifts. It gives the top crate of each stack.
func part1(in string) any {
	return rearrange(in, true)
}

// part2 moves all the crates of one step together, so they keep their
// order. It gives the top crate of each stack.
func part2(in string) any {
	return rearrange(in, false)
}

// rearrange does all the moves on the stacks and gives the letters of the
// top crates, from the first stack to the last.
//
// Each stack is a slice with the bottom crate first. A move takes the last
// count crates off the source stack and appends them to the target stack.
// When the crane lifts one crate at a time, the moved crates are in reverse
// order.
func rearrange(in string, oneAtATime bool) string {
	blocks := input.Blocks(in)
	stacks := parseStacks(blocks[0])

	for _, line := range blocks[1] {
		n := input.UInts(line)
		count, from, to := n[0], n[1]-1, n[2]-1

		cut := len(stacks[from]) - count
		moved := stacks[from][cut:]

		// moved shares memory with the source stack. The loop copies each
		// crate into the target stack, so a later append to the source
		// can write over that memory without harm.
		for i := range moved {
			if oneAtATime {
				stacks[to] = append(stacks[to], moved[len(moved)-1-i])
			} else {
				stacks[to] = append(stacks[to], moved[i])
			}
		}

		stacks[from] = stacks[from][:cut]
	}

	var tops strings.Builder

	for _, stack := range stacks {
		if len(stack) > 0 {
			tops.WriteByte(stack[len(stack)-1])
		}
	}

	return tops.String()
}

// parseStacks reads the drawing of the stacks. The last line has the stack
// numbers, and its field count gives the number of stacks. The letter of
// stack i is in column 1+4*i. The drawing lists the top crates first, so
// the rows are read from the bottom up.
//
// A line can be shorter than the full width when an editor removes its
// trailing spaces, so the column check makes sure it is in the line.
func parseStacks(drawing []string) [][]byte {
	numbers := drawing[len(drawing)-1]
	stacks := make([][]byte, len(strings.Fields(numbers)))

	for row := len(drawing) - 2; row >= 0; row-- {
		line := drawing[row]

		for i := range stacks {
			col := 1 + 4*i

			if col < len(line) && line[col] != ' ' {
				stacks[i] = append(stacks[i], line[col])
			}
		}
	}

	return stacks
}
