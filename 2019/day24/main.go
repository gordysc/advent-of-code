// Advent of Code 2019, day 24: Planet of Discord.
// https://adventofcode.com/2019/day/24
package main

import (
	"math/bits"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 24, part1, part2)
}

// size is the width and height of one grid, tiles is its tile count, and mid
// is the row and column of its middle tile. In part 2 the middle tile holds
// the next grid down, so it is never a bug.
const (
	size  = 5
	tiles = size * size
	mid   = size / 2
)

// minutes is how long part 2 runs. The worked example in the puzzle text runs
// for 10 minutes instead and has 99 bugs then. After 200 minutes the example
// has 1922 bugs.
const minutes = 200

// dirs are the four steps to an adjacent tile: up, down, left and right.
var dirs = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// links holds, for one tile, bit masks of the tiles it touches. same is on its
// own level, outer is on the level that holds this grid, and inner is on the
// level that this grid holds in its middle tile.
type links struct {
	same, outer, inner uint32
}

// part1 finds the first layout that appears twice and gives its biodiversity
// rating.
//
// A layout is a 25-bit mask where bit i is set when tile i (counted row by
// row) has a bug. Tile i is worth 2^i points, and 2^i is the value of bit i,
// so the mask itself is the biodiversity rating.
func part1(in string) any {
	adj := makeLinks(false)
	layout := parse(in)
	seen := map[uint32]bool{}

	for !seen[layout] {
		seen[layout] = true
		layout = step(0, layout, 0, &adj)
	}

	return int(layout)
}

// part2 counts the bugs after the given minutes on the recursive grids.
//
// levels holds one mask per level, with the start grid in the middle. A lower
// index is further out: levels[d] sits in the middle tile of levels[d-1], and
// levels[d+1] sits in the middle tile of levels[d]. In one minute bugs can spread at most one level out or in,
// so minutes+1 levels on each side is enough. The outermost and innermost
// levels stay empty, which lets the loop read levels[d-1] and levels[d+1]
// with no bounds checks.
func part2(in string) any {
	adj := makeLinks(true)
	depth := minutes + 1
	levels := make([]uint32, 2*depth+1)
	next := make([]uint32, len(levels))
	levels[depth] = parse(in)

	for range minutes {
		for d := 1; d < len(levels)-1; d++ {
			next[d] = step(levels[d-1], levels[d], levels[d+1], &adj)
		}

		levels, next = next, levels
	}

	count := 0
	for _, level := range levels {
		count += bits.OnesCount32(level)
	}

	return count
}

// parse reads the grid into a mask with bit r*size+c set for a bug in row r,
// column c.
func parse(in string) uint32 {
	var mask uint32

	for r, line := range input.Lines(in) {
		for c, ch := range line {
			if ch == '#' {
				mask |= 1 << (r*size + c)
			}
		}
	}

	return mask
}

// makeLinks works out which tiles each tile touches.
//
// With recursive false this is the plain grid of part 1: up to four tiles on
// the same level. With recursive true a step off the edge lands on the outer
// level, on the tile next to its middle in that direction. A step onto the
// middle tile lands on the whole near edge of the inner level: five tiles. The
// middle tile itself gets no links, so it stays empty.
func makeLinks(recursive bool) [tiles]links {
	var adj [tiles]links

	for i := range tiles {
		r, c := i/size, i%size
		if recursive && r == mid && c == mid {
			continue
		}

		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]

			switch {
			case nr < 0 || nr >= size || nc < 0 || nc >= size:
				if recursive {
					adj[i].outer |= bit(mid+d[0], mid+d[1])
				}

			case recursive && nr == mid && nc == mid:
				adj[i].inner |= innerEdge(d)

			default:
				adj[i].same |= bit(nr, nc)
			}
		}
	}

	return adj
}

// innerEdge returns the mask of the inner level's edge that a step in
// direction d meets. A step down meets its top row, a step right meets its
// left column, and so on.
func innerEdge(d [2]int) uint32 {
	var mask uint32

	for k := range size {
		switch d {
		case [2]int{1, 0}:
			mask |= bit(0, k)
		case [2]int{-1, 0}:
			mask |= bit(size-1, k)
		case [2]int{0, 1}:
			mask |= bit(k, 0)
		case [2]int{0, -1}:
			mask |= bit(k, size-1)
		}
	}

	return mask
}

// bit returns the mask with only the bit for row r, column c set.
func bit(r, c int) uint32 {
	return 1 << (r*size + c)
}

// step runs one minute on a level, given the level outside it and the level
// inside it. For each tile, ANDing a level with the tile's link mask keeps
// only the bugs next to it, and OnesCount32 counts them.
func step(outer, cur, inner uint32, adj *[tiles]links) uint32 {
	var next uint32

	for i, l := range adj {
		n := bits.OnesCount32(cur&l.same) +
			bits.OnesCount32(outer&l.outer) +
			bits.OnesCount32(inner&l.inner)

		bug := cur&(1<<i) != 0
		if n == 1 || (!bug && n == 2) {
			next |= 1 << i
		}
	}

	return next
}
