// Advent of Code 2015, day 13: Knights of the Dinner Table.
// https://adventofcode.com/2015/day/13
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
	runner.Run(2015, 13, part1, part2)
}

// part1 finds the seating order around the round table that gives the largest
// total change in happiness.
func part1(in string) any {
	return bestSeating(parseTable(in))
}

// part2 adds yourself to the table first. You gain and lose nothing next to
// anyone, and nobody gains or loses anything next to you.
func part2(in string) any {
	t := parseTable(in)
	t.add("Me")

	return bestSeating(t)
}

// table holds the happiness change for every pair of neighbours. Guests are
// numbered in the order they first appear in the input, so the table is a
// small square matrix and the search can work with ints instead of strings.
//
// pair[a][b] is the change for the pair as a whole: what a gains next to b
// plus what b gains next to a. Two neighbours always affect each other, so the
// search only ever needs the combined value.
type table struct {
	index map[string]int
	names []string
	pair  [][]int
}

// add returns the number for a guest, adding them to the table if this is the
// first time they appear. Growing the matrix one row and column at a time
// keeps the parse to a single pass. New cells start at zero, which is exactly
// what part 2 needs for the extra guest.
//
// The receiver is a pointer because append can return a new slice, which add
// stores back into the table. With a value receiver those assignments would
// go to a copy and the caller would never see the new guest.
func (t *table) add(name string) int {
	if i, ok := t.index[name]; ok {
		return i
	}

	i := len(t.names)
	t.index[name] = i
	t.names = append(t.names, name)

	for r := range t.pair {
		t.pair[r] = append(t.pair[r], 0)
	}

	t.pair = append(t.pair, make([]int, i+1))

	return i
}

// parseTable reads lines like
// "Alice would gain 54 happiness units by sitting next to Bob." into a table.
func parseTable(in string) table {
	t := table{index: map[string]int{}}

	for _, line := range input.Lines(in) {
		fields := strings.Fields(line)

		// The second name is the last word and has the full stop of the sentence
		// attached to it.
		a := t.add(fields[0])
		b := t.add(strings.TrimSuffix(fields[len(fields)-1], "."))

		units := input.Int(fields[3])
		if fields[2] == "lose" {
			units = -units
		}

		// Each line gives one direction only. Adding it to both cells builds up
		// the combined value for the pair once both lines have been read.
		t.pair[a][b] += units
		t.pair[b][a] += units
	}

	return t
}

// bestSeating tries the seating orders and returns the largest total change in
// happiness.
//
// The table is round, so turning a whole arrangement by one seat gives the
// same neighbours. Guest 0 therefore stays in the first seat and only the
// other guests are permuted. That cuts the work from n! to (n-1)! orders, at
// most 8! = 40320 for the nine guests of part 2.
func bestSeating(t table) int {
	if len(t.names) == 0 {
		return 0
	}

	others := make([]int, 0, len(t.names)-1)
	for i := 1; i < len(t.names); i++ {
		others = append(others, i)
	}

	best := math.MinInt

	for order := range slicesx.Permutations(others) {
		// Start at guest 0 and walk around the table. The last step closes the
		// circle back to guest 0.
		total := 0
		prev := 0

		for _, guest := range order {
			total += t.pair[prev][guest]
			prev = guest
		}

		total += t.pair[prev][0]

		best = max(best, total)
	}

	return best
}
