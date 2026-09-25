// Advent of Code 2018, day 3: No Matter How You Slice It.
// https://adventofcode.com/2018/day/3
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 3, part1, part2)
}

// size is the width and height of the fabric, in inches.
const size = 1000

// claim is one Elf's rectangle of fabric.
type claim struct {
	id, x, y, w, h int
}

// part1 counts the square inches that two or more claims cover.
func part1(in string) any {
	fabric := cover(parse(in))
	overlaps := 0

	for _, n := range fabric {
		if n > 1 {
			overlaps++
		}
	}

	return overlaps
}

// part2 finds the one claim where every square inch is covered only once.
func part2(in string) any {
	claims := parse(in)
	fabric := cover(claims)

	for _, c := range claims {
		if alone(c, fabric) {
			return c.id
		}
	}

	return nil
}

// parse reads lines such as "#1 @ 1,3: 4x4" into claims.
func parse(in string) []claim {
	var claims []claim

	for _, line := range input.Lines(in) {
		n := input.UInts(line)
		claims = append(claims, claim{id: n[0], x: n[1], y: n[2], w: n[3], h: n[4]})
	}

	return claims
}

// cover counts, for each square inch of fabric, how many claims include it.
// The fabric is stored row by row in a flat slice.
func cover(claims []claim) []int {
	fabric := make([]int, size*size)

	for _, c := range claims {
		for y := c.y; y < c.y+c.h; y++ {
			for x := c.x; x < c.x+c.w; x++ {
				fabric[y*size+x]++
			}
		}
	}

	return fabric
}

// alone reports whether no other claim shares a square inch with c.
func alone(c claim, fabric []int) bool {
	for y := c.y; y < c.y+c.h; y++ {
		for x := c.x; x < c.x+c.w; x++ {
			if fabric[y*size+x] > 1 {
				return false
			}
		}
	}

	return true
}
