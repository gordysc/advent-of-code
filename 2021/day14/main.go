// Advent of Code 2021, day 14: Extended Polymerization.
// https://adventofcode.com/2021/day/14
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 14, part1, part2)
}

// part1 runs 10 steps of pair insertion. It returns the count of the most
// common element minus the count of the least common element.
func part1(in string) any {
	return grow(in, 10)
}

// part2 does the same as part1, but for 40 steps.
func part2(in string) any {
	return grow(in, 40)
}

// pair is two elements next to each other in the polymer.
type pair [2]byte

// grow runs the given number of insertion steps and returns the difference
// between the most and least common element counts.
//
// The polymer doubles in length at each step, so it is too long to build.
// But each insertion only changes one pair into two pairs: AB with rule
// AB -> C becomes AC and CB. The order of the pairs does not matter, so it is
// enough to count how many times each pair occurs.
//
// At the end, each element is the first element of exactly one pair, except
// for the last element of the polymer. Insertion never changes the last
// element, so it comes from the template.
func grow(in string, steps int) int {
	blocks := input.Blocks(in)
	template := blocks[0][0]

	rules := map[pair]byte{}
	for _, line := range blocks[1] {
		from, to, _ := strings.Cut(line, " -> ")
		rules[pair{from[0], from[1]}] = to[0]
	}

	pairs := map[pair]int{}
	for i := range len(template) - 1 {
		pairs[pair{template[i], template[i+1]}]++
	}

	for range steps {
		next := make(map[pair]int, len(pairs))

		for p, n := range pairs {
			c, ok := rules[p]
			if !ok {
				next[p] += n
				continue
			}

			next[pair{p[0], c}] += n
			next[pair{c, p[1]}] += n
		}

		pairs = next
	}

	elements := map[byte]int{template[len(template)-1]: 1}
	for p, n := range pairs {
		elements[p[0]] += n
	}

	counts := make([]int, 0, len(elements))
	for _, n := range elements {
		counts = append(counts, n)
	}

	return slices.Max(counts) - slices.Min(counts)
}
