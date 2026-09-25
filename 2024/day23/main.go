// Advent of Code 2024, day 23: LAN Party.
// https://adventofcode.com/2024/day/23
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 23, part1, part2)
}

// part1 counts the sets of three computers that all connect to each other,
// where at least one name starts with "t".
func part1(in string) any {
	g := parse(in)
	count := 0

	// Visit each triangle once, with its names in order a < b < c.
	for a, na := range g.adj {
		for b := range na.All() {
			if b <= a {
				continue
			}

			for c := range g.adj[b].All() {
				if c <= b || !na.Has(c) {
					continue
				}

				if g.names[a][0] == 't' || g.names[b][0] == 't' || g.names[c][0] == 't' {
					count++
				}
			}
		}
	}

	return count
}

// part2 finds the largest set of computers that all connect to each other.
// It returns the password: the names in order, joined by commas.
func part2(in string) any {
	g := parse(in)

	all := set.New[int]()
	for v := range g.adj {
		all.Add(v)
	}

	var best []int
	g.bronKerbosch(nil, all, set.New[int](), &best)

	names := make([]string, len(best))
	for i, v := range best {
		names[i] = g.names[v]
	}

	slices.Sort(names)

	return strings.Join(names, ",")
}

// graph holds the network. Each computer has a number, and adj[v] is the set
// of computers that v connects to.
type graph struct {
	names []string
	adj   []set.Set[int]
}

// parse reads one connection per line, such as "kh-tc". It numbers the
// computers in the order they first appear.
func parse(in string) graph {
	ids := map[string]int{}
	var g graph

	id := func(name string) int {
		if n, ok := ids[name]; ok {
			return n
		}

		ids[name] = len(g.names)
		g.names = append(g.names, name)
		g.adj = append(g.adj, set.New[int]())

		return ids[name]
	}

	for _, line := range input.Lines(in) {
		left, right, ok := strings.Cut(line, "-")
		if !ok {
			continue
		}

		a, b := id(left), id(right)
		g.adj[a].Add(b)
		g.adj[b].Add(a)
	}

	return g
}

// bronKerbosch finds the largest clique (a set where all members connect to
// each other) and stores it in best.
//
// r is the clique so far. p holds the computers that can still join r. x
// holds the computers that could join r but were already tried, so a clique
// from them would repeat an earlier result.
//
// The pivot u cuts the search down. Any maximal clique must hold u or a
// computer that is not a neighbor of u. So we only try the computers in p
// that are not neighbors of u.
func (g graph) bronKerbosch(r []int, p, x set.Set[int], best *[]int) {
	if p.Len() == 0 {
		if x.Len() == 0 && len(r) > len(*best) {
			// slices.Clone makes a copy, because the caller changes r later.
			*best = slices.Clone(r)
		}

		return
	}

	// Pick the pivot with the most neighbors in p, so that the fewest
	// computers are left to try.
	pivot, most := -1, -1
	for u := range p.Union(x).All() {
		if n := countIn(g.adj[u], p); n > most {
			pivot, most = u, n
		}
	}

	for v := range p.Difference(g.adj[pivot]).All() {
		// Intersect loops over its receiver, so the small neighbor set goes
		// first. append may reuse the backing array of r for each v. That is
		// safe, because each call finishes before the next v overwrites the
		// slot.
		g.bronKerbosch(append(r, v), g.adj[v].Intersect(p), g.adj[v].Intersect(x), best)

		p.Remove(v)
		x.Add(v)
	}
}

// countIn counts the members of small that are also in big. Each computer
// has few neighbors, so this is much cheaper than a full intersection.
func countIn(small, big set.Set[int]) int {
	n := 0
	for v := range small.All() {
		if big.Has(v) {
			n++
		}
	}

	return n
}
