// Advent of Code 2018, day 17: Reservoir Research.
// https://adventofcode.com/2018/day/17
package main

import (
	"math"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 17, part1, part2)
}

// Tiles in the scan. Sand is dry, flowing water is wet sand that water only
// passes through, and settled water is water that has come to rest.
const (
	sand    = '.'
	clay    = '#'
	flowing = '|'
	settled = '~'
)

// springX is the column of the spring. The spring is at y=0.
const springX = 500

// scan is the ground slice. Column 0 of the grid is x = offset, and the grid
// rows are the real y values from 0 down to maxY.
type scan struct {
	g          grid.Grid[byte]
	offset     int
	minY, maxY int
}

// part1 counts every tile the water reaches inside the clay's y range.
func part1(in string) any {
	s := run(in)

	return s.count(func(b byte) bool { return b == flowing || b == settled })
}

// part2 counts only the water that stays when the spring stops.
func part2(in string) any {
	s := run(in)

	return s.count(func(b byte) bool { return b == settled })
}

// run reads the clay veins and lets the water flow from the spring until it
// can go nowhere new.
func run(in string) scan {
	var veins [][4]int // x1, x2, y1, y2 of each vein
	minX, maxX := springX, springX
	minY, maxY := math.MaxInt, 0

	for _, line := range input.Lines(in) {
		n := input.Ints(line)

		v := [4]int{n[0], n[0], n[1], n[2]}
		if line[0] == 'y' {
			v = [4]int{n[1], n[2], n[0], n[0]}
		}

		veins = append(veins, v)
		minX, maxX = min(minX, v[0]), max(maxX, v[1])
		minY, maxY = min(minY, v[2]), max(maxY, v[3])
	}

	// Water can spill one column past the outermost clay on each side.
	s := scan{
		g:      grid.New[byte](maxX-minX+3, maxY+1),
		offset: minX - 1,
		minY:   minY,
		maxY:   maxY,
	}

	for i := range s.g.Cells {
		s.g.Cells[i] = sand
	}

	for _, v := range veins {
		for x := v[0]; x <= v[1]; x++ {
			for y := v[2]; y <= v[3]; y++ {
				s.set(x, y, clay)
			}
		}
	}

	s.flow(springX, 0, 0)

	return s
}

// flow moves water into the tile at (x, y) and returns the x where the water
// stopped. dir is 0 when the water falls into the tile, and -1 or +1 when it
// spreads sideways into it.
//
// Falling water goes down as far as it can. When it lands on clay or settled
// water, it spreads both ways. If both sides end at a wall of clay, the row is
// a basin floor, so it fills with settled water. The recursion then unwinds
// to the tile above, which sees settled water below it and spreads in turn.
// That is how a basin fills up row by row. Recursion depth is about the depth
// of the scan, and Go grows goroutine stacks as needed, so this is safe.
func (s *scan) flow(x, y, dir int) int {
	here := s.at(x, y)
	if here == clay {
		return x
	}

	if here == sand {
		s.set(x, y, flowing)
	}

	if y == s.maxY {
		return x
	}

	if s.at(x, y+1) == sand {
		s.flow(x, y+1, 0)
	}

	// Water that still has nothing solid below it falls away here.
	if below := s.at(x, y+1); below != clay && below != settled {
		return x
	}

	if dir != 0 {
		return s.flow(x+dir, y, dir)
	}

	left := s.flow(x-1, y, -1)
	right := s.flow(x+1, y, 1)

	if s.at(left, y) == clay && s.at(right, y) == clay {
		for fx := left + 1; fx < right; fx++ {
			s.set(fx, y, settled)
		}
	}

	return x
}

// at returns the tile at the real coordinates (x, y).
func (s *scan) at(x, y int) byte {
	return s.g.At(grid.P(x-s.offset, y))
}

// set changes the tile at the real coordinates (x, y).
func (s *scan) set(x, y int, b byte) {
	s.g.Set(grid.P(x-s.offset, y), b)
}

// count counts the tiles that match keep between minY and maxY. Rows above
// the highest clay do not count, even though water falls through them.
func (s *scan) count(keep func(byte) bool) int {
	n := 0

	for y := s.minY; y <= s.maxY; y++ {
		for _, b := range s.g.Cells[y*s.g.W : (y+1)*s.g.W] {
			if keep(b) {
				n++
			}
		}
	}

	return n
}
