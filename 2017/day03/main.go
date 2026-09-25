// Advent of Code 2017, day 3: Spiral Memory.
// https://adventofcode.com/2017/day/3
package main

import (
	"iter"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 3, part1, part2)
}

// part1 walks the spiral to the input square and returns its Manhattan
// distance from square 1, which sits at the origin.
func part1(in string) any {
	target := input.Int(in)
	square := 0

	for p := range spiral() {
		square++

		if square == target {
			return p.Manhattan(grid.P(0, 0))
		}
	}

	return nil
}

// part2 fills the spiral again, but each square now holds the sum of its
// eight neighbours that already have a value. It returns the first value that
// is larger than the input.
func part2(in string) any {
	target := input.Int(in)
	values := map[grid.Point]int{grid.P(0, 0): 1}

	for p := range spiral() {
		if p == grid.P(0, 0) {
			continue
		}

		sum := 0

		for _, n := range p.Neighbors8() {
			sum += values[n]
		}

		if sum > target {
			return sum
		}

		values[p] = sum
	}

	return nil
}

// spiral yields the position of every square in order, starting with square 1
// at the origin. The legs of the spiral grow by one square after every second
// turn: 1, 1, 2, 2, 3, 3, and so on.
func spiral() iter.Seq[grid.Point] {
	return func(yield func(grid.Point) bool) {
		p, dir := grid.P(0, 0), grid.Right

		if !yield(p) {
			return
		}

		for leg := 1; ; leg++ {
			for range 2 {
				for range leg {
					p = p.Add(dir)

					if !yield(p) {
						return
					}
				}

				dir = dir.TurnLeft()
			}
		}
	}
}
