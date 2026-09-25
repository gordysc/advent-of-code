// Advent of Code 2022, day 17: Pyroclastic Flow.
// https://adventofcode.com/2022/day/17
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 17, part1, part2)
}

// part1 finds the height of the tower after 2022 rocks fall.
func part1(in string) any {
	return towerHeight(strings.TrimSpace(in), 2022)
}

// part2 finds the height of the tower after one trillion rocks fall.
//
// That is too many rocks to drop one by one. But the rock shapes and the jets
// both repeat. After some time the chamber comes back to the same state: the
// same next rock, the same next jet, and the same shape at the top of the
// tower. From then on, the tower grows by the same amount in each cycle.
func part2(in string) any {
	return towerHeight(strings.TrimSpace(in), 1_000_000_000_000)
}

// width is the width of the chamber.
const width = 7

// rocks are the five rock shapes, in the order they fall. Each shape is a list
// of rows from the bottom up. Bit x of a row is set when the rock fills
// column x, counted from the left edge of the rock.
var rocks = [][]uint8{
	{0b1111},              // -
	{0b010, 0b111, 0b010}, // +
	{0b111, 0b100, 0b100}, // backwards L; bit 2 is the right column
	{0b1, 0b1, 0b1, 0b1},  // |
	{0b11, 0b11},          // square
}

// profileDepth is the number of rows below the top that the surface profile
// looks at. A column with no rock in these rows gets this depth.
const profileDepth = 64

// state is what decides how the tower grows from now on.
type state struct {
	rock, jet int
	profile   [width]int
}

// seen records when a state occurred.
type seen struct {
	rocks, height int
}

// chamber is the tower of settled rock. Each row is a bit mask of the filled
// columns, and row 0 is at the bottom.
type chamber struct {
	rows []uint8
}

// towerHeight drops count rocks and returns the height of the tower.
//
// After each rock it records the state of the chamber. When a state occurs
// again, the rocks between the two times make one cycle. We skip as many
// whole cycles as fit, then drop the last few rocks one by one.
func towerHeight(jets string, count int) int {
	var c chamber
	history := map[state]seen{}
	jet := 0
	skipped := 0

	for n := 0; n < count; n++ {
		jet = c.drop(rocks[n%len(rocks)], jets, jet)

		// Look for a cycle only until we find one.
		if skipped > 0 {
			continue
		}

		key := state{rock: (n + 1) % len(rocks), jet: jet, profile: c.profile()}
		now := seen{rocks: n + 1, height: len(c.rows)}

		prev, ok := history[key]
		if !ok {
			history[key] = now
			continue
		}

		cycleRocks := now.rocks - prev.rocks
		cycles := (count - now.rocks) / cycleRocks
		skipped = cycles * (now.height - prev.height)
		n += cycles * cycleRocks
	}

	return len(c.rows) + skipped
}

// drop lets one rock fall until it comes to rest, and adds it to the tower.
// Jets push the rock before each step down. It returns the index of the next
// jet.
func (c *chamber) drop(rock []uint8, jets string, jet int) int {
	// The rock starts two units from the left wall and three rows above the
	// top of the tower.
	x, y := 2, len(c.rows)+3

	for {
		push := 1
		if jets[jet] == '<' {
			push = -1
		}
		jet = (jet + 1) % len(jets)

		if c.fits(rock, x+push, y) {
			x += push
		}

		if !c.fits(rock, x, y-1) {
			break
		}

		y--
	}

	for i, row := range rock {
		for len(c.rows) <= y+i {
			c.rows = append(c.rows, 0)
		}

		c.rows[y+i] |= row << x
	}

	return jet
}

// fits reports whether the rock can be at column x and row y without hitting
// a wall, the floor, or settled rock.
func (c *chamber) fits(rock []uint8, x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}

	for i, row := range rock {
		mask := row << x

		// A bit past column 6 means the rock goes through the right wall.
		if mask >= 1<<width {
			return false
		}

		if y+i < len(c.rows) && c.rows[y+i]&mask != 0 {
			return false
		}
	}

	return true
}

// profile returns, for each column, how many rows below the top the first
// rock in that column is. It gives the shape of the top of the tower.
func (c *chamber) profile() [width]int {
	var p [width]int

	for x := range width {
		p[x] = profileDepth

		for d := 0; d < profileDepth && d < len(c.rows); d++ {
			if c.rows[len(c.rows)-1-d]&(1<<x) != 0 {
				p[x] = d
				break
			}
		}
	}

	return p
}
