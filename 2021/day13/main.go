// Advent of Code 2021, day 13: Transparent Origami.
// https://adventofcode.com/2021/day/13
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 13, part1, part2)
}

// on and off are the characters for a dot and an empty space in the drawing.
const (
	on  = '#'
	off = ' '
)

// part1 does the first fold only and counts the dots that are visible.
// Dots that land on top of each other count as one.
func part1(in string) any {
	dots, folds := parse(in)

	return foldAll(dots, folds[:1]).Len()
}

// part2 does all the folds and draws the dots. They spell out eight capital
// letters. Grid.String ends every row with a newline, so the last one is
// trimmed.
func part2(in string) any {
	dots, folds := parse(in)
	paper := foldAll(dots, folds)

	w, h := 0, 0
	for p := range paper.All() {
		w = max(w, p.X+1)
		h = max(h, p.Y+1)
	}

	screen := grid.New[byte](w, h)
	for i := range screen.Cells {
		screen.Cells[i] = off
	}

	for p := range paper.All() {
		screen.Set(p, on)
	}

	return strings.TrimRight(screen.String(), "\n")
}

// fold is one fold instruction. A vertical fold (x=n) has alongX set.
type fold struct {
	alongX bool
	line   int
}

// parse reads the dot positions and the fold instructions.
func parse(in string) (set.Set[grid.Point], []fold) {
	blocks := input.Blocks(in)
	dots := set.New[grid.Point]()

	for _, line := range blocks[0] {
		n := input.Ints(line)
		dots.Add(grid.P(n[0], n[1]))
	}

	var folds []fold
	for _, line := range blocks[1] {
		folds = append(folds, fold{
			alongX: strings.Contains(line, "x="),
			line:   input.Ints(line)[0],
		})
	}

	return dots, folds
}

// foldAll applies the folds in order and returns the new set of dots.
//
// A dot past the fold line moves to its mirror position on the other side.
// For a line at n, the coordinate c becomes 2n - c. Dots before the line do
// not move.
func foldAll(dots set.Set[grid.Point], folds []fold) set.Set[grid.Point] {
	for _, f := range folds {
		next := set.New[grid.Point]()

		for p := range dots.All() {
			if f.alongX && p.X > f.line {
				p.X = 2*f.line - p.X
			}

			if !f.alongX && p.Y > f.line {
				p.Y = 2*f.line - p.Y
			}

			next.Add(p)
		}

		dots = next
	}

	return dots
}
