// Advent of Code 2023, day 22: Sand Slabs.
// https://adventofcode.com/2023/day/22
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 22, part1, part2)
}

// brick is a line of cubes from (x1, y1, z1) to (x2, y2, z2). The parser
// makes sure that each first coordinate is not more than the second.
type brick struct {
	x1, y1, z1 int
	x2, y2, z2 int
}

// stack records, after the bricks fall, which bricks hold up which.
// below[i] lists the bricks directly under brick i that touch it.
// above[i] lists the bricks directly on top of brick i.
type stack struct {
	below [][]int
	above [][]int
}

// part1 counts the bricks we can remove without making any other brick fall.
//
// We can remove a brick when each brick on top of it has another brick
// under it too.
func part1(in string) any {
	s := settle(parse(in))
	count := 0

	for i := range s.above {
		safe := true

		for _, j := range s.above[i] {
			if len(s.below[j]) == 1 {
				safe = false
				break
			}
		}

		if safe {
			count++
		}
	}

	return count
}

// part2 adds, for each brick, the number of other bricks that fall when we
// remove it.
//
// A brick falls when all the bricks under it fall. We start from the removed
// brick and walk up. For each brick above a falling brick, we count how many
// of its supports fall. When that count gets to the number of its supports,
// it falls too.
func part2(in string) any {
	s := settle(parse(in))
	total := 0

	// fallen[j] is the number of falling supports of brick j. We reuse the
	// slice for each removed brick and reset only the entries we touched.
	fallen := make([]int, len(s.below))

	for i := range s.below {
		queue := []int{i}
		var touched []int

		for head := 0; head < len(queue); head++ {
			for _, j := range s.above[queue[head]] {
				if fallen[j] == 0 {
					touched = append(touched, j)
				}

				fallen[j]++

				if fallen[j] == len(s.below[j]) {
					queue = append(queue, j)
				}
			}
		}

		// The queue holds the removed brick and all the bricks that fell.
		total += len(queue) - 1

		for _, j := range touched {
			fallen[j] = 0
		}
	}

	return total
}

// settle lets all the bricks fall to rest and records which bricks touch.
//
// We drop the bricks from the lowest to the highest. A height map holds, for
// each (x, y) column, the top of the stack and the brick at that top. A brick
// comes to rest one above the highest top under it. The bricks at that
// highest top are its supports.
func settle(bricks []brick) stack {
	slices.SortFunc(bricks, func(a, b brick) int {
		return a.z1 - b.z1
	})

	width, depth := 0, 0

	for _, b := range bricks {
		width = max(width, b.x2+1)
		depth = max(depth, b.y2+1)
	}

	// top and owner are flat slices indexed by y*width + x. An owner of -1
	// means the floor.
	top := make([]int, width*depth)
	owner := make([]int, width*depth)

	for i := range owner {
		owner[i] = -1
	}

	s := stack{
		below: make([][]int, len(bricks)),
		above: make([][]int, len(bricks)),
	}

	for i, b := range bricks {
		rest := 0

		for y := b.y1; y <= b.y2; y++ {
			for x := b.x1; x <= b.x2; x++ {
				rest = max(rest, top[y*width+x])
			}
		}

		// One support can touch the brick in several cells, so skip the
		// ones we already have.
		for y := b.y1; y <= b.y2; y++ {
			for x := b.x1; x <= b.x2; x++ {
				k := y*width + x
				j := owner[k]

				if top[k] == rest && j >= 0 && !slices.Contains(s.below[i], j) {
					s.below[i] = append(s.below[i], j)
					s.above[j] = append(s.above[j], i)
				}
			}
		}

		// The brick now rests with its bottom at rest+1.
		height := b.z2 - b.z1 + 1

		for y := b.y1; y <= b.y2; y++ {
			for x := b.x1; x <= b.x2; x++ {
				top[y*width+x] = rest + height
				owner[y*width+x] = i
			}
		}
	}

	return s
}

// parse reads lines like "1,0,1~1,2,1" into bricks.
func parse(in string) []brick {
	var bricks []brick

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		bricks = append(bricks, brick{
			x1: min(n[0], n[3]), y1: min(n[1], n[4]), z1: min(n[2], n[5]),
			x2: max(n[0], n[3]), y2: max(n[1], n[4]), z2: max(n[2], n[5]),
		})
	}

	return bricks
}
