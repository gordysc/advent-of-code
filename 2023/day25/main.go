// Advent of Code 2023, day 25: Snowverload.
// https://adventofcode.com/2023/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 25, part1, nil)
}

// part1 cuts three wires to split the components into two groups, and returns
// the product of the two group sizes.
//
// The three wires are a minimum cut. By the max-flow min-cut theorem, the
// number of paths without a shared wire between two components is equal to
// the smallest number of wires that separates them. So we fix one source
// component and try each other component as the sink. A sink in the same group
// as the source has four or more such paths. A sink in the other group has
// exactly three. After the three paths use up their wires, the components that
// the source can still reach are its group.
func part1(in string) any {
	g := parse(in)
	if len(g.adj) < 2 {
		return nil
	}

	for sink := 1; sink < len(g.adj); sink++ {
		flow := make([]int, len(g.head))

		// Stop after a fourth path: then this sink is in the source's group.
		for paths := 0; paths <= 3; paths++ {
			found, reached := g.augment(0, sink, flow)
			if found {
				continue
			}

			// No more paths. If there were exactly three, the source's group
			// is the set of components it could still reach.
			if paths == 3 {
				return reached * (len(g.adj) - reached)
			}

			break
		}
	}

	return nil
}

// graph stores the wires as arcs. Each wire becomes two arcs, one in each
// direction, at indexes 2k and 2k+1. So arc a^1 is always the reverse of arc a.
type graph struct {
	// adj lists, for each component, the arcs that leave it.
	adj [][]int
	// head gives the component at the end of each arc.
	head []int
}

// parse reads the wiring diagram. It numbers the components in the order they
// first appear.
func parse(in string) graph {
	ids := map[string]int{}
	var g graph

	id := func(name string) int {
		if n, ok := ids[name]; ok {
			return n
		}

		ids[name] = len(g.adj)
		g.adj = append(g.adj, nil)

		return ids[name]
	}

	for _, line := range input.Lines(in) {
		// strings.Cut splits at the first ":" and returns false if there is none.
		left, right, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		from := id(strings.TrimSpace(left))

		for _, name := range strings.Fields(right) {
			to := id(name)

			g.adj[from] = append(g.adj[from], len(g.head))
			g.head = append(g.head, to)
			g.adj[to] = append(g.adj[to], len(g.head))
			g.head = append(g.head, from)
		}
	}

	return g
}

// augment looks for one more path from source to sink with a breadth-first
// search, and adds it to flow. It returns true when it found a path. When it
// finds none, it also returns the number of components that source can reach.
//
// Each wire has room for one unit of flow in either direction. flow[a] is the
// flow along arc a, and flow[a^1] is always its negative. So arc a has room
// for more when flow[a] < 1. This also lets a new path cancel the flow of an
// old path on the same wire.
func (g graph) augment(source, sink int, flow []int) (bool, int) {
	// via holds the arc that the search used to get to each component.
	// -1 means not reached yet.
	via := make([]int, len(g.adj))
	for i := range via {
		via[i] = -1
	}

	queue := []int{source}

	// The source has no arc into it, so mark it with any value except -1.
	via[source] = len(g.head)

	for i := 0; i < len(queue); i++ {
		node := queue[i]

		for _, a := range g.adj[node] {
			next := g.head[a]
			if flow[a] >= 1 || via[next] != -1 {
				continue
			}

			via[next] = a
			queue = append(queue, next)
		}
	}

	if via[sink] == -1 {
		return false, len(queue)
	}

	// Walk back from the sink and push one unit along each arc of the path.
	for node := sink; node != source; {
		a := via[node]
		flow[a]++
		flow[a^1]--
		node = g.head[a^1]
	}

	return true, 0
}
