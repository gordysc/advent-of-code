// Advent of Code 2018, day 10: The Stars Align.
// https://adventofcode.com/2018/day/10
package main

import (
	"math"
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 10, part1, part2)
}

// on and off are the characters for a cell with a point and an empty cell.
const (
	on  = '#'
	off = ' '
)

// star is a point of light with its position and its velocity per second.
type star struct {
	pos, vel grid.Point
}

// part1 draws the stars at the moment they spell the message.
// Grid.String ends every row with a newline, so the last one is trimmed.
func part1(in string) any {
	stars := parse(in)
	t := align(stars)

	return strings.TrimRight(draw(stars, t).String(), "\n")
}

// part2 returns how many seconds the stars take to spell the message.
func part2(in string) any {
	return align(parse(in))
}

// parse reads one star from each line.
func parse(in string) []star {
	var stars []star

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		stars = append(stars, star{pos: grid.P(n[0], n[1]), vel: grid.P(n[2], n[3])})
	}

	return stars
}

// align finds the second when the stars spell the message. The stars move
// toward each other, then apart again. The message shows when they are
// closest, so this is the second when the bounding box is at its lowest.
func align(stars []star) int {
	t := 0

	for height(stars, t+1) < height(stars, t) {
		t++
	}

	return t
}

// height returns the height of the bounding box of the stars after t seconds.
func height(stars []star, t int) int {
	top, bottom := math.MaxInt, math.MinInt

	for _, s := range stars {
		y := s.pos.Y + s.vel.Y*t
		top, bottom = min(top, y), max(bottom, y)
	}

	return bottom - top + 1
}

// draw puts the stars after t seconds on a grid that just fits them.
func draw(stars []star, t int) grid.Grid[byte] {
	pos := make([]grid.Point, len(stars))
	lo := grid.P(math.MaxInt, math.MaxInt)
	hi := grid.P(math.MinInt, math.MinInt)

	for i, s := range stars {
		p := s.pos.Add(s.vel.Scale(t))
		pos[i] = p
		lo = grid.P(min(lo.X, p.X), min(lo.Y, p.Y))
		hi = grid.P(max(hi.X, p.X), max(hi.Y, p.Y))
	}

	sky := grid.New[byte](hi.X-lo.X+1, hi.Y-lo.Y+1)
	for i := range sky.Cells {
		sky.Cells[i] = off
	}

	for _, p := range pos {
		sky.Set(p.Sub(lo), on)
	}

	return sky
}
