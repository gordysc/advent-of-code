// Advent of Code 2018, day 18: Settlers of The North Pole.
// https://adventofcode.com/2018/day/18
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 18, part1, part2)
}

// Acres of the lumber collection area.
const (
	open       = '.'
	trees      = '|'
	lumberyard = '#'
)

// Minutes to run for each part.
const (
	shortRun = 10
	longRun  = 1_000_000_000
)

// part1 returns the resource value after ten minutes.
func part1(in string) any {
	return resourceValue(in, shortRun)
}

// part2 returns the resource value after a billion minutes. The puzzle text
// gives no answer for the example. The example's forest dies out and every
// acre becomes open ground, so this gives 0 there.
func part2(in string) any {
	return resourceValue(in, longRun)
}

// resourceValue runs the area for the given number of minutes and multiplies
// the wooded acres by the lumberyards. The area soon falls into a loop, so
// each state is stored with the minute it first appeared. When a state comes
// back, the loop length is known. Whole loops change nothing, so only the
// minutes left over after the last whole loop still need to run.
func resourceValue(in string, minutes int) int {
	g := grid.Parse(in)
	seen := map[string]int{}

	for minute := 0; minute < minutes; minute++ {
		// A byte slice cannot be a map key, but a string copy of it can.
		key := string(g.Cells)

		if first, ok := seen[key]; ok {
			loop := minute - first

			for range (minutes - minute) % loop {
				g = step(g)
			}

			break
		}

		seen[key] = minute
		g = step(g)
	}

	return g.Count(isByte(trees)) * g.Count(isByte(lumberyard))
}

// step returns the area one minute later. Every acre changes at the same
// time, so the new state goes into a fresh grid.
func step(g grid.Grid[byte]) grid.Grid[byte] {
	next := grid.New[byte](g.W, g.H)

	for p, acre := range g.All() {
		wooded, yards := 0, 0

		for _, q := range g.Neighbors8(p) {
			switch g.At(q) {
			case trees:
				wooded++
			case lumberyard:
				yards++
			}
		}

		switch {
		case acre == open && wooded >= 3:
			acre = trees
		case acre == trees && yards >= 3:
			acre = lumberyard
		case acre == lumberyard && (yards == 0 || wooded == 0):
			acre = open
		}

		next.Set(p, acre)
	}

	return next
}

// isByte returns a test that matches one kind of acre, for Grid.Count.
func isByte(b byte) func(byte) bool {
	return func(c byte) bool { return c == b }
}
