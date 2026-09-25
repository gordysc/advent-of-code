// Advent of Code 2025, day 12: Christmas Tree Farm.
// https://adventofcode.com/2025/day/12
//
// Day 12 has no second puzzle: part 2 is unlocked by collecting the other 23
// stars, so only part 1 is written here.
package main

import (
	"encoding/binary"
	"slices"
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 12, part1, nil)
}

// part1 returns the number of regions that have room for all of their
// presents.
//
// Packing shapes is a hard problem in general. But the real inputs are easy:
// each region either has room for every present as a full 3x3 block, or it has
// fewer cells than all its presents together. Two area checks decide those
// regions at once. The example is not like that, so for any region that the
// checks do not decide, we fall back to a backtracking search. The search is
// only fast for small regions like the ones in the example.
func part1(in string) any {
	shapes, regions := parse(in)
	if len(shapes) == 0 || len(regions) == 0 {
		return nil
	}

	fits := 0

	for _, r := range regions {
		if r.fits(shapes) {
			fits++
		}
	}

	return fits
}

// shape is one present, as the cells of each way that it can lie. Each
// orientation lists its cells in row-major order, relative to its first cell.
// So the first cell of each orientation is always (0, 0).
type shape struct {
	orientations [][]grid.Point
	// block is the side of the square that holds the shape.
	block int
}

// region is one space under a tree, with the number of presents of each shape
// that must fit in it.
type region struct {
	w, h   int
	counts []int
}

// parse reads the shapes and the regions. A shape block starts with a line
// such as "4:". The other block holds one region per line, such as
// "12x5: 1 0 1 0 2 2".
func parse(in string) ([]shape, []region) {
	var shapes []shape
	var regions []region

	for _, block := range input.Blocks(in) {
		if strings.Contains(block[0], "x") {
			for _, line := range block {
				n := input.UInts(line)
				if len(n) < 2 {
					continue
				}

				regions = append(regions, region{w: n[0], h: n[1], counts: n[2:]})
			}

			continue
		}

		shapes = append(shapes, parseShape(block[1:]))
	}

	return shapes, regions
}

// parseShape reads the drawing of one shape and finds all its orientations.
// It turns the shape four times, and does the same for its mirror image.
// Symmetric shapes give the same orientation more than once, so it keeps only
// one copy of each.
func parseShape(rows []string) shape {
	var cells []grid.Point

	for y, row := range rows {
		for x, c := range row {
			if c == '#' {
				cells = append(cells, grid.P(x, y))
			}
		}
	}

	s := shape{block: max(len(rows), len(rows[0]))}

	for range 2 {
		for range 4 {
			o := normalize(cells)
			if !slices.ContainsFunc(s.orientations, func(have []grid.Point) bool {
				return slices.Equal(have, o)
			}) {
				s.orientations = append(s.orientations, o)
			}

			// Turn a quarter: (x, y) goes to (-y, x).
			for i, p := range cells {
				cells[i] = grid.P(-p.Y, p.X)
			}
		}

		// Mirror: (x, y) goes to (-x, y).
		for i, p := range cells {
			cells[i] = grid.P(-p.X, p.Y)
		}
	}

	return s
}

// normalize sorts the cells in row-major order and moves them so that the
// first cell is at (0, 0). It returns a new slice.
func normalize(cells []grid.Point) []grid.Point {
	out := slices.Clone(cells)

	slices.SortFunc(out, func(a, b grid.Point) int {
		if a.Y != b.Y {
			return a.Y - b.Y
		}

		return a.X - b.X
	})

	first := out[0]
	for i, p := range out {
		out[i] = p.Sub(first)
	}

	return out
}

// fits reports whether all the presents of the region fit in it.
func (r region) fits(shapes []shape) bool {
	need, presents, block := 0, 0, 0

	for i, n := range r.counts {
		need += n * len(shapes[i].orientations[0])
		presents += n
		block = max(block, shapes[i].block)
	}

	// Too few cells: no packing can work.
	if need > r.w*r.h {
		return false
	}

	// Room for each present as a full block, side by side in a simple grid.
	if (r.w/block)*(r.h/block) >= presents {
		return true
	}

	// The search fills the region row by row. Short rows find dead ends sooner,
	// so turn the region to make its rows the short side. The shapes can turn
	// too, so this does not change the answer.
	if r.w > r.h {
		r.w, r.h = r.h, r.w
	}

	p := packer{
		region: r,
		shapes: shapes,
		used:   grid.New[bool](r.w, r.h),
		left:   slices.Clone(r.counts),
		todo:   presents,
		block:  block,
		failed: map[string]bool{},
	}

	return p.pack(0, r.w*r.h-need)
}

// packer holds the state of the backtracking search for one region.
type packer struct {
	region region
	shapes []shape
	// used marks the cells that a present covers, or that the search decided
	// to leave empty.
	used grid.Grid[bool]
	// left gives the number of presents of each shape that still need a place.
	left []int
	// todo is the sum of left.
	todo int
	// block is the largest side of a shape's square.
	block int
	// failed keeps the keys of the states that cannot be finished. Different
	// orders of placements often lead to the same state, and this memo stops
	// the search from trying such a state again.
	failed map[string]bool
}

// pack tries to place the remaining presents. It goes through the cells in
// row-major order, starting at index from. The first free cell must either
// hold the first cell of some present, or stay empty. So the search never
// tries the same packing twice. slack is the number of cells that can still
// stay empty. When it reaches zero, every free cell must get a present.
func (p *packer) pack(from, slack int) bool {
	if p.todo == 0 {
		return true
	}

	// Skip to the first free cell.
	for from < len(p.used.Cells) && p.used.Cells[from] {
		from++
	}

	if from == len(p.used.Cells) {
		return false
	}

	key := p.key(from, slack)
	if p.failed[key] {
		return false
	}

	if !p.try(from, slack) {
		p.failed[key] = true

		return false
	}

	return true
}

// try places a present on the first free cell, or leaves the cell empty, and
// then continues the search. It returns true when the search finds a
// packing.
func (p *packer) try(from, slack int) bool {
	at := grid.P(from%p.region.w, from/p.region.w)

	for i, s := range p.shapes {
		if p.left[i] == 0 {
			continue
		}

		for _, o := range s.orientations {
			if !p.free(at, o) {
				continue
			}

			p.place(at, o, true)
			p.left[i]--
			p.todo--

			ok := p.pack(from+1, slack)

			// Undo the placement before the next try, or before we return.
			p.place(at, o, false)
			p.left[i]++
			p.todo++

			if ok {
				return true
			}
		}
	}

	if slack == 0 {
		return false
	}

	// Leave this cell empty and go on.
	p.used.Cells[from] = true
	ok := p.pack(from+1, slack-1)
	p.used.Cells[from] = false

	return ok
}

// key returns a string that tells the full state of the search at cell from.
//
// All the cells before from are already used. A present covers cells at most
// block-1 rows below its first cell, and its first cell is before from. So the
// used cells from index from onward all fit in a window of (block-1)*w+block
// cells. The window, the presents that are left, and the slack give the full
// state.
func (p *packer) key(from, slack int) string {
	b := binary.AppendUvarint(nil, uint64(from))
	b = binary.AppendUvarint(b, uint64(slack))

	for _, n := range p.left {
		b = binary.AppendUvarint(b, uint64(n))
	}

	// Pack the window into bytes, eight cells in each byte.
	window := p.used.Cells[from:min(len(p.used.Cells), from+(p.block-1)*p.region.w+p.block)]
	for i := 0; i < len(window); i += 8 {
		var bits byte

		for j := i; j < min(i+8, len(window)); j++ {
			if window[j] {
				bits |= 1 << (j - i)
			}
		}

		b = append(b, bits)
	}

	return string(b)
}

// free reports whether the orientation o fits with its first cell at at,
// inside the region and on free cells only.
func (p *packer) free(at grid.Point, o []grid.Point) bool {
	for _, c := range o {
		q := at.Add(c)

		// Get returns false as its second value when q is outside the grid.
		used, ok := p.used.Get(q)
		if !ok || used {
			return false
		}
	}

	return true
}

// place marks the cells of orientation o, with its first cell at at, as used
// or free.
func (p *packer) place(at grid.Point, o []grid.Point, used bool) {
	for _, c := range o {
		p.used.Set(at.Add(c), used)
	}
}
