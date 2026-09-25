// Advent of Code 2023, day 23: A Long Walk.
// https://adventofcode.com/2023/day/23
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 23, part1, part2)
}

// part1 finds the longest hike from the start to the end when slopes are
// one-way: a step off a slope tile must go downhill.
func part1(in string) any {
	return longestHike(grid.Parse(in), true)
}

// part2 finds the longest hike when the slopes are ordinary path tiles.
func part2(in string) any {
	return longestHike(grid.Parse(in), false)
}

// edge is a corridor from one junction to another.
type edge struct {
	to    int
	steps int
}

// trail is the maze compressed to a graph. The nodes are the start, the end,
// and every junction (a path tile with three or more path neighbours). Each
// node has a list of the corridors that leave it.
type trail struct {
	edges [][]edge
	start int
	end   int
}

// longestHike compresses the maze to a graph and then tries every route.
//
// Longest path has no fast general method, so we must try all routes. But the
// maze has only a few dozen junctions, and the corridors between them have no
// choices. On the compressed graph the search is small enough to do in full.
func longestHike(g grid.Grid[byte], slopes bool) any {
	t, ok := compress(g, slopes)
	if !ok {
		return nil
	}

	// The search keeps its visited nodes as bits in one uint64.
	if len(t.edges) > 64 {
		return nil
	}

	best := t.longest(t.start, 1<<t.start, 0)
	if best < 0 {
		return nil
	}

	return best
}

// compress finds the start, the end and the junctions, and walks each
// corridor to measure it. It returns false when the maze has no start or end.
func compress(g grid.Grid[byte], slopes bool) (trail, bool) {
	start, okStart := openInRow(g, 0)
	end, okEnd := openInRow(g, g.H-1)
	if !okStart || !okEnd {
		return trail{}, false
	}

	// Give each node a number. The numbers are bit positions in the visited mask.
	nodes := []grid.Point{start, end}
	for p, c := range g.All() {
		if c != '#' && len(openNeighbors(g, p)) >= 3 {
			nodes = append(nodes, p)
		}
	}

	ids := make(map[grid.Point]int, len(nodes))
	for id, p := range nodes {
		ids[p] = id
	}

	edges := make([][]edge, len(nodes))
	for id, p := range nodes {
		for _, next := range openNeighbors(g, p) {
			if to, steps, ok := walk(g, p, next, ids, slopes); ok {
				edges[id] = append(edges[id], edge{to, steps})
			}
		}
	}

	t := trail{edges: edges, start: 0, end: 1}
	t.forceExit()

	return t, true
}

// forceExit prunes the search at the node before the end.
//
// When only one node has a corridor to the end, a hike that reaches that node
// must go to the end at once. If it went anywhere else, it could not come back
// through the node, so it could never finish. Thus we remove the other
// corridors from that node.
//
// The receiver is a copy of trail, but the copy shares the same edges slice
// memory. So the change to t.edges[last] is seen by the caller too.
func (t trail) forceExit() {
	last := -1

	for from, list := range t.edges {
		for _, e := range list {
			if e.to != t.end {
				continue
			}

			if last >= 0 && last != from {
				return
			}
			last = from
		}
	}

	if last < 0 {
		return
	}

	for _, e := range t.edges[last] {
		if e.to == t.end {
			t.edges[last] = []edge{e}

			return
		}
	}
}

// longest returns the most steps from node to the end without a return to a
// node in seen. dist is the number of steps already done. It returns -1 when
// no route gets to the end.
func (t trail) longest(node int, seen uint64, dist int) int {
	if node == t.end {
		return dist
	}

	best := -1

	for _, e := range t.edges[node] {
		bit := uint64(1) << e.to
		if seen&bit != 0 {
			continue
		}

		best = max(best, t.longest(e.to, seen|bit, dist+e.steps))
	}

	return best
}

// walk follows the corridor that starts with the step from junction to first.
// It returns the node at the other end and the length of the corridor. It
// returns false when the corridor is a dead end or a slope blocks it.
func walk(g grid.Grid[byte], junction, first grid.Point, ids map[grid.Point]int, slopes bool) (int, int, bool) {
	if !canStep(g, junction, first, slopes) {
		return 0, 0, false
	}

	prev, cur := junction, first
	steps := 1

	for {
		if id, ok := ids[cur]; ok {
			return id, steps, true
		}

		// A corridor tile has at most two path neighbours, and one of them is
		// the tile we came from. So there is at most one way on.
		moved := false

		for _, next := range openNeighbors(g, cur) {
			if next == prev || !canStep(g, cur, next, slopes) {
				continue
			}

			prev, cur = cur, next
			steps++
			moved = true

			break
		}

		if !moved {
			return 0, 0, false
		}
	}
}

// canStep tells if a hiker on from can step to to. With slopes on, a hiker on
// a slope tile can only step in the direction of the slope.
func canStep(g grid.Grid[byte], from, to grid.Point, slopes bool) bool {
	c := g.At(from)
	if !slopes || c == '.' {
		return true
	}

	return grid.DirFromRune[rune(c)] == to.Sub(from)
}

// openNeighbors returns the neighbours of p that are not forest.
func openNeighbors(g grid.Grid[byte], p grid.Point) []grid.Point {
	var open []grid.Point

	for _, n := range g.Neighbors4(p) {
		if g.At(n) != '#' {
			open = append(open, n)
		}
	}

	return open
}

// openInRow finds the first path tile in row y.
func openInRow(g grid.Grid[byte], y int) (grid.Point, bool) {
	for x := range g.W {
		if p := grid.P(x, y); g.At(p) == '.' {
			return p, true
		}
	}

	return grid.Point{}, false
}
