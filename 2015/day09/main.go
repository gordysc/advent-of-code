// Advent of Code 2015, day 9: All in a Single Night.
// https://adventofcode.com/2015/day/9
package main

import (
	"math"
	"strings"

	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 9, part1, part2)
}

// part1 finds the shortest route that visits every location exactly once.
// The route can start and end anywhere; it does not return to the start.
func part1(in string) any {
	shortest, _ := routeLengths(parseMap(in))

	return shortest
}

// part2 finds the longest such route instead.
func part2(in string) any {
	_, longest := routeLengths(parseMap(in))

	return longest
}

// distMap holds the distance between every pair of locations. Locations are
// numbered in the order they first appear in the input, so the table is a
// small square matrix and the search can work with ints instead of strings.
type distMap struct {
	names []string
	dist  [][]int
}

// parseMap reads lines like "London to Dublin = 464" into a distMap. Every
// pair appears once in the input, but the table stores both directions since
// the routes are undirected.
func parseMap(in string) distMap {
	index := map[string]int{}
	m := distMap{}

	// id returns the number for a location, adding it to the table if this is
	// the first time it appears. Growing the matrix one row and column at a
	// time keeps the parse to a single pass.
	id := func(name string) int {
		if i, ok := index[name]; ok {
			return i
		}

		i := len(m.names)
		index[name] = i
		m.names = append(m.names, name)

		for r := range m.dist {
			m.dist[r] = append(m.dist[r], 0)
		}

		m.dist = append(m.dist, make([]int, i+1))

		return i
	}

	for _, line := range input.Lines(in) {
		fields := strings.Fields(line)
		a, b := id(fields[0]), id(fields[2])
		d := input.Int(fields[4])

		m.dist[a][b] = d
		m.dist[b][a] = d
	}

	return m
}

// routeLengths tries every ordering of the locations and returns the shortest
// and longest total distance. The inputs have at most eight locations, so the
// 8! = 40320 permutations are cheap to check by brute force.
func routeLengths(m distMap) (shortest, longest int) {
	shortest = math.MaxInt
	longest = 0

	order := make([]int, len(m.names))
	for i := range order {
		order[i] = i
	}

	for route := range slicesx.Permutations(order) {
		total := 0

		for i := 1; i < len(route); i++ {
			total += m.dist[route[i-1]][route[i]]
		}

		shortest = min(shortest, total)
		longest = max(longest, total)
	}

	return shortest, longest
}
