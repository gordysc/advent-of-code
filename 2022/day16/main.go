// Advent of Code 2022, day 16: Proboscidea Volcanium.
// https://adventofcode.com/2022/day/16
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 16, part1, part2)
}

// part1 finds the most pressure that one person can release in 30 minutes.
//
// Most valves have a flow rate of zero, so opening them is useless. We keep
// only the valves with flow (about 15 in real inputs) and the shortest travel
// time between each pair. Then a depth-first search tries every order to
// open them in, and it stops when the time is used up.
func part1(in string) any {
	return slices.Max(bestPerSet(parse(in), 30))
}

// part2 finds the most pressure that you and an elephant can release in 26
// minutes, when you both work at the same time.
//
// You and the elephant never open the same valve, so you open two disjoint
// sets of valves. We find the best pressure for each set of opened valves in
// 26 minutes. Then we try each set for you, and give the elephant the best
// set among the valves that are left.
func part2(in string) any {
	best := bestPerSet(parse(in), 26)

	// Change best[set] into the best pressure for any subset of set. We go
	// through the sets in increasing order, so each smaller set is final
	// before a larger set reads it.
	for set := range best {
		for bit := 1; bit < len(best); bit <<= 1 {
			if set&bit != 0 {
				best[set] = max(best[set], best[set^bit])
			}
		}
	}

	all := len(best) - 1
	total := 0

	for set := range best {
		total = max(total, best[set]+best[all^set])
	}

	return total
}

// network holds the valves with a flow rate above zero. Valve i is bit i in a
// set of opened valves.
type network struct {
	rates []int
	// dist[i][j] is the number of minutes to walk from valve i to valve j.
	// Index len(rates) is the start valve AA.
	dist [][]int
}

// parse reads the valves and tunnels. It keeps only the valves with flow and
// the start valve, and finds the walking time between each pair of them.
func parse(in string) network {
	rates := map[string]int{}
	tunnels := map[string][]string{}
	var useful []string

	for _, line := range input.Lines(in) {
		// A line looks like "Valve AA has flow rate=0; tunnels lead to
		// valves DD, II, BB". The names of the next valves start at field 9.
		fields := strings.Fields(line)
		name := fields[1]
		rates[name] = input.Ints(line)[0]

		for _, f := range fields[9:] {
			tunnels[name] = append(tunnels[name], strings.TrimSuffix(f, ","))
		}

		if rates[name] > 0 {
			useful = append(useful, name)
		}
	}

	// The start valve goes last, so the valve bits stay 0 to n-1.
	names := append(useful, "AA")
	next := func(name string) []string { return tunnels[name] }

	net := network{dist: make([][]int, len(names))}
	for _, name := range useful {
		net.rates = append(net.rates, rates[name])
	}

	for i, from := range names {
		steps := search.Flood(from, next)

		net.dist[i] = make([]int, len(names))
		for j, to := range names {
			net.dist[i][j] = steps[to]
		}
	}

	return net
}

// bestPerSet returns, for each set of opened valves, the most pressure one
// worker can release in the given minutes by opening exactly that set. Sets
// that the worker cannot open in time stay at 0.
func bestPerSet(net network, minutes int) []int {
	n := len(net.rates)
	best := make([]int, 1<<n)

	// visit is recursive, so it must be declared before it is assigned.
	var visit func(at, left, opened, pressure int)
	visit = func(at, left, opened, pressure int) {
		best[opened] = max(best[opened], pressure)

		for next := range n {
			if opened&(1<<next) != 0 {
				continue
			}

			// Walk to the valve and spend one minute to open it. It then
			// releases pressure in each minute that is left.
			remaining := left - net.dist[at][next] - 1
			if remaining <= 0 {
				continue
			}

			visit(next, remaining, opened|1<<next, pressure+remaining*net.rates[next])
		}
	}

	visit(n, minutes, 0, 0)

	return best
}
