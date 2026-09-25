// Advent of Code 2024, day 13: Claw Contraption.
// https://adventofcode.com/2024/day/13
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 13, part1, part2)
}

// part1 adds the fewest tokens that win every prize that can be won, when
// each button can be pressed at most 100 times.
func part1(in string) any {
	total := 0
	for _, m := range parse(in) {
		a, b, ok := m.solve()
		if ok && a <= 100 && b <= 100 {
			total += 3*a + b
		}
	}

	return total
}

// part2 adds the fewest tokens that win every prize that can be won, after
// the prize moves 10000000000000 units further on both axes.
//
// The puzzle text gives no example answer for part 2. On the example, this
// part gives 875318608908.
func part2(in string) any {
	const shift = 10000000000000

	total := 0
	for _, m := range parse(in) {
		m.px += shift
		m.py += shift

		if a, b, ok := m.solve(); ok {
			total += 3*a + b
		}
	}

	return total
}

// machine holds how far each button moves the claw, and where the prize is.
type machine struct {
	ax, ay int
	bx, by int
	px, py int
}

// parse reads one machine from each block of three lines.
func parse(in string) []machine {
	var out []machine

	for _, block := range input.Blocks(in) {
		n := input.Ints(strings.Join(block, "\n"))
		if len(n) != 6 {
			continue
		}

		out = append(out, machine{n[0], n[1], n[2], n[3], n[4], n[5]})
	}

	return out
}

// solve finds the number of A presses and B presses that put the claw on the
// prize. It returns false when no whole, non-negative solution exists.
//
// The presses a and b must fit two equations:
//
//	a*ax + b*bx = px
//	a*ay + b*by = py
//
// Cramer's rule solves this 2x2 system directly. When the determinant is not
// zero, there is exactly one solution, so "the fewest tokens" has only one
// candidate. A determinant of zero means the two buttons move in the same
// direction. The real inputs do not have such machines, so we skip them.
func (m machine) solve() (int, int, bool) {
	det := m.ax*m.by - m.ay*m.bx
	if det == 0 {
		return 0, 0, false
	}

	aNum := m.px*m.by - m.py*m.bx
	bNum := m.ax*m.py - m.ay*m.px

	// Go's % keeps the sign of the left operand, but a remainder of zero is
	// zero for any sign. So this test works for negative values too.
	if aNum%det != 0 || bNum%det != 0 {
		return 0, 0, false
	}

	a, b := aNum/det, bNum/det
	if a < 0 || b < 0 {
		return 0, 0, false
	}

	return a, b, true
}
