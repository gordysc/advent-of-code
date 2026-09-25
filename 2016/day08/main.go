// Advent of Code 2016, day 8: Two-Factor Authentication.
// https://adventofcode.com/2016/day/8
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 8, part1, part2)
}

// width and height are the size of the screen in the puzzle. The worked
// example in the puzzle text uses a 7x3 screen instead. Its pixels end up in
// different places on this screen, but the lit count is the same.
const (
	width  = 50
	height = 6
)

// on and off are the characters for the two states of a pixel.
const (
	on  = '#'
	off = ' '
)

// part1 counts the lit pixels after every instruction has run.
func part1(in string) any {
	return run(in).Count(func(b byte) bool { return b == on })
}

// part2 draws the screen. The lit pixels spell out the code as capital letters.
// Grid.String ends every row with a newline, so the last one is trimmed.
func part2(in string) any {
	return strings.TrimRight(run(in).String(), "\n")
}

// run starts with every pixel off and applies each instruction in turn.
func run(in string) grid.Grid[byte] {
	screen := grid.New[byte](width, height)
	for i := range screen.Cells {
		screen.Cells[i] = off
	}

	for _, line := range input.Lines(in) {
		nums := input.Ints(line)

		switch {
		case strings.HasPrefix(line, "rect"):
			rect(screen, nums[0], nums[1])
		case strings.HasPrefix(line, "rotate row"):
			rotateRow(screen, nums[0], nums[1])
		case strings.HasPrefix(line, "rotate column"):
			rotateColumn(screen, nums[0], nums[1])
		}
	}

	return screen
}

// rect turns on every pixel in the w×h rectangle at the top-left corner.
func rect(screen grid.Grid[byte], w, h int) {
	for y := range h {
		for x := range w {
			screen.Set(grid.Point{X: x, Y: y}, on)
		}
	}
}

// rotateRow shifts row y right by n pixels. Pixels that fall off the right
// edge come back on the left.
func rotateRow(screen grid.Grid[byte], y, n int) {
	row := make([]byte, screen.W)
	for x := range row {
		row[x] = screen.At(grid.Point{X: x, Y: y})
	}

	for x, b := range row {
		screen.Set(grid.Point{X: (x + n) % screen.W, Y: y}, b)
	}
}

// rotateColumn shifts column x down by n pixels. Pixels that fall off the
// bottom edge come back at the top.
func rotateColumn(screen grid.Grid[byte], x, n int) {
	col := make([]byte, screen.H)
	for y := range col {
		col[y] = screen.At(grid.Point{X: x, Y: y})
	}

	for y, b := range col {
		screen.Set(grid.Point{X: x, Y: (y + n) % screen.H}, b)
	}
}
