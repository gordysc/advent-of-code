// Advent of Code 2015, day 3: Perfectly Spherical Houses in a Vacuum.
// https://adventofcode.com/2015/day/3
package main

import (
	"aoc/lib/grid"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 3, part1, part2)
}

// part1 counts the houses that get at least one present when Santa alone
// follows every move.
func part1(in string) any {
	return deliver(in, 1)
}

// part2 counts the houses when Santa and Robo-Santa take turns: Santa follows
// the first move, Robo-Santa the second, and so on.
func part2(in string) any {
	return deliver(in, 2)
}

// deliver spreads the moves across `santas` deliverers who take turns, and
// returns how many distinct houses were visited. Everyone starts at the origin,
// which counts as a visited house.
func deliver(in string, santas int) int {
	positions := make([]grid.Point, santas)
	visited := set.Of(grid.Point{})

	for i, r := range in {
		who := i % santas
		positions[who] = positions[who].Add(grid.DirFromRune[r])

		visited.Add(positions[who])
	}

	return visited.Len()
}
