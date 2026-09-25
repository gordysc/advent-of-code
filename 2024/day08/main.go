// Advent of Code 2024, day 8: Resonant Collinearity.
// https://adventofcode.com/2024/day/8
package main

import (
	"aoc/lib/grid"
	"aoc/lib/mathx"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 8, part1, part2)
}

// part1 counts the map cells that hold an antinode. Each pair of antennas
// with the same frequency makes two antinodes: one on each side of the pair,
// at the same distance as the distance between the two antennas.
func part1(in string) any {
	city, antennas := parse(in)
	antinodes := grid.New[bool](city.W, city.H)

	for _, group := range antennas {
		// Pairs gives each pair one time, so we mark both sides here.
		for a, b := range slicesx.Pairs(group) {
			for _, p := range []grid.Point{a.Add(a.Sub(b)), b.Add(b.Sub(a))} {
				if antinodes.InBounds(p) {
					antinodes.Set(p, true)
				}
			}
		}
	}

	return antinodes.Count(func(on bool) bool { return on })
}

// part2 counts the map cells that hold an antinode, where an antinode is any
// cell exactly in line with two antennas of the same frequency.
//
// The step between two cells on the line is the distance between the antennas
// divided by the GCD of its X and Y parts. This gives the smallest move that
// lands exactly on a cell of the line. On real inputs the GCD is usually 1,
// so the step is the full distance, but dividing is correct in all cases.
func part2(in string) any {
	city, antennas := parse(in)
	antinodes := grid.New[bool](city.W, city.H)

	for _, group := range antennas {
		for a, b := range slicesx.Pairs(group) {
			d := b.Sub(a)
			g := mathx.GCD(d.X, d.Y)
			step := grid.P(d.X/g, d.Y/g)

			// Walk from a in both directions until the walk leaves the map.
			// The walk forward also goes through b.
			for _, dir := range []grid.Point{step, step.Reverse()} {
				for p := a; antinodes.InBounds(p); p = p.Add(dir) {
					antinodes.Set(p, true)
				}
			}
		}
	}

	return antinodes.Count(func(on bool) bool { return on })
}

// parse reads the city map, and groups the antenna positions by frequency.
// A frequency is a letter or a digit; '.' is an empty cell.
func parse(in string) (grid.Grid[byte], map[byte][]grid.Point) {
	city := grid.Parse(in)
	antennas := map[byte][]grid.Point{}

	for p, c := range city.All() {
		if c == '.' {
			continue
		}

		// Appending to a nil slice from a missing map key is allowed. append
		// makes a new slice, and we store it back in the map.
		antennas[c] = append(antennas[c], p)
	}

	return city, antennas
}
