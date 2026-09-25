// Advent of Code 2017, day 24: Electromagnetic Moat.
// https://adventofcode.com/2017/day/24
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 24, part1, part2)
}

// component is a piece of bridge with a port at each end. The number on a
// port is how many pins it has.
type component struct {
	a, b int
}

// best holds the two answers that one search finds: the strongest bridge,
// and the length and strength of the longest bridge.
type best struct {
	strongest       int
	longest         int
	longestStrength int
}

// part1 finds the strength of the strongest bridge.
func part1(in string) any {
	return search(parse(in)).strongest
}

// part2 finds the strength of the longest bridge. When two bridges have the
// same length, the stronger one wins.
func part2(in string) any {
	return search(parse(in)).longestStrength
}

// parse reads one "a/b" component from every line.
func parse(in string) []component {
	var parts []component

	for _, line := range input.Lines(in) {
		n := input.UInts(line)
		parts = append(parts, component{n[0], n[1]})
	}

	return parts
}

// search tries every bridge that starts at the 0-pin port and returns the
// best ones for both parts.
//
// It is a depth-first search. A bit mask marks the components that the
// current bridge uses, so each one goes in at most once. The inputs have
// fewer than 64 components, so the mask fits in one uint64. At each step the
// search tries every unused component that has a port which matches the open
// end of the bridge, and the other port of that component becomes the new
// open end.
func search(parts []component) best {
	var res best

	// build extends a bridge whose open end has port pins, that has length
	// components, and that has the given strength.
	var build func(used uint64, port, length, strength int)
	build = func(used uint64, port, length, strength int) {
		res.strongest = max(res.strongest, strength)
		if length > res.longest || length == res.longest && strength > res.longestStrength {
			res.longest = length
			res.longestStrength = strength
		}

		for i, c := range parts {
			if used&(1<<i) != 0 {
				continue
			}

			next := strength + c.a + c.b
			switch port {
			case c.a:
				build(used|1<<i, c.b, length+1, next)
			case c.b:
				build(used|1<<i, c.a, length+1, next)
			}
		}
	}

	build(0, 0, 0, 0)

	return res
}
