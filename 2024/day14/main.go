// Advent of Code 2024, day 14: Restroom Redoubt.
// https://adventofcode.com/2024/day/14
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 14, part1, part2)
}

// part1 moves the robots for 100 seconds and returns the safety factor: the
// product of the robot counts in the four quadrants.
func part1(in string) any {
	robots := parse(in)
	w, h := size(robots)

	var quad [4]int
	for _, r := range robots {
		x := mathx.Mod(r.px+100*r.vx, w)
		y := mathx.Mod(r.py+100*r.vy, h)

		// Robots on the middle row or the middle column are in no quadrant.
		if x == w/2 || y == h/2 {
			continue
		}

		i := 0
		if x > w/2 {
			i++
		}

		if y > h/2 {
			i += 2
		}

		quad[i]++
	}

	return quad[0] * quad[1] * quad[2] * quad[3]
}

// part2 finds the first second where the robots draw a Christmas tree.
//
// The puzzle does not say what the tree looks like, and there is no example
// answer. So this part returns nil for the 11x7 example grid.
//
// In the tree picture, most robots are in a small box, so the spread of their
// positions is much smaller than at other times. The X positions repeat every
// w seconds and the Y positions repeat every h seconds, and each axis moves
// independently of the other. So we find the second in [0, w) with the lowest
// X variance, and the second in [0, h) with the lowest Y variance. Then the
// Chinese remainder theorem gives the one second in [0, w*h) that matches
// both. This checks w + h seconds, not all w*h seconds.
//
// This assumes w and h are coprime, which is true for 101 and 103.
func part2(in string) any {
	robots := parse(in)
	w, h := size(robots)
	if w != realW {
		return nil
	}

	tx := calmest(robots, w, func(r robot) (int, int) { return r.px, r.vx })
	ty := calmest(robots, h, func(r robot) (int, int) { return r.py, r.vy })

	// Step through the seconds that match tx until one also matches ty.
	for t := tx; t < w*h; t += w {
		if t%h == ty {
			return t
		}
	}

	return nil
}

// The grid sizes. Every robot of the example fits in the small grid. The
// real inputs use the large grid.
const (
	exampleW, exampleH = 11, 7
	realW, realH       = 101, 103
)

// robot holds the start position and the velocity of one robot.
type robot struct {
	px, py int
	vx, vy int
}

// parse reads one robot per line. input.Ints keeps the minus signs of the
// negative velocities.
func parse(in string) []robot {
	var out []robot

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		if len(n) != 4 {
			continue
		}

		out = append(out, robot{n[0], n[1], n[2], n[3]})
	}

	return out
}

// size picks the grid size. The input does not give it, so we use the small
// example grid when every robot starts inside it, and the real grid if not.
func size(robots []robot) (int, int) {
	for _, r := range robots {
		if r.px >= exampleW || r.py >= exampleH {
			return realW, realH
		}
	}

	return exampleW, exampleH
}

// calmest returns the second in [0, n) where the robots have the lowest
// variance along one axis. axis gives the position and the velocity of a
// robot on that axis, and n is the length of the grid on that axis.
func calmest(robots []robot, n int, axis func(robot) (int, int)) int {
	best, bestVar := 0, -1

	for t := range n {
		sum, sumSq := 0, 0

		for _, r := range robots {
			p, v := axis(r)
			x := mathx.Mod(p+t*v, n)
			sum += x
			sumSq += x * x
		}

		// This is the variance times len(robots) squared. The factor is the
		// same for every t, so it does not change which t is the lowest, and
		// it keeps the math in integers.
		variance := len(robots)*sumSq - sum*sum
		if bestVar < 0 || variance < bestVar {
			best, bestVar = t, variance
		}
	}

	return best
}
