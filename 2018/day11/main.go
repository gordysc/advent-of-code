// Advent of Code 2018, day 11: Chronal Charge.
// https://adventofcode.com/2018/day/11
package main

import (
	"fmt"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 11, part1, part2)
}

// size is the width and height of the grid of fuel cells.
const size = 300

// square is the top-left cell of a square, its size and its total power.
type square struct {
	x, y, size, power int
}

// part1 finds the 3x3 square with the most total power.
func part1(in string) any {
	best := bestSquare(table(input.Int(in)), 3)

	return fmt.Sprintf("%d,%d", best.x, best.y)
}

// part2 finds the square of any size with the most total power.
func part2(in string) any {
	sums := table(input.Int(in))
	best := square{power: -1 << 62} // lower than any real total

	for s := 1; s <= size; s++ {
		if sq := bestSquare(sums, s); sq.power > best.power {
			best = sq
		}
	}

	return fmt.Sprintf("%d,%d,%d", best.x, best.y, best.size)
}

// power returns the power level of the fuel cell at x,y.
func power(x, y, serial int) int {
	rack := x + 10
	p := (rack*y + serial) * rack

	return p/100%10 - 5
}

// table builds a summed-area table. sums[y][x] is the total power of every
// cell from 1,1 to x,y. Row 0 and column 0 stay zero, so the lookups in
// bestSquare need no bounds checks.
func table(serial int) [][]int {
	sums := make([][]int, size+1)
	for y := range sums {
		sums[y] = make([]int, size+1)
	}

	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			// The two rectangles above and to the left share the area above-left
			// of the cell, so that area is counted twice and one copy comes off.
			sums[y][x] = power(x, y, serial) + sums[y-1][x] + sums[y][x-1] - sums[y-1][x-1]
		}
	}

	return sums
}

// bestSquare finds the s×s square with the most total power. With the
// summed-area table, the total of any square takes four lookups.
func bestSquare(sums [][]int, s int) square {
	best := square{power: -1 << 62} // lower than any real total

	for y := 1; y+s-1 <= size; y++ {
		for x := 1; x+s-1 <= size; x++ {
			x2, y2 := x+s-1, y+s-1

			// Take the rectangle to x2,y2, remove the strips above and to the
			// left, then add back the corner that both strips removed.
			total := sums[y2][x2] - sums[y-1][x2] - sums[y2][x-1] + sums[y-1][x-1]

			if total > best.power {
				best = square{x: x, y: y, size: s, power: total}
			}
		}
	}

	return best
}
