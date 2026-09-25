// Advent of Code 2023, day 3: Gear Ratios.
// https://adventofcode.com/2023/day/3
package main

import (
	"aoc/lib/grid"
	"aoc/lib/strx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 3, part1, part2)
}

// part1 adds the part numbers. A part number is a number with a symbol in
// one of the cells around it, diagonals included.
func part1(in string) any {
	total := 0

	for _, n := range findNumbers(grid.Parse(in)) {
		if len(n.symbols) > 0 {
			total += n.value
		}
	}

	return total
}

// part2 adds the gear ratios. A gear is a "*" next to exactly two part
// numbers, and its ratio is the product of these two numbers.
//
// Each number already knows the symbols next to it. So the part turns that
// list around: it collects, for each "*", the numbers that touch it.
func part2(in string) any {
	g := grid.Parse(in)
	touching := map[grid.Point][]int{}

	for _, n := range findNumbers(g) {
		for _, p := range n.symbols {
			if g.At(p) == '*' {
				touching[p] = append(touching[p], n.value)
			}
		}
	}

	total := 0

	for _, values := range touching {
		if len(values) == 2 {
			total += values[0] * values[1]
		}
	}

	return total
}

// number is one number in the schematic, with the positions of all symbols
// next to it.
type number struct {
	value   int
	symbols []grid.Point
}

// findNumbers reads every number in the grid from left to right, row by row.
// For each number it also records the symbols next to any of its digits.
func findNumbers(g grid.Grid[byte]) []number {
	var numbers []number

	for y := range g.H {
		x := 0

		for x < g.W {
			if !strx.IsDigit(g.At(grid.P(x, y))) {
				x++
				continue
			}

			// Read all digits of the number and remember where it starts
			// and ends, so the border around it is known.
			start := x
			value := 0

			for x < g.W && strx.IsDigit(g.At(grid.P(x, y))) {
				value = value*10 + int(g.At(grid.P(x, y))-'0')
				x++
			}

			numbers = append(numbers, number{value, symbolsAround(g, start, x-1, y)})
		}
	}

	return numbers
}

// symbolsAround gives the symbols in the box of cells around a number on row
// y, from column x0 to column x1. A symbol is any cell that is not a digit
// and not a ".". The box includes the digits, but they are never symbols.
func symbolsAround(g grid.Grid[byte], x0, x1, y int) []grid.Point {
	var symbols []grid.Point

	for yy := y - 1; yy <= y+1; yy++ {
		for xx := x0 - 1; xx <= x1+1; xx++ {
			c, ok := g.Get(grid.P(xx, yy))
			if ok && c != '.' && !strx.IsDigit(c) {
				symbols = append(symbols, grid.P(xx, yy))
			}
		}
	}

	return symbols
}
