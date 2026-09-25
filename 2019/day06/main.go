// Advent of Code 2019, day 6: Universal Orbit Map.
// https://adventofcode.com/2019/day/6
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 6, part1, part2)
}

// part1 counts the direct and indirect orbits. Each object orbits every
// object on its path down to COM, so the count for one object is its depth,
// and the answer is the sum of all depths.
//
// example.txt is the part 2 example, which adds YOU and SAN to the part 1 map.
// Those two extra objects add 7 and 5 orbits, so part 1 gives 54 for it and
// not the 42 from the part 1 example.
func part1(in string) any {
	parents := parse(in)
	depths := map[string]int{}

	total := 0
	for obj := range parents {
		total += depth(obj, parents, depths)
	}

	return total
}

// part2 counts the orbital transfers that move YOU to orbit the same object
// as SAN. The path goes down from YOU's parent to the closest object that
// both YOU and SAN orbit, then up to SAN's parent. The part 1 example has no
// YOU or SAN, so part 2 has no answer for it.
func part2(in string) any {
	parents := parse(in)

	if _, ok := parents["YOU"]; !ok {
		return nil
	}

	if _, ok := parents["SAN"]; !ok {
		return nil
	}

	// Record how many steps down from YOU's parent each ancestor is.
	steps := map[string]int{}
	for obj, n := parents["YOU"], 0; obj != ""; obj, n = parents[obj], n+1 {
		steps[obj] = n
	}

	// Walk down from SAN's parent until the path meets YOU's path. The first
	// shared object is the closest common one.
	for obj, n := parents["SAN"], 0; obj != ""; obj, n = parents[obj], n+1 {
		if m, ok := steps[obj]; ok {
			return m + n
		}
	}

	return nil
}

// parse reads the "A)B" lines into a map from each object to the object it
// orbits. COM orbits nothing, so it is not a key, and a lookup for it gives
// the empty string. The walks in part 2 use that to stop.
func parse(in string) map[string]string {
	parents := map[string]string{}

	for _, line := range input.Lines(in) {
		center, obj, _ := strings.Cut(line, ")")
		parents[obj] = center
	}

	return parents
}

// depth returns how many objects obj orbits, directly and indirectly. It
// saves each result in depths, so every object is counted only once even
// though many objects share the same path down to COM.
func depth(obj string, parents map[string]string, depths map[string]int) int {
	parent, ok := parents[obj]
	if !ok {
		return 0
	}

	if d, ok := depths[obj]; ok {
		return d
	}

	d := depth(parent, parents, depths) + 1
	depths[obj] = d

	return d
}
