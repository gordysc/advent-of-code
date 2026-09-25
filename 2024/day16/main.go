// Advent of Code 2024, day 16: Reindeer Maze.
// https://adventofcode.com/2024/day/16
package main

import (
	"aoc/lib/ds"
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 16, part1, part2)
}

// part1 returns the lowest score a reindeer can get from S to E. A step
// forward costs 1 point and a 90 degree turn costs 1000 points.
//
// The puzzle has two examples. The file example.txt holds the first one, where
// part 1 gives 7036 and part 2 gives 45. On the second example, part 1 gives
// 11048 and part 2 gives 64.
func part1(in string) any {
	m := parseMaze(in)
	from := m.search([]int{m.state(m.start, east)}, false)

	best, ok := m.bestAtEnd(from)
	if !ok {
		return nil
	}

	return best
}

// part2 counts the tiles that are on at least one path with the lowest score.
//
// We search once forward from the start and once backward from the end. A
// state is on a best path when its cost from the start plus its cost to the
// end is equal to the best score.
func part2(in string) any {
	m := parseMaze(in)
	from := m.search([]int{m.state(m.start, east)}, false)

	best, ok := m.bestAtEnd(from)
	if !ok {
		return nil
	}

	// The reindeer can arrive at E facing any direction, so the backward
	// search starts from all four.
	var ends []int
	for d := range grid.Dirs4 {
		ends = append(ends, m.state(m.end, d))
	}

	to := m.search(ends, true)

	tiles := 0

	for i := range m.Cells {
		for d := range grid.Dirs4 {
			s := i*4 + d
			if from[s] >= 0 && to[s] >= 0 && from[s]+to[s] == best {
				tiles++
				break
			}
		}
	}

	return tiles
}

// east is the index of Right in grid.Dirs4. The reindeer starts facing east.
const east = 1

// maze is the map with its start and end tiles. It embeds grid.Grid, so the
// grid fields and methods, such as m.W and m.At, are available on maze
// directly.
type maze struct {
	grid.Grid[byte]
	start, end grid.Point
}

// parseMaze reads the map and finds the S and E tiles.
func parseMaze(in string) maze {
	g := grid.Parse(in)
	start, _ := g.FindByte('S')
	end, _ := g.FindByte('E')

	return maze{Grid: g, start: start, end: end}
}

// state packs a tile and a facing direction into one index. The index of
// direction d is an index into grid.Dirs4.
func (m maze) state(p grid.Point, d int) int {
	return (p.Y*m.W+p.X)*4 + d
}

// bestAtEnd returns the lowest cost to reach E in any direction. It returns
// false when E cannot be reached.
func (m maze) bestAtEnd(dist []int) (int, bool) {
	best := -1

	for d := range grid.Dirs4 {
		c := dist[m.state(m.end, d)]
		if c >= 0 && (best < 0 || c < best) {
			best = c
		}
	}

	return best, best >= 0
}

// search runs Dijkstra's algorithm from all the start states and returns the
// lowest cost to each state, or -1 when a state cannot be reached.
//
// When backward is true, the search follows each move in reverse. A step then
// goes from (p, d) back to (p-d, d). A turn is its own reverse, so it does not
// change. The costs are then the costs to get from each state to the starts.
func (m maze) search(starts []int, backward bool) []int {
	// A flat slice indexed by state is much faster than a map here.
	dist := make([]int, len(m.Cells)*4)
	for i := range dist {
		dist[i] = -1
	}

	pq := ds.NewPriorityQueue[int]()

	for _, s := range starts {
		dist[s] = 0
		pq.Push(s, 0)
	}

	for pq.Len() > 0 {
		s, cost := pq.Pop()

		// Skip old queue entries for a state that got a lower cost later.
		if cost > dist[s] {
			continue
		}

		tile, d := s/4, s%4
		p := grid.P(tile%m.W, tile/m.W)

		step := grid.Dirs4[d]
		if backward {
			step = step.Reverse()
		}

		// (d+1)%4 turns right and (d+3)%4 turns left, because Dirs4 is in
		// clockwise order.
		moves := [][2]int{{d, 1}, {(d + 1) % 4, 1000}, {(d + 3) % 4, 1000}}

		for i, mv := range moves {
			next := m.state(p, mv[0])

			// Only the first move is a step. The other two turn in place.
			// The maze has walls all around, so q is always in the grid.
			if i == 0 {
				q := p.Add(step)
				if m.At(q) == '#' {
					continue
				}

				next = m.state(q, d)
			}

			c := cost + mv[1]
			if dist[next] >= 0 && dist[next] <= c {
				continue
			}

			dist[next] = c
			pq.Push(next, c)
		}
	}

	return dist
}
