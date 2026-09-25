// Advent of Code 2019, day 3: Crossed Wires.
// https://adventofcode.com/2019/day/3
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 3, part1, part2)
}

// crossing is a point where the two wires meet. steps is the number of steps
// each wire takes to first reach it, added together.
type crossing struct {
	at    grid.Point
	steps int
}

// part1 returns the Manhattan distance from the central port to the closest
// crossing.
func part1(in string) any {
	best := -1

	for _, c := range crossings(in) {
		d := c.at.Manhattan(grid.P(0, 0))
		if best < 0 || d < best {
			best = d
		}
	}

	if best < 0 {
		return nil
	}

	return best
}

// part2 returns the fewest combined steps that the wires take to reach a
// crossing.
func part2(in string) any {
	best := -1

	for _, c := range crossings(in) {
		if best < 0 || c.steps < best {
			best = c.steps
		}
	}

	if best < 0 {
		return nil
	}

	return best
}

// crossings follows both wires from the central port and returns every point
// they share. The central port itself does not count.
func crossings(in string) []crossing {
	lines := input.Lines(in)
	first := trace(lines[0])

	var found []crossing

	// Walk the second wire again step by step instead of building a second
	// map. Each point it reaches that the first wire also reached is a
	// crossing. A wire can pass a crossing more than once, so seen keeps only
	// the first visit, which has the fewest steps.
	seen := map[grid.Point]bool{}
	walk(lines[1], func(p grid.Point, steps int) {
		firstSteps, ok := first[p]
		if !ok || seen[p] {
			return
		}

		seen[p] = true
		found = append(found, crossing{at: p, steps: firstSteps + steps})
	})

	return found
}

// trace returns every point a wire reaches, with the number of steps it took
// to reach that point the first time.
func trace(path string) map[grid.Point]int {
	steps := map[grid.Point]int{}

	walk(path, func(p grid.Point, n int) {
		if _, ok := steps[p]; !ok {
			steps[p] = n
		}
	})

	return steps
}

// walk follows a wire path such as "R8,U5" one grid step at a time. It calls
// visit with each point it reaches and the step count so far. The start point
// is not visited, because the wires do not cross at the central port.
//
// visit is a function value, so the callers pass in a closure. The closure
// can read and change the caller's local variables (such as a map), which
// lets trace and crossings share this walking code.
func walk(path string, visit func(p grid.Point, steps int)) {
	pos := grid.P(0, 0)
	steps := 0

	for _, move := range strings.Split(path, ",") {
		dir := grid.DirFromRune[rune(move[0])]

		for range input.Int(move[1:]) {
			pos = pos.Add(dir)
			steps++
			visit(pos, steps)
		}
	}
}
