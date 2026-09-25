// Advent of Code 2020, day 11: Seating System.
// https://adventofcode.com/2020/day/11
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 11, part1, part2)
}

// part1 counts the occupied seats after the layout stops changing. Each seat
// looks at the seats in the eight cells next to it, and a seat empties when
// four or more of them are occupied.
func part1(in string) any {
	return settle(grid.Parse(in), false, 4)
}

// part2 counts the occupied seats after the layout stops changing. Each seat
// now looks along the eight directions to the first seat it can see, and a
// seat empties when five or more of those seats are occupied.
func part2(in string) any {
	return settle(grid.Parse(in), true, 5)
}

// settle runs the seating rules until no seat changes and returns the number
// of occupied seats. The floor never changes, so the function first finds
// the seats and the neighbours of each seat one time. After that, each round
// is a pass over flat slices with no grid lookups.
//
// When sight is true, a neighbour is the first seat in each direction.
// Otherwise it is the seat in the next cell. A seat empties when it has
// tolerance or more occupied neighbours.
func settle(g grid.Grid[byte], sight bool, tolerance int) int {
	index := make(map[grid.Point]int)
	var seats []grid.Point
	for p, c := range g.All() {
		if c == 'L' {
			index[p] = len(seats)
			seats = append(seats, p)
		}
	}

	neighbours := make([][]int, len(seats))
	for i, p := range seats {
		for _, d := range grid.Dirs8 {
			q := p.Add(d)
			for sight && g.InBounds(q) && g.At(q) == '.' {
				q = q.Add(d)
			}

			if j, ok := index[q]; ok {
				neighbours[i] = append(neighbours[i], j)
			}
		}
	}

	occupied := make([]bool, len(seats))
	next := make([]bool, len(seats))
	for {
		changed := false
		for i, ns := range neighbours {
			count := 0
			for _, j := range ns {
				if occupied[j] {
					count++
				}
			}

			switch {
			case !occupied[i] && count == 0:
				next[i] = true
				changed = true
			case occupied[i] && count >= tolerance:
				next[i] = false
				changed = true
			default:
				next[i] = occupied[i]
			}
		}

		occupied, next = next, occupied
		if !changed {
			break
		}
	}

	total := 0
	for _, o := range occupied {
		if o {
			total++
		}
	}

	return total
}
