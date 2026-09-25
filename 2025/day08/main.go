// Advent of Code 2025, day 8: Playground.
// https://adventofcode.com/2025/day/8
package main

import (
	"cmp"
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 8, part1, part2)
}

// part1 connects the closest pairs of junction boxes, and returns the product
// of the sizes of the three largest circuits.
//
// The real input has 1000 boxes and makes 1000 connections. The example has
// only 20 boxes and makes 10 connections. So a small input (fewer than 1000
// boxes) uses 10. A pair that is already in the same circuit still counts as a
// connection, but it does not change the circuits.
func part1(in string) any {
	boxes := parse(in)
	if len(boxes) < 3 {
		return nil
	}

	connections := 1000
	if len(boxes) < 1000 {
		connections = 10
	}

	pairs := sortedPairs(boxes)
	circuits := newUnionFind(len(boxes))

	for _, p := range pairs[:min(connections, len(pairs))] {
		circuits.union(p.a, p.b)
	}

	// Each root of the union-find is one circuit, and its size is the number
	// of boxes in it.
	var sizes []int
	for i := range boxes {
		if circuits.find(i) == i {
			sizes = append(sizes, circuits.size[i])
		}
	}

	if len(sizes) < 3 {
		return nil
	}

	slices.Sort(sizes)
	slices.Reverse(sizes)

	return sizes[0] * sizes[1] * sizes[2]
}

// part2 connects the closest pairs until all boxes are in one circuit. It
// returns the product of the X coordinates of the last two boxes joined.
func part2(in string) any {
	boxes := parse(in)
	if len(boxes) < 2 {
		return nil
	}

	circuits := newUnionFind(len(boxes))
	count := len(boxes)

	for _, p := range sortedPairs(boxes) {
		if !circuits.union(p.a, p.b) {
			continue
		}

		// Each join makes two circuits into one. The last join leaves one.
		count--
		if count == 1 {
			return boxes[p.a][0] * boxes[p.b][0]
		}
	}

	return nil
}

// box is the X, Y and Z position of a junction box.
type box [3]int

// pair is two boxes, by index, and the square of their distance.
type pair struct {
	a, b int
	dist int
}

// parse reads one "X,Y,Z" box per line.
func parse(in string) []box {
	var boxes []box

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		if len(n) < 3 {
			continue
		}

		boxes = append(boxes, box{n[0], n[1], n[2]})
	}

	return boxes
}

// sortedPairs returns all pairs of boxes, closest first.
//
// It compares the square of the distance and does not take the square root.
// The order is the same, and integers have no rounding errors. With 1000 boxes
// there are about 500,000 pairs, so sorting all of them is fast enough.
func sortedPairs(boxes []box) []pair {
	pairs := make([]pair, 0, len(boxes)*(len(boxes)-1)/2)

	for i := range boxes {
		for j := i + 1; j < len(boxes); j++ {
			d := 0
			for k := range 3 {
				diff := boxes[i][k] - boxes[j][k]
				d += diff * diff
			}

			pairs = append(pairs, pair{i, j, d})
		}
	}

	slices.SortFunc(pairs, func(x, y pair) int {
		return cmp.Compare(x.dist, y.dist)
	})

	return pairs
}

// unionFind keeps track of which boxes are in the same circuit. Each circuit
// is a tree, and the root of the tree names the circuit.
type unionFind struct {
	// parent gives the next box toward the root. A root is its own parent.
	parent []int
	// size gives the number of boxes in a circuit. It is correct only for roots.
	size []int
}

// newUnionFind returns n circuits that each hold one box.
func newUnionFind(n int) *unionFind {
	u := &unionFind{parent: make([]int, n), size: make([]int, n)}

	for i := range n {
		u.parent[i] = i
		u.size[i] = 1
	}

	return u
}

// find returns the root of the circuit that holds box x.
//
// It also points each box on the path straight to the root (path
// compression), so later calls are faster. The method has a pointer receiver
// so that these changes stay in the struct.
func (u *unionFind) find(x int) int {
	for u.parent[x] != x {
		u.parent[x] = u.parent[u.parent[x]]
		x = u.parent[x]
	}

	return x
}

// union joins the circuits of boxes a and b. It returns false when they are
// already in the same circuit.
//
// The smaller tree goes below the root of the larger tree. This keeps the
// trees short.
func (u *unionFind) union(a, b int) bool {
	ra, rb := u.find(a), u.find(b)
	if ra == rb {
		return false
	}

	if u.size[ra] < u.size[rb] {
		ra, rb = rb, ra
	}

	u.parent[rb] = ra
	u.size[ra] += u.size[rb]

	return true
}
