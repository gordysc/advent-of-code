// Advent of Code 2020, day 24: Lobby Layout.
// https://adventofcode.com/2020/day/24
package main

import (
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 24, part1, part2)
}

// hexDirs gives the step for each of the six hex directions in axial
// coordinates. X moves east, and Y moves south-east. With these two axes
// every hex tile has one unique Point, and the six neighbours are the six
// steps below.
var hexDirs = map[string]grid.Point{
	"e":  grid.P(1, 0),
	"w":  grid.P(-1, 0),
	"ne": grid.P(1, -1),
	"nw": grid.P(0, -1),
	"se": grid.P(0, 1),
	"sw": grid.P(-1, 1),
}

// part1 counts the black tiles after all the flips.
func part1(in string) any {
	return flipAll(in).Len()
}

// part2 runs the daily art exhibit for 100 days and counts the black tiles.
//
// Each day, a black tile with zero or more than two black neighbours turns
// white, and a white tile with exactly two black neighbours turns black.
// Only tiles next to a black tile can be black the next day, so each day
// counts the black neighbours of those tiles only.
func part2(in string) any {
	black := flipAll(in)

	for range 100 {
		counts := map[grid.Point]int{}
		for p := range black {
			for _, d := range hexDirs {
				counts[p.Add(d)]++
			}
		}

		next := set.New[grid.Point]()
		for p, n := range counts {
			if n == 2 || (n == 1 && black.Has(p)) {
				next.Add(p)
			}
		}

		black = next
	}

	return black.Len()
}

// flipAll follows each line to a tile and flips it. A tile that is flipped
// twice is white again, so the set holds only the tiles flipped an odd
// number of times.
func flipAll(in string) set.Set[grid.Point] {
	black := set.New[grid.Point]()

	for _, line := range input.Lines(in) {
		p := walk(line)

		if black.Has(p) {
			black.Remove(p)
		} else {
			black.Add(p)
		}
	}

	return black
}

// walk follows a line of directions from the reference tile and returns the
// tile where it stops. The directions have no separators. "n" and "s" always
// start a two-letter direction, and "e" and "w" alone are one-letter
// directions.
func walk(line string) grid.Point {
	p := grid.P(0, 0)

	for i := 0; i < len(line); i++ {
		step := line[i : i+1]
		if step == "n" || step == "s" {
			step = line[i : i+2]
			i++
		}

		p = p.Add(hexDirs[step])
	}

	return p
}
