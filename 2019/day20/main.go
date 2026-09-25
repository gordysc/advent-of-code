// Advent of Code 2019, day 20: Donut Maze.
// https://adventofcode.com/2019/day/20
package main

import (
	"aoc/lib/ds"
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 20, part1, part2)
}

// end is one end of a portal: the open tile next to a two-letter label.
// outer is true for ends on the outside edge of the donut. partner is the
// index of the other end with the same label, or -1 for AA and ZZ.
type end struct {
	label   string
	p       grid.Point
	outer   bool
	partner int
}

// maze is the donut reduced to its portal ends. dist[i][j] is the number of
// steps from end i to end j without using a portal, or 0 when j cannot be
// reached from i.
type maze struct {
	ends        []end
	dist        [][]int
	start, goal int
}

// state is one moment in the search: the portal end we stand on, and the
// level of the maze we are on. Level 0 is the outermost maze.
type state struct {
	at, level int
}

// part1 finds the fewest steps from AA to ZZ when every portal connects two
// places in the same maze.
func part1(in string) any {
	return solve(parse(in), false)
}

// part2 finds the fewest steps from AA to ZZ in the recursive maze. There, an
// inner portal leads one level down and an outer portal one level up. AA and
// ZZ exist only on level 0, and the outer portals on level 0 are walls. It
// returns nil when there is no way out.
func part2(in string) any {
	return solve(parse(in), true)
}

// parse reads the maze, finds the portal ends, and measures the distances
// between them.
func parse(in string) maze {
	g := readGrid(in)

	// The donut's outer edge is the box around every wall tile. Labels sit
	// outside it, so an end on this box belongs to an outer portal.
	lo, hi := grid.P(g.W, g.H), grid.P(0, 0)

	for p, c := range g.All() {
		if c == '#' {
			lo = grid.P(min(lo.X, p.X), min(lo.Y, p.Y))
			hi = grid.P(max(hi.X, p.X), max(hi.Y, p.Y))
		}
	}

	var m maze
	byLabel := map[string][]int{}
	endAt := map[grid.Point]int{}

	for p, c := range g.All() {
		if !isLetter(c) {
			continue
		}

		label, open, ok := portal(g, p)
		if !ok {
			continue
		}

		i := len(m.ends)
		outer := open.X == lo.X || open.X == hi.X || open.Y == lo.Y || open.Y == hi.Y
		m.ends = append(m.ends, end{label: label, p: open, outer: outer, partner: -1})
		byLabel[label] = append(byLabel[label], i)
		endAt[open] = i
	}

	for _, pair := range byLabel {
		if len(pair) == 2 {
			m.ends[pair[0]].partner = pair[1]
			m.ends[pair[1]].partner = pair[0]
		}
	}

	m.start = byLabel["AA"][0]
	m.goal = byLabel["ZZ"][0]

	m.dist = make([][]int, len(m.ends))

	for i, e := range m.ends {
		m.dist[i] = distances(g, e.p, endAt, len(m.ends))
	}

	return m
}

// readGrid builds a grid from the maze text. The rows have different lengths
// when trailing spaces are missing, so grid.Parse would reject them. Every row
// is padded with spaces to the length of the longest row instead.
func readGrid(in string) grid.Grid[byte] {
	lines := input.Lines(in)

	w := 0

	for _, line := range lines {
		w = max(w, len(line))
	}

	g := grid.New[byte](w, len(lines))

	for i := range g.Cells {
		g.Cells[i] = ' '
	}

	for y, line := range lines {
		copy(g.Cells[y*w:], line)
	}

	return g
}

// portal checks whether the letter at p is the first letter of a label, which
// is the left letter of a label read across or the top letter of a label read
// down. It returns the label and the open tile next to it.
func portal(g grid.Grid[byte], p grid.Point) (string, grid.Point, bool) {
	for _, d := range []grid.Point{grid.Right, grid.Down} {
		q := p.Add(d)
		if c, ok := g.Get(q); !ok || !isLetter(c) {
			continue
		}

		label := string([]byte{g.At(p), g.At(q)})

		// The open tile is either just before the first letter or just
		// after the second one.
		for _, open := range []grid.Point{p.Sub(d), q.Add(d)} {
			if c, ok := g.Get(open); ok && c == '.' {
				return label, open, true
			}
		}
	}

	return "", grid.Point{}, false
}

// distances runs a breadth-first search over the open tiles from p. It returns
// the steps to each of the n portal ends, with 0 for the ends it cannot reach.
func distances(g grid.Grid[byte], p grid.Point, endAt map[grid.Point]int, n int) []int {
	out := make([]int, n)
	steps := map[grid.Point]int{p: 0}
	queue := ds.NewQueue(p)

	for !queue.Empty() {
		cur := queue.Pop()

		for _, next := range g.Neighbors4(cur) {
			if _, seen := steps[next]; seen || g.At(next) != '.' {
				continue
			}
			steps[next] = steps[cur] + 1

			if i, ok := endAt[next]; ok {
				out[i] = steps[next]
			}

			queue.Push(next)
		}
	}

	return out
}

// solve finds the fewest steps from AA to ZZ, or nil when ZZ cannot be
// reached. When recursive is false, the level always stays 0.
//
// Each move walks from the current end to another end, and then steps through
// that portal for 1 more step. Dijkstra's algorithm fits because the walks
// have different lengths.
//
// A recursive maze has no bottom, so the search could go down forever when
// there is no way out. The levels are capped at the number of portal ends: a
// shortest path in these puzzles does not go that deep, and the cap makes the
// search stop.
func solve(m maze, recursive bool) any {
	maxLevel := len(m.ends)

	// next is a closure, so it can read m and recursive from solve.
	next := func(s state) []search.Edge[state] {
		var edges []search.Edge[state]

		for j, d := range m.dist[s.at] {
			if d == 0 {
				continue
			}

			// ZZ is the exit only on level 0. On other levels it is a wall.
			if j == m.goal {
				if s.level == 0 {
					edges = append(edges, search.Edge[state]{To: state{j, 0}, Cost: d})
				}
				continue
			}

			// AA is the entrance: there is no reason to walk back to it.
			if j == m.start {
				continue
			}

			level := s.level
			if recursive {
				if m.ends[j].outer {
					level--
				} else {
					level++
				}
			}

			if level < 0 || level > maxLevel {
				continue
			}

			to := state{m.ends[j].partner, level}
			edges = append(edges, search.Edge[state]{To: to, Cost: d + 1})
		}

		return edges
	}

	goal := func(s state) bool { return s.at == m.goal }

	steps, ok := search.Dijkstra(state{m.start, 0}, next, goal)
	if !ok {
		return nil
	}

	return steps
}

// isLetter reports whether c is part of a portal label.
func isLetter(c byte) bool {
	return c >= 'A' && c <= 'Z'
}
