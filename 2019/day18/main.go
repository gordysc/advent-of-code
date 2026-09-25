// Advent of Code 2019, day 18: Many-Worlds Interpretation.
// https://adventofcode.com/2019/day/18
//
// The example is the 136-step vault from part 1. Its single @ sits in a
// corridor, not in the middle of an open 3x3 area, so part 2 cannot split it
// into four robots and has no answer for the example.
package main

import (
	"aoc/lib/ds"
	"aoc/lib/grid"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 18, part1, part2)
}

// numKeys is the number of possible keys, a to z. maxRobots is the number of
// robots in part 2.
const (
	numKeys   = 26
	maxRobots = 4
)

// route is the shortest walk from a start or a key to one key. doors and keys
// are bitmasks: bit i of doors is set when the walk goes through door 'A'+i,
// and bit i of keys is set when the walk steps on key 'a'+i, the last key
// included. A route with 0 steps does not exist.
type route struct {
	steps int
	doors uint32
	keys  uint32
}

// routeTable holds the routes out of every point of interest. Indexes 0 to 25
// are the keys a to z, and index numKeys+r is the start of robot r.
type routeTable [numKeys + maxRobots][numKeys]route

// state is one moment in the search: the point of interest where each robot
// stands, and the keys collected so far. It holds only an array and a number,
// so it is comparable and the search can use it as a map key. The indexes fit
// in a uint8, and a small key makes the search's map faster.
type state struct {
	at   [maxRobots]uint8
	keys uint32
}

// part1 finds the fewest steps for one robot to collect every key.
func part1(in string) any {
	g := grid.Parse(in)
	start, _ := g.FindByte('@')

	return collect(g, []grid.Point{start})
}

// part2 walls off the middle of the vault, puts four robots there, and finds
// the fewest steps for them to collect every key together. It returns nil when
// the vault has no open 3x3 area around a single @ to split.
func part2(in string) any {
	g := grid.Parse(in)

	starts, ok := split(g)
	if !ok {
		return nil
	}

	return collect(g, starts)
}

// split changes g for part 2 and returns the four robot starts. The 3x3 area
// around the @ becomes
//
//	@#@
//	###
//	@#@
//
// A map that already has four @ is used as it is.
func split(g grid.Grid[byte]) ([]grid.Point, bool) {
	var starts []grid.Point

	for p, c := range g.All() {
		if c == '@' {
			starts = append(starts, p)
		}
	}

	if len(starts) == maxRobots {
		return starts, true
	}

	if len(starts) != 1 {
		return nil, false
	}

	mid := starts[0]

	for _, n := range mid.Neighbors8() {
		if c, ok := g.Get(n); !ok || c != '.' {
			return nil, false
		}
	}

	g.Set(mid, '#')
	for _, d := range grid.Dirs4 {
		g.Set(mid.Add(d), '#')
	}

	return []grid.Point{
		mid.Add(grid.UpLeft), mid.Add(grid.UpRight),
		mid.Add(grid.DownLeft), mid.Add(grid.DownRight),
	}, true
}

// collect finds the fewest steps for the robots at starts to collect every
// key, or nil when some key cannot be reached.
//
// Walking cell by cell would make the search huge. Instead, the routes between
// the starts and the keys are found once. The search then only moves a robot
// straight to a key it does not have yet, through doors it has keys for.
// Dijkstra's algorithm fits because the moves have different step counts.
func collect(g grid.Grid[byte], starts []grid.Point) any {
	var routes routeTable
	var all uint32 // bit i is set when key 'a'+i is in the vault

	for p, c := range g.All() {
		if isKey(c) {
			routes[c-'a'] = routesFrom(g, p)
			all |= 1 << (c - 'a')
		}
	}

	var start state
	for r := range maxRobots {
		// Robots that do not exist stand on a start with no routes, so they
		// never move.
		start.at[r] = uint8(numKeys + r)
	}

	for r, p := range starts {
		routes[numKeys+r] = routesFrom(g, p)
	}

	// next is a closure: it can read routes from the enclosing function, so
	// the search only has to pass it a state.
	next := func(s state) []search.Edge[state] {
		var edges []search.Edge[state]

		for r, from := range s.at {
			for k, rt := range routes[from] {
				// Skip keys already held, routes that do not exist, and
				// routes with a door whose key is not held. x &^ y clears
				// the bits of y in x, so it is 0 when every door is opened.
				if s.keys&(1<<k) != 0 || rt.steps == 0 || rt.doors&^s.keys != 0 {
					continue
				}

				// Arrays are values in Go, so to is a full copy of s and
				// changing to.at does not change s. Keys stepped on along
				// the way are picked up too.
				to := s
				to.at[r] = uint8(k)
				to.keys |= rt.keys
				edges = append(edges, search.Edge[state]{To: to, Cost: rt.steps})
			}
		}

		return edges
	}

	steps, ok := search.Dijkstra(start, next, func(s state) bool { return s.keys == all })
	if !ok {
		return nil
	}

	return steps
}

// routesFrom runs a breadth-first search from p and returns the route to each
// key it can reach. Every cell in the queue carries the route that reached it,
// so the doors and keys along the way are known when the walk finds a key.
func routesFrom(g grid.Grid[byte], p grid.Point) [numKeys]route {
	type step struct {
		p  grid.Point
		rt route
	}

	var out [numKeys]route
	seen := grid.New[bool](g.W, g.H)
	seen.Set(p, true)
	queue := ds.NewQueue(step{p: p})

	for !queue.Empty() {
		cur := queue.Pop()

		for _, n := range g.Neighbors4(cur.p) {
			c := g.At(n)
			if c == '#' || seen.At(n) {
				continue
			}
			seen.Set(n, true)

			rt := cur.rt
			rt.steps++

			switch {
			case isDoor(c):
				rt.doors |= 1 << (c - 'A')
			case isKey(c):
				rt.keys |= 1 << (c - 'a')
				out[c-'a'] = rt
			}

			queue.Push(step{n, rt})
		}
	}

	return out
}

// isKey reports whether c is a key, a to z.
func isKey(c byte) bool {
	return c >= 'a' && c <= 'z'
}

// isDoor reports whether c is a door, A to Z.
func isDoor(c byte) bool {
	return c >= 'A' && c <= 'Z'
}
