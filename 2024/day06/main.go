// Advent of Code 2024, day 6: Guard Gallivant.
// https://adventofcode.com/2024/day/6
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 6, part1, part2)
}

// part1 counts the distinct positions that the guard visits before she
// leaves the map.
func part1(in string) any {
	lab, g, ok := parse(in)
	if !ok {
		return nil
	}

	visited := make([]bool, len(lab.Cells))
	visited[index(lab, g.pos)] = true
	count := 1

	for {
		next, inside := g.step(lab)
		if !inside {
			return count
		}

		g = next

		if i := index(lab, g.pos); !visited[i] {
			visited[i] = true
			count++
		}
	}
}

// part2 counts the positions where one new obstruction traps the guard in a
// loop. The guard's start position is not allowed.
//
// A new obstruction changes the walk only if it is on the guard's route. So we
// walk the route from part 1 and try each cell the first time the guard is
// about to step into it. Up to that step, the walk with the obstruction is the
// same as the walk without it. So each try starts from the guard's current
// state and not from the start.
func part2(in string) any {
	lab, g, ok := parse(in)
	if !ok {
		return nil
	}

	start := g.pos

	// tried marks the cells that already had an obstruction. A cell counts
	// only once, and only the first visit is a correct place to start a try.
	tried := make([]bool, len(lab.Cells))
	tried[index(lab, start)] = true

	// seen[cell*4+dir] holds the number of the last try that saw the guard
	// turn at that cell to face dir. Each try uses a new number, so we never
	// have to clear the slice.
	seen := make([]int, len(lab.Cells)*4)
	try := 0
	count := 0

	for {
		ahead := g.pos.Add(grid.Dirs4[g.dir])
		if c, inside := lab.Get(ahead); inside && c != '#' && !tried[index(lab, ahead)] {
			tried[index(lab, ahead)] = true
			try++

			lab.Set(ahead, '#')
			if loops(lab, g, seen, try) {
				count++
			}

			lab.Set(ahead, '.')
		}

		next, inside := g.step(lab)
		if !inside {
			return count
		}

		g = next
	}
}

// guard is the guard's position and the direction she faces. dir is an index
// into grid.Dirs4, so dir 0 is up and dir+1 is a right turn.
type guard struct {
	pos grid.Point
	dir int
}

// parse reads the lab map. It returns the map with the guard's cell cleared,
// and the guard at her start. The guard always starts facing up ('^').
func parse(in string) (grid.Grid[byte], guard, bool) {
	lab := grid.Parse(in)

	start, ok := lab.FindByte('^')
	if !ok {
		return lab, guard{}, false
	}

	lab.Set(start, '.')

	return lab, guard{pos: start, dir: 0}, true
}

// step moves the guard one time. If an obstruction is ahead she turns right
// and stays on her cell; if not, she moves forward. It returns false when the
// move takes her off the map.
//
// g has a value receiver, so g is a copy. Changes to it do not change the
// caller's guard, and we return the copy as the new state.
func (g guard) step(lab grid.Grid[byte]) (guard, bool) {
	ahead := g.pos.Add(grid.Dirs4[g.dir])

	c, inside := lab.Get(ahead)
	if !inside {
		return g, false
	}

	if c == '#' {
		g.dir = (g.dir + 1) % 4
	} else {
		g.pos = ahead
	}

	return g, true
}

// loops reports whether the guard walks in a loop from state g.
//
// Every loop has at least one turn, because a straight walk leaves the map.
// So it is enough to record the states just after a turn. If the same state
// comes back, the walk repeats forever. There are far fewer turns than
// steps, so this is faster than recording every step.
func loops(lab grid.Grid[byte], g guard, seen []int, try int) bool {
	for {
		next, inside := g.step(lab)
		if !inside {
			return false
		}

		if next.dir != g.dir {
			key := index(lab, next.pos)*4 + next.dir
			if seen[key] == try {
				return true
			}

			seen[key] = try
		}

		g = next
	}
}

// index turns a point into its position in lab.Cells. The grid stores its
// rows one after the other in a single slice.
func index(lab grid.Grid[byte], p grid.Point) int {
	return p.Y*lab.W + p.X
}
