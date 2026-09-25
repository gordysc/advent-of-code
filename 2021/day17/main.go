// Advent of Code 2021, day 17: Trick Shot.
// https://adventofcode.com/2021/day/17
package main

import (
	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 17, part1, part2)
}

// part1 finds the highest position that a probe can reach and still hit the
// target area.
//
// A probe that goes up with speed vy comes back to y=0 with speed -(vy+1).
// Its highest point is the triangular number of vy. So the largest vy that
// hits the target gives the answer.
func part1(in string) any {
	best := 0

	for _, v := range hits(parse(in)) {
		// A probe that starts downwards never goes above y=0.
		best = max(best, mathx.Triangular(max(v.Y, 0)))
	}

	return best
}

// part2 counts the initial velocities that hit the target area.
//
// The search space is small. The x speed cannot be negative or more than the
// far edge of the target, otherwise the probe passes it on the first step.
// The y speed cannot be less than the bottom of the target, for the same
// reason. It cannot be more than the depth of the bottom minus one, because
// the probe then passes the target when it comes back down through y=0.
func part2(in string) any {
	return len(hits(parse(in)))
}

// target is the rectangular target area.
type target struct {
	minX, maxX, minY, maxY int
}

// parse reads the target area. The target is always to the right of and
// below the submarine.
func parse(in string) target {
	n := input.Ints(in)

	return target{minX: n[0], maxX: n[1], minY: n[2], maxY: n[3]}
}

// hits returns all the initial velocities that put the probe in the target
// area after some step. Velocities use grid.Point, with Y as the up speed.
func hits(t target) []grid.Point {
	var result []grid.Point

	for vx := 0; vx <= t.maxX; vx++ {
		for vy := t.minY; vy < -t.minY; vy++ {
			if hitsTarget(t, vx, vy) {
				result = append(result, grid.P(vx, vy))
			}
		}
	}

	return result
}

// hitsTarget simulates one launch and tells if the probe is ever inside the
// target area at the end of a step.
func hitsTarget(t target, vx, vy int) bool {
	x, y := 0, 0

	// Stop when the probe is past the far edge or below the bottom. After
	// that, it can never come back.
	for x <= t.maxX && y >= t.minY {
		if x >= t.minX && y <= t.maxY {
			return true
		}

		x += vx
		y += vy

		vx -= mathx.Sign(vx)
		vy--
	}

	return false
}
