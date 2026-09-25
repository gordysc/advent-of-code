// Advent of Code 2020, day 17: Conway Cubes.
// https://adventofcode.com/2020/day/17
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 17, part1, part2)
}

// cycles is the number of boot cycles that both parts run.
const cycles = 6

// part1 runs the boot process in three dimensions and counts the active cubes.
func part1(in string) any {
	return simulate(parse(in), 3)
}

// part2 runs the same process in four dimensions.
func part2(in string) any {
	return simulate(parse(in), 4)
}

// cube is a position in up to four dimensions. The unused dimensions stay 0.
type cube [4]int

// simulate runs the boot cycles on a sparse set of active cubes and returns
// how many are active at the end.
//
// Each cycle, every active cube adds 1 to the count of each of its
// neighbours. Only cubes with a count can become or stay active, so we never
// visit the empty space. A cube is active next cycle if its count is 3, or if
// its count is 2 and it is active now.
func simulate(active map[cube]bool, dims int) int {
	offsets := neighbourOffsets(dims)

	for range cycles {
		counts := map[cube]int{}
		for c := range active {
			for _, d := range offsets {
				counts[cube{c[0] + d[0], c[1] + d[1], c[2] + d[2], c[3] + d[3]}]++
			}
		}

		next := map[cube]bool{}
		for c, n := range counts {
			if n == 3 || (n == 2 && active[c]) {
				next[c] = true
			}
		}

		active = next
	}

	return len(active)
}

// neighbourOffsets returns every offset in {-1, 0, 1} for the first dims
// dimensions, except the zero offset. That gives 26 offsets in 3D and 80 in
// 4D.
func neighbourOffsets(dims int) []cube {
	offsets := []cube{{}}
	for d := range dims {
		var grown []cube
		for _, o := range offsets {
			for _, delta := range []int{-1, 0, 1} {
				// cube is an array, so o is a copy. Setting o[d] and
				// appending o stores a new value each time.
				o[d] = delta
				grown = append(grown, o)
			}
		}

		offsets = grown
	}

	// The loop above builds every combination, including the cube itself.
	// Drop the zero offset so that a cube does not count as its own neighbour.
	var result []cube
	for _, o := range offsets {
		if o != (cube{}) {
			result = append(result, o)
		}
	}

	return result
}

// parse reads the 2D starting slice. Each '#' is an active cube at z = 0 and
// w = 0.
func parse(in string) map[cube]bool {
	active := map[cube]bool{}
	for y, line := range input.Lines(in) {
		for x, ch := range line {
			if ch == '#' {
				active[cube{x, y, 0, 0}] = true
			}
		}
	}

	return active
}
