// Advent of Code 2021, day 23: Amphipod.
// https://adventofcode.com/2021/day/23
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 23, part1, part2)
}

// part1 finds the least energy to sort the amphipods in rooms of depth 2.
//
// Each burrow state is a node in a graph, and each legal move is an edge with
// its energy cost. Dijkstra's algorithm finds the cheapest path to the
// sorted state.
func part1(in string) any {
	return solve(parseRows(in))
}

// part2 inserts two more rows into each room and solves again.
//
// The folded part of the diagram is "DCBA" and "DBAC". It goes between the
// first and second rows of the input.
func part2(in string) any {
	rows := parseRows(in)
	rows = append([]string{rows[0], "DCBA", "DBAC"}, rows[1:]...)

	return solve(rows)
}

// hallLen is the number of cells in the hallway.
const hallLen = 11

// energy is the cost of one step for each amphipod type, A to D.
var energy = [4]int{1, 10, 100, 1000}

// door returns the hallway cell outside room r.
func door(r int) int {
	return 2 + 2*r
}

// isDoor reports if amphipods may not stop at hallway cell h.
func isDoor(h int) bool {
	return h >= 2 && h <= 8 && h%2 == 0
}

// burrow holds the room depth, which is the same for all states in a search.
//
// A state is a string: the 11 hallway cells, then the cells of room 0 from
// top to bottom, then room 1, and so on. An empty cell is '.'.
type burrow struct {
	depth int
}

// solve returns the least energy to sort the amphipods. rows holds the
// amphipods in each room row, from top to bottom, as four letters each.
func solve(rows []string) any {
	b := burrow{depth: len(rows)}

	start := []byte(strings.Repeat(".", hallLen))
	goal := []byte(strings.Repeat(".", hallLen))

	for r := range 4 {
		for d := range b.depth {
			start = append(start, rows[d][r])
			goal = append(goal, byte('A'+r))
		}
	}

	target := string(goal)
	cost, ok := search.Dijkstra(string(start), b.moves, func(s string) bool {
		return s == target
	})

	if !ok {
		return nil
	}

	return cost
}

// cell returns the index in the state of room r at depth d (0 is the top).
func (b burrow) cell(r, d int) int {
	return hallLen + r*b.depth + d
}

// moves returns all the legal moves from state s.
//
// An amphipod leaves its room only to stop in the hallway, and leaves the
// hallway only to go into its own room when that room holds only its own
// type. If an amphipod can go into its own room, that move is always part of
// a best solution, so only that move is returned.
func (b burrow) moves(s string) []search.Edge[string] {
	// Moves from the hallway into a room.
	for h := range hallLen {
		if s[h] == '.' {
			continue
		}

		t := int(s[h] - 'A')
		d, ok := b.freeSlot(s, t)

		if !ok || !b.hallClear(s, h, door(t)) {
			continue
		}

		steps := mathx.Abs(h-door(t)) + d + 1
		next := []byte(s)
		next[b.cell(t, d)] = s[h]
		next[h] = '.'

		return []search.Edge[string]{{To: string(next), Cost: steps * energy[t]}}
	}

	// Moves from the top of a room into the hallway.
	var edges []search.Edge[string]

	for r := range 4 {
		d, ok := b.topToMove(s, r)
		if !ok {
			continue
		}

		c := b.cell(r, d)
		t := int(s[c] - 'A')

		for h := range hallLen {
			if isDoor(h) || !b.hallClear(s, door(r), h) || s[h] != '.' {
				continue
			}

			steps := d + 1 + mathx.Abs(h-door(r))
			next := []byte(s)
			next[h] = s[c]
			next[c] = '.'

			edges = append(edges, search.Edge[string]{To: string(next), Cost: steps * energy[t]})
		}
	}

	return edges
}

// hallClear reports if the hallway cells between from and to are empty. The
// from cell is not checked, but the to cell is.
func (b burrow) hallClear(s string, from, to int) bool {
	step := 1
	if to < from {
		step = -1
	}

	for h := from + step; h != to+step; h += step {
		if s[h] != '.' {
			return false
		}
	}

	return true
}

// freeSlot returns the deepest empty cell of room t, if the room holds only
// amphipods of type t.
func (b burrow) freeSlot(s string, t int) (int, bool) {
	own := byte('A' + t)
	slot := -1

	for d := range b.depth {
		switch s[b.cell(t, d)] {
		case '.':
			slot = d
		case own:
		default:
			return 0, false
		}
	}

	return slot, slot >= 0
}

// topToMove returns the depth of the top amphipod in room r, if that
// amphipod must move. It must move when it, or one below it, is in the wrong
// room.
func (b burrow) topToMove(s string, r int) (int, bool) {
	own := byte('A' + r)

	for d := range b.depth {
		if s[b.cell(r, d)] == '.' {
			continue
		}

		for below := d; below < b.depth; below++ {
			if s[b.cell(r, below)] != own {
				return d, true
			}
		}

		return 0, false
	}

	return 0, false
}

// parseRows returns the letters in each room row of the diagram, from top to
// bottom.
func parseRows(in string) []string {
	var rows []string

	for _, line := range input.Lines(in) {
		letters := strings.Map(func(r rune) rune {
			if r >= 'A' && r <= 'D' {
				return r
			}

			return -1
		}, line)

		if letters != "" {
			rows = append(rows, letters)
		}
	}

	return rows
}
