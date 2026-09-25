// Advent of Code 2018, day 13: Mine Cart Madness.
// https://adventofcode.com/2018/day/13
package main

import (
	"fmt"
	"slices"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 13, part1, part2)
}

// cart is one mine cart. turns counts the intersections it has passed, which
// picks its next turn: left, straight, right, and then again from the start.
type cart struct {
	pos, dir grid.Point
	turns    int
	crashed  bool
}

// cartDirs maps each cart symbol to the direction the cart faces.
var cartDirs = map[byte]grid.Point{
	'^': grid.Up,
	'v': grid.Down,
	'<': grid.Left,
	'>': grid.Right,
}

// part1 returns the location of the first crash as "x,y".
func part1(in string) any {
	track, carts := parse(in)

	for {
		if crashes := tick(track, carts); len(crashes) > 0 {
			return format(crashes[0])
		}
	}
}

// part2 removes the carts that crash and returns the location of the one
// cart that is left at the end of the tick of the last crash. The first
// example has only two carts, which crash into each other, so part 2 has no
// answer for it.
func part2(in string) any {
	track, carts := parse(in)

	for {
		tick(track, carts)

		var left []*cart
		for _, c := range carts {
			if !c.crashed {
				left = append(left, c)
			}
		}

		switch len(left) {
		case 0:
			return nil
		case 1:
			return format(left[0].pos)
		}

		carts = left
	}
}

// parse reads the track and the carts on it. The rows can have different
// lengths, so the track stays as text rows. Each cart is replaced by the
// straight piece of track under it.
func parse(in string) ([][]byte, []*cart) {
	var track [][]byte
	var carts []*cart

	for y, line := range input.Lines(in) {
		row := []byte(line)

		for x, c := range row {
			dir, ok := cartDirs[c]
			if !ok {
				continue
			}

			carts = append(carts, &cart{pos: grid.P(x, y), dir: dir})

			row[x] = '-'
			if dir.X == 0 {
				row[x] = '|'
			}
		}

		track = append(track, row)
	}

	return track, carts
}

// tick moves every cart one step, in reading order, and returns the places
// where carts crashed. A cart that crashes stops at once, and so does the
// cart it hits, even if that cart has not moved yet in this tick.
func tick(track [][]byte, carts []*cart) []grid.Point {
	slices.SortFunc(carts, func(a, b *cart) int {
		if a.pos.Y != b.pos.Y {
			return a.pos.Y - b.pos.Y
		}

		return a.pos.X - b.pos.X
	})

	var crashes []grid.Point

	for _, c := range carts {
		if c.crashed {
			continue
		}

		c.move(track)

		for _, other := range carts {
			if other != c && !other.crashed && other.pos == c.pos {
				c.crashed = true
				other.crashed = true
				crashes = append(crashes, c.pos)
			}
		}
	}

	return crashes
}

// move takes the cart one step forward and turns it to match the track
// under its new position.
func (c *cart) move(track [][]byte) {
	c.pos = c.pos.Add(c.dir)

	switch track[c.pos.Y][c.pos.X] {
	case '/':
		// Up turns right, right turns up, down turns left, left turns down.
		c.dir = grid.P(-c.dir.Y, -c.dir.X)
	case '\\':
		// Up turns left, left turns up, down turns right, right turns down.
		c.dir = grid.P(c.dir.Y, c.dir.X)
	case '+':
		switch c.turns % 3 {
		case 0:
			c.dir = c.dir.TurnLeft()
		case 2:
			c.dir = c.dir.TurnRight()
		}

		c.turns++
	}
}

// format writes a location the way the puzzle wants it: "x,y".
func format(p grid.Point) string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}
