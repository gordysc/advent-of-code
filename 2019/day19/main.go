// Advent of Code 2019, day 19: Tractor Beam.
// https://adventofcode.com/2019/day/19
//
// The puzzle text has no example program, only drawings of a beam. example.txt
// holds a small handmade Intcode program. It reads x and y and outputs 1 when
// y <= x <= 2*y, so its beam is a cone between the lines x = y and x = 2*y.
// It checks x < y and 2*y < x with less-than instructions, adds the two
// results, and outputs 1 when the sum is 0. In the 50x50 area 650 points are
// pulled, and the first 100x100 square fits at x = 297, y = 198, so the
// answers are 650 and 2970198.
package main

import (
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 19, part1, part2)
}

// scanSize is the side of the area that part 1 scans, and shipSize is the side
// of the ship that part 2 fits into the beam. The drawings in the puzzle text
// use a 10x10 area and a 10x10 ship instead.
const (
	scanSize = 50
	shipSize = 100
)

// part1 counts the points in the scan area that the beam pulls.
func part1(in string) any {
	return countPulled(drone(intcode.Parse(in)), scanSize)
}

// part2 finds the square for the ship that is closest to the emitter, and
// returns the x of its top-left corner times 10000 plus its y.
func part2(in string) any {
	x, y := closestSquare(drone(intcode.Parse(in)), shipSize)

	return x*10000 + y
}

// drone returns a pulled function for the program: pulled(x, y) reports
// whether the beam pulls the point (x, y). The program stops after one
// answer, so each call runs a fresh machine. The returned function is a
// closure that keeps program for all later calls.
func drone(program []int) func(x, y int) bool {
	return func(x, y int) bool {
		m := intcode.New(program)
		m.Send(x, y)
		out := m.Run()

		return len(out) > 0 && out[0] == 1
	}
}

// countPulled counts the pulled points in the size x size area at the emitter.
func countPulled(pulled func(x, y int) bool, size int) int {
	n := 0

	for y := range size {
		for x := range size {
			if pulled(x, y) {
				n++
			}
		}
	}

	return n
}

// closestSquare finds the first size x size square that fits in the beam, and
// returns its top-left corner.
//
// The beam is a cone that widens away from the emitter. The walk goes down
// the rows and follows the left edge of the beam. The point there is the
// bottom-left corner of a possible square, so the top-right corner is size-1
// to the right and size-1 up. Both edges of the beam move right as y grows,
// so when these two corners are pulled, the other two corners are pulled too,
// and a cone that holds all four corners holds the whole square. Because the
// left edge never moves left, each row continues the scan from the x of the
// row above. The walk starts at row
// size-1, because a square cannot fit above it. The beam has points in every
// row that far from the emitter, so the scan along a row always stops.
func closestSquare(pulled func(x, y int) bool, size int) (int, int) {
	x := 0

	for y := size - 1; ; y++ {
		for !pulled(x, y) {
			x++
		}

		top := y - (size - 1)
		if pulled(x+size-1, top) {
			return x, top
		}
	}
}
