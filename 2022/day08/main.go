// Advent of Code 2022, day 8: Treetop Tree House.
// https://adventofcode.com/2022/day/8
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 8, part1, part2)
}

// part1 counts the trees that you can see from outside the grid.
//
// A tree is visible when, in at least one of the four directions, all trees
// between it and the edge are shorter than it.
func part1(in string) any {
	trees := grid.Parse(in)

	count := 0

	for p := range trees.Points() {
		for _, dir := range grid.Dirs4 {
			_, blocked := look(trees, p, dir)
			if !blocked {
				count++
				break
			}
		}
	}

	return count
}

// part2 finds the highest scenic score of all the trees.
//
// The scenic score of a tree is the product of its viewing distances in the
// four directions.
func part2(in string) any {
	trees := grid.Parse(in)

	best := 0

	for p := range trees.Points() {
		score := 1
		for _, dir := range grid.Dirs4 {
			dist, _ := look(trees, p, dir)
			score *= dist
		}

		best = max(best, score)
	}

	return best
}

// look walks from the tree at p in direction dir. It returns the number of
// trees that the tree at p can see, and true when a tree of the same height
// or taller stops the view before the edge.
//
// The digits are bytes, but their order is the same as the order of the
// heights, so the code compares the bytes directly.
func look(trees grid.Grid[byte], p, dir grid.Point) (int, bool) {
	height := trees.At(p)

	dist := 0

	for q := p.Add(dir); trees.InBounds(q); q = q.Add(dir) {
		dist++

		if trees.At(q) >= height {
			return dist, true
		}
	}

	return dist, false
}
