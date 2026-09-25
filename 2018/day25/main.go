// Advent of Code 2018, day 25: Four-Dimensional Adventure.
// https://adventofcode.com/2018/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 25, part1, nil)
}

// reach is the largest Manhattan distance at which two points join the same
// constellation.
const reach = 3

// point is a position in four dimensions.
type point [4]int

// part1 counts the constellations. Every pair of points that are close
// enough are joined in a union-find structure, and each root that is left
// is one constellation.
func part1(in string) any {
	var points []point
	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		points = append(points, point{n[0], n[1], n[2], n[3]})
	}

	parent := make([]int, len(points))
	for i := range parent {
		parent[i] = i
	}

	// find returns the root of the set that holds i. It points every node
	// it passes at its grandparent, so later calls take fewer steps.
	find := func(i int) int {
		for parent[i] != i {
			parent[i] = parent[parent[i]]
			i = parent[i]
		}

		return i
	}

	count := len(points)
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			if distance(points[i], points[j]) > reach {
				continue
			}

			a, b := find(i), find(j)
			if a != b {
				parent[a] = b
				count--
			}
		}
	}

	return count
}

// distance returns the Manhattan distance between two points.
func distance(a, b point) int {
	d := 0
	for i := range a {
		d += mathx.Abs(a[i] - b[i])
	}

	return d
}
