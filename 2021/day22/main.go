// Advent of Code 2021, day 22: Reactor Reboot.
// https://adventofcode.com/2021/day/22
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 22, part1, part2)
}

// cuboid is a box of cubes. Both ends of each range are included.
type cuboid struct {
	x1, x2, y1, y2, z1, z2 int
}

// step is one reboot step: it turns a cuboid on or off.
type step struct {
	on bool
	c  cuboid
}

// part1 counts the cubes that are on in the region -50..50 on all axes.
//
// It clips each step to the region and ignores steps that are fully outside
// it. Then it counts the cubes the same way as part 2.
func part1(in string) any {
	region := cuboid{-50, 50, -50, 50, -50, 50}

	var steps []step
	for _, s := range parse(in) {
		if c, ok := s.c.intersect(region); ok {
			steps = append(steps, step{s.on, c})
		}
	}

	return countOn(steps)
}

// part2 counts the cubes that are on after all the steps.
//
// The cuboids are too large to store each cube. countOn keeps a list of
// signed cuboids instead, and adds their volumes at the end.
func part2(in string) any {
	return countOn(parse(in))
}

// signed is a cuboid with a sign. Its volume is added (+1) or taken away
// (-1) from the total.
type signed struct {
	c    cuboid
	sign int
}

// countOn applies the steps and returns the number of cubes that are on.
//
// For each step, it adds the overlap with each cuboid in the list, with the
// opposite sign. This cancels the cubes of the step that are already counted.
// Then, if the step turns cubes on, it adds the step's cuboid with sign +1.
func countOn(steps []step) int {
	var list []signed

	for _, s := range steps {
		var added []signed
		for _, prev := range list {
			if o, ok := s.c.intersect(prev.c); ok {
				added = append(added, signed{o, -prev.sign})
			}
		}

		if s.on {
			added = append(added, signed{s.c, 1})
		}

		list = append(list, added...)
	}

	total := 0
	for _, sc := range list {
		total += sc.sign * sc.c.volume()
	}

	return total
}

// intersect returns the overlap of two cuboids, and false if they do not
// overlap.
func (a cuboid) intersect(b cuboid) (cuboid, bool) {
	o := cuboid{
		max(a.x1, b.x1), min(a.x2, b.x2),
		max(a.y1, b.y1), min(a.y2, b.y2),
		max(a.z1, b.z1), min(a.z2, b.z2),
	}

	ok := o.x1 <= o.x2 && o.y1 <= o.y2 && o.z1 <= o.z2

	return o, ok
}

// volume returns the number of cubes in the cuboid.
func (c cuboid) volume() int {
	return (c.x2 - c.x1 + 1) * (c.y2 - c.y1 + 1) * (c.z2 - c.z1 + 1)
}

// parse reads the reboot steps, one for each line.
func parse(in string) []step {
	var steps []step

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		c := cuboid{n[0], n[1], n[2], n[3], n[4], n[5]}

		steps = append(steps, step{strings.HasPrefix(line, "on"), c})
	}

	return steps
}
