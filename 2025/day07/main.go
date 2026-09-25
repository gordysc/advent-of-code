// Advent of Code 2025, day 7: Laboratories.
// https://adventofcode.com/2025/day/7
package main

import (
	"aoc/lib/grid"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 7, part1, part2)
}

// part1 counts how many times a beam is split.
func part1(in string) any {
	splits, _, ok := simulate(in)
	if !ok {
		return nil
	}

	return splits
}

// part2 counts the timelines of one tachyon particle. At each splitter the
// particle takes both paths, so each split makes one timeline into two.
func part2(in string) any {
	_, timelines, ok := simulate(in)
	if !ok {
		return nil
	}

	return timelines
}

// simulate moves the beams down the manifold one row at a time. It returns the
// number of splits and the number of timelines at the bottom. It returns false
// when the manifold has no start S.
//
// Beams that reach the same column merge into one beam. So it is not necessary
// to follow each timeline: for each column, we keep only the number of
// timelines that have a beam there. A splitter sends that number to the
// columns on its left and right. This makes the time linear in the grid size,
// although the number of timelines grows very fast.
func simulate(in string) (splits, timelines int, ok bool) {
	g := grid.Parse(in)

	start, found := g.FindByte('S')
	if !found {
		return 0, 0, false
	}

	// counts[x] is the number of timelines with a beam in column x.
	counts := make([]int, g.W)
	counts[start.X] = 1

	for y := start.Y + 1; y < g.H; y++ {
		next := make([]int, g.W)

		for x, c := range counts {
			if c == 0 {
				continue
			}

			if g.At(grid.P(x, y)) != '^' {
				next[x] += c
				continue
			}

			// Beams from a merged column split only one time, so this adds 1
			// and not c.
			splits++

			// A beam that leaves the side of the manifold is lost.
			if x > 0 {
				next[x-1] += c
			}

			if x+1 < g.W {
				next[x+1] += c
			}
		}

		counts = next
	}

	return splits, slicesx.Sum(counts), true
}
