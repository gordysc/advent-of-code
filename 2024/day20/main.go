// Advent of Code 2024, day 20: Race Condition.
// https://adventofcode.com/2024/day/20
package main

import (
	"aoc/lib/grid"
	"aoc/lib/mathx"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 20, part1, part2)
}

// part1 counts the cheats of up to 2 picoseconds that save at least 100
// picoseconds.
//
// The example saves less than 100 picoseconds with every cheat, so on the
// 15x15 example we count the cheats that save at least 64 picoseconds. The
// puzzle text lists exactly 1 such cheat.
func part1(in string) any {
	return countCheats(in, 2, 64)
}

// part2 counts the cheats of up to 20 picoseconds that save at least 100
// picoseconds.
//
// On the 15x15 example we count the cheats that save at least 50
// picoseconds. The puzzle text lists 285 such cheats in total.
func part2(in string) any {
	return countCheats(in, 20, 50)
}

// countCheats counts the cheats that go through walls for at most maxCheat
// picoseconds and save at least the threshold. The threshold is 100 for real
// inputs and exampleMin for the 15x15 example.
//
// A cheat goes from track tile p to track tile q in Manhattan(p, q)
// picoseconds. The track is one path without branches, so the time from the
// start to each tile tells us everything. The cheat saves
// dist[q] - dist[p] - Manhattan(p, q) picoseconds.
func countCheats(in string, maxCheat, exampleMin int) any {
	g := grid.Parse(in)

	start, ok := g.FindByte('S')
	if !ok {
		return nil
	}

	threshold := 100
	if g.W == 15 && g.H == 15 {
		threshold = exampleMin
	}

	// Flood returns the number of steps from start to each tile it reaches.
	dist := search.Flood(start, func(p grid.Point) []grid.Point {
		var out []grid.Point

		for _, q := range g.Neighbors4(p) {
			if g.At(q) != '#' {
				out = append(out, q)
			}
		}

		return out
	})

	// A grid of times is much faster to look up than the map. -1 marks walls.
	times := grid.New[int](g.W, g.H)
	for i := range times.Cells {
		times.Cells[i] = -1
	}

	for p, d := range dist {
		times.Set(p, d)
	}

	count := 0

	// Try every end tile q within maxCheat steps of each track tile p. Each
	// cheat is counted once, from the tile where it starts.
	for p, dp := range dist {
		for dy := -maxCheat; dy <= maxCheat; dy++ {
			span := maxCheat - mathx.Abs(dy)

			for dx := -span; dx <= span; dx++ {
				// Get returns false when q is outside the grid.
				dq, ok := times.Get(p.Add(grid.P(dx, dy)))
				if !ok || dq < 0 {
					continue
				}

				saved := dq - dp - mathx.Abs(dx) - mathx.Abs(dy)
				if saved >= threshold {
					count++
				}
			}
		}
	}

	return count
}
