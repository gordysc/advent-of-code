// Advent of Code 2022, day 12: Hill Climbing Algorithm.
// https://adventofcode.com/2022/day/12
package main

import (
	"aoc/lib/grid"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 12, part1, part2)
}

// part1 finds the fewest steps from the start S to the best signal E.
func part1(in string) any {
	heights, start, end := parse(in)

	steps, _ := search.BFS(start, uphill(heights), func(p grid.Point) bool {
		return p == end
	})

	return steps
}

// part2 finds the fewest steps from any square of height 'a' to E.
//
// A search from each 'a' square is slow. Instead, the search goes backward
// from E, with the climb rule turned around, and stops at the first 'a'
// square. Breadth-first search finds the nearest one first.
func part2(in string) any {
	heights, _, end := parse(in)

	steps, _ := search.BFS(end, downhill(heights), func(p grid.Point) bool {
		return heights.At(p) == 'a'
	})

	return steps
}

// uphill returns the moves from a square: the neighbors that are at most one
// higher.
func uphill(heights grid.Grid[byte]) func(grid.Point) []grid.Point {
	return func(p grid.Point) []grid.Point {
		var next []grid.Point

		for _, n := range heights.Neighbors4(p) {
			if int(heights.At(n)) <= int(heights.At(p))+1 {
				next = append(next, n)
			}
		}

		return next
	}
}

// downhill returns the backward moves from a square: the neighbors from which
// a forward move to this square is possible.
func downhill(heights grid.Grid[byte]) func(grid.Point) []grid.Point {
	return func(p grid.Point) []grid.Point {
		var next []grid.Point

		for _, n := range heights.Neighbors4(p) {
			if int(heights.At(p)) <= int(heights.At(n))+1 {
				next = append(next, n)
			}
		}

		return next
	}
}

// parse reads the height map. It returns the grid with S changed to 'a' and E
// changed to 'z', and the positions of S and E.
func parse(in string) (grid.Grid[byte], grid.Point, grid.Point) {
	heights := grid.Parse(in)

	start, _ := heights.FindByte('S')
	end, _ := heights.FindByte('E')

	heights.Set(start, 'a')
	heights.Set(end, 'z')

	return heights, start, end
}
