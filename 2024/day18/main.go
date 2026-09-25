// Advent of Code 2024, day 18: RAM Run.
// https://adventofcode.com/2024/day/18
package main

import (
	"fmt"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 18, part1, part2)
}

// part1 returns the fewest steps from the top-left corner to the bottom-right
// corner after the first bytes fall. The real input uses a 71x71 grid and the
// first 1024 bytes. The example uses a 7x7 grid and the first 12 bytes, and
// gives 22.
func part1(in string) any {
	bytes := parse(in)
	size, count := dimensions(bytes)

	steps, ok := shortest(bytes[:min(count, len(bytes))], size)
	if !ok {
		return nil
	}

	return steps
}

// part2 returns the "x,y" position of the first byte that blocks all paths to
// the exit. The example gives "6,1".
//
// More bytes can only remove paths, never add them. So "the first n bytes
// block the exit" is false up to some n and true from there on, and a binary
// search finds that n with a few path searches.
func part2(in string) any {
	bytes := parse(in)
	size, _ := dimensions(bytes)

	n := search.BinarySearch(1, len(bytes), func(n int) bool {
		_, ok := shortest(bytes[:n], size)
		return !ok
	})

	// BinarySearch returns hi+1 when no prefix blocks the exit.
	if n > len(bytes) {
		return nil
	}

	b := bytes[n-1]

	return fmt.Sprintf("%d,%d", b.X, b.Y)
}

// parse reads one "x,y" byte position per line.
func parse(in string) []grid.Point {
	var bytes []grid.Point

	for _, line := range input.Lines(in) {
		nums := input.Ints(line)
		if len(nums) < 2 {
			continue
		}

		bytes = append(bytes, grid.P(nums[0], nums[1]))
	}

	return bytes
}

// dimensions returns the grid size and the number of bytes for part 1. The
// example is the only input with all coordinates below 7.
func dimensions(bytes []grid.Point) (size, count int) {
	for _, b := range bytes {
		if b.X >= 7 || b.Y >= 7 {
			return 71, 1024
		}
	}

	return 7, 12
}

// shortest returns the fewest steps from (0,0) to the far corner when the
// given bytes are corrupted. It returns false when no path exists.
func shortest(bytes []grid.Point, size int) (int, bool) {
	corrupt := grid.New[bool](size, size)
	for _, b := range bytes {
		corrupt.Set(b, true)
	}

	exit := grid.P(size-1, size-1)

	next := func(p grid.Point) []grid.Point {
		var out []grid.Point

		// Neighbors4 leaves out the points outside the grid.
		for _, q := range corrupt.Neighbors4(p) {
			if !corrupt.At(q) {
				out = append(out, q)
			}
		}

		return out
	}

	return search.BFS(grid.P(0, 0), next, func(p grid.Point) bool { return p == exit })
}
