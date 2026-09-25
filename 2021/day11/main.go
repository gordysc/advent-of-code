// Advent of Code 2021, day 11: Dumbo Octopus.
// https://adventofcode.com/2021/day/11
package main

import (
	"aoc/lib/ds"
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 11, part1, part2)
}

// part1 counts all the flashes in the first 100 steps.
func part1(in string) any {
	octopuses := parse(in)
	total := 0

	for range 100 {
		total += step(octopuses)
	}

	return total
}

// part2 finds the first step in which every octopus flashes at the same time.
func part2(in string) any {
	octopuses := parse(in)
	all := len(octopuses.Cells)

	for n := 1; ; n++ {
		if step(octopuses) == all {
			return n
		}
	}
}

// parse reads the energy level of each octopus into a grid.
func parse(in string) grid.Grid[int] {
	return grid.Map(in, func(r rune) int { return int(r - '0') })
}

// step runs one step in place and returns the number of flashes.
//
// Each energy level goes up by one. An octopus flashes when its level goes
// above 9, and the flash adds one to each of its eight neighbours. That can
// make more octopuses flash. An octopus is put on the stack only at the moment
// its level becomes 10, so it flashes one time only. At the end, each octopus
// that flashed goes back to 0.
func step(g grid.Grid[int]) int {
	var flashing ds.Stack[grid.Point]

	charge := func(p grid.Point) {
		level := g.At(p) + 1
		g.Set(p, level)

		if level == 10 {
			flashing.Push(p)
		}
	}

	for p := range g.Points() {
		charge(p)
	}

	for !flashing.Empty() {
		for _, n := range g.Neighbors8(flashing.Pop()) {
			charge(n)
		}
	}

	flashes := 0

	for i, level := range g.Cells {
		if level > 9 {
			g.Cells[i] = 0
			flashes++
		}
	}

	return flashes
}
