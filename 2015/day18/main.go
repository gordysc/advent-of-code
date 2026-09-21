// Advent of Code 2015, day 18: Like a GIF For Your Yard.
// https://adventofcode.com/2015/day/18
package main

import (
	"aoc/lib/grid"
	"aoc/runner"
)

// animationSteps is the number of steps to animate in the puzzle. The worked
// example in the puzzle text stops after 4 steps (part 1) and 5 steps (part 2)
// instead, so the example answers from this program differ from the ones in
// the text.
const animationSteps = 100

// on and off are the characters for the two states of a light.
const (
	on  = '#'
	off = '.'
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 18, part1, part2)
}

// part1 counts the lights that are on after 100 steps of the animation.
func part1(in string) any {
	return countOn(animate(grid.Parse(in), animationSteps, false))
}

// part2 runs the same animation, but the four corner lights are stuck on.
func part2(in string) any {
	return countOn(animate(grid.Parse(in), animationSteps, true))
}

// animate runs the given number of steps and returns the final grid. When
// stuckCorners is true, the four corner lights are on at all times, which
// includes the initial state.
func animate(lights grid.Grid[byte], steps int, stuckCorners bool) grid.Grid[byte] {
	if stuckCorners {
		lightCorners(lights)
	}

	// All lights change at the same time, so a step must read from the old
	// state and write to a different grid. Two grids that swap roles after
	// each step do this without a new allocation for every step.
	next := grid.New[byte](lights.W, lights.H)

	for range steps {
		for p, state := range lights.All() {
			next.Set(p, nextState(state, onNeighbors(lights, p)))
		}

		if stuckCorners {
			lightCorners(next)
		}

		lights, next = next, lights
	}

	return lights
}

// nextState applies the rules of the animation to one light. A light that is
// on stays on with 2 or 3 neighbors that are on. A light that is off turns on
// with exactly 3 neighbors that are on.
func nextState(state byte, neighborsOn int) byte {
	if neighborsOn == 3 || (state == on && neighborsOn == 2) {
		return on
	}

	return off
}

// onNeighbors counts the lights around p that are on. Neighbors8 skips the
// points outside the grid, so the missing neighbors of an edge count as off.
func onNeighbors(lights grid.Grid[byte], p grid.Point) int {
	count := 0

	for _, q := range lights.Neighbors8(p) {
		if lights.At(q) == on {
			count++
		}
	}

	return count
}

// lightCorners turns on the four corner lights. Grid is passed by value, but
// its Cells field is a slice, and a copy of a slice points at the same memory.
// Thus Set changes the grid of the caller too.
func lightCorners(lights grid.Grid[byte]) {
	for _, x := range []int{0, lights.W - 1} {
		for _, y := range []int{0, lights.H - 1} {
			lights.Set(grid.P(x, y), on)
		}
	}
}

// countOn returns the number of lights that are on.
func countOn(lights grid.Grid[byte]) int {
	return lights.Count(func(state byte) bool { return state == on })
}
