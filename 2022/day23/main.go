// Advent of Code 2022, day 23: Unstable Diffusion.
// https://adventofcode.com/2022/day/23
//
// example.txt holds the larger example from the puzzle text. It gives 110
// for part 1 and 20 for part 2.
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 23, part1, part2)
}

// part1 runs 10 rounds, then counts the empty tiles in the smallest
// rectangle that holds all the elves.
func part1(in string) any {
	s := newSwarm(grid.Parse(in))

	for r := 0; r < 10; r++ {
		s.round(r)
	}

	lo, hi := s.bounds()
	area := (hi.X - lo.X + 1) * (hi.Y - lo.Y + 1)

	return area - len(s.elves)
}

// part2 finds the first round in which no elf moves.
func part2(in string) any {
	s := newSwarm(grid.Parse(in))

	for r := 0; ; r++ {
		if !s.round(r) {
			return r + 1
		}
	}
}

// margin is the number of free tiles that newSwarm and grow put around the
// elves.
const margin = 32

// check lists the three tiles an elf looks at for each move direction, in
// the order the elves consider them in the first round. The first tile of
// each group is the move direction.
var check = [4][3]grid.Point{
	{grid.Up, grid.UpLeft, grid.UpRight},
	{grid.Down, grid.DownLeft, grid.DownRight},
	{grid.Left, grid.UpLeft, grid.DownLeft},
	{grid.Right, grid.UpRight, grid.DownRight},
}

// swarm holds the elves on a fixed-size grid.
//
// A grid is much faster than a map for the "is this tile taken" tests. The
// elves spread out over time, so grow makes a larger grid when an elf comes
// near the edge.
type swarm struct {
	elves []grid.Point
	taken grid.Grid[bool]

	// votes counts how many elves propose each tile in the current round.
	votes grid.Grid[uint8]
}

// newSwarm finds the elves ('#') in the scan and puts them on a grid with
// free space around them.
func newSwarm(scan grid.Grid[byte]) *swarm {
	s := &swarm{}

	for p, c := range scan.All() {
		if c == '#' {
			s.elves = append(s.elves, p)
		}
	}

	s.grow()

	return s
}

// grow makes a new grid with margin free tiles on each side of the elves,
// and moves the elves onto it.
func (s *swarm) grow() {
	lo, hi := s.bounds()
	shift := grid.P(margin, margin).Sub(lo)

	s.taken = grid.New[bool](hi.X-lo.X+1+2*margin, hi.Y-lo.Y+1+2*margin)
	s.votes = grid.New[uint8](s.taken.W, s.taken.H)

	for i := range s.elves {
		s.elves[i] = s.elves[i].Add(shift)
		s.taken.Set(s.elves[i], true)
	}
}

// bounds returns the top-left and bottom-right corners of the smallest
// rectangle that holds all the elves.
func (s *swarm) bounds() (grid.Point, grid.Point) {
	lo, hi := s.elves[0], s.elves[0]

	for _, e := range s.elves {
		lo = grid.P(min(lo.X, e.X), min(lo.Y, e.Y))
		hi = grid.P(max(hi.X, e.X), max(hi.Y, e.Y))
	}

	return lo, hi
}

// round runs round number r (from 0) and reports if any elf moved.
//
// First, each elf proposes a tile. Then each elf moves to its tile only if
// no other elf proposed the same tile.
func (s *swarm) round(r int) bool {
	proposals := make([]grid.Point, len(s.elves))

	for i, e := range s.elves {
		proposals[i] = s.propose(e, r)

		if proposals[i] != e {
			s.votes.Set(proposals[i], s.votes.At(proposals[i])+1)
		}
	}

	moved := false
	nearEdge := false

	for i, e := range s.elves {
		to := proposals[i]
		if to == e {
			continue
		}

		if s.votes.At(to) == 1 {
			s.taken.Set(e, false)
			s.taken.Set(to, true)
			s.elves[i] = to
			moved = true

			if to.X < 1 || to.Y < 1 || to.X >= s.taken.W-1 || to.Y >= s.taken.H-1 {
				nearEdge = true
			}
		}
	}

	// Clear only the tiles that got votes. This is much faster than a
	// clear of the full grid.
	for _, to := range proposals {
		s.votes.Set(to, 0)
	}

	// The proposals use the old grid positions, so grow only after the
	// votes are clear.
	if nearEdge {
		s.grow()
	}

	return moved
}

// propose returns the tile the elf at e wants to move to in round r. It
// returns e when the elf stays: when it has no neighbours, or when all four
// directions are blocked.
func (s *swarm) propose(e grid.Point, r int) grid.Point {
	alone := true

	for _, d := range grid.Dirs8 {
		if s.taken.At(e.Add(d)) {
			alone = false
			break
		}
	}

	if alone {
		return e
	}

	// Each round, the first direction the elves consider moves to the end
	// of the list.
	for k := range 4 {
		group := check[(r+k)%4]

		if !s.taken.At(e.Add(group[0])) && !s.taken.At(e.Add(group[1])) && !s.taken.At(e.Add(group[2])) {
			return e.Add(group[0])
		}
	}

	return e
}
