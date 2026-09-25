// Advent of Code 2025, day 11: Reactor.
// https://adventofcode.com/2025/day/11
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 11, part1, part2)
}

// part1 returns the number of paths from the device "you" to the device
// "out".
//
// The puzzle has two examples. The file example.txt holds the part 1 example.
// It has no device "svr", so part 2 returns nil on it. On the part 2 example,
// part 1 returns nil, because that example has no device "you".
func part1(in string) any {
	g := parse(in)
	if _, ok := g["you"]; !ok {
		return nil
	}

	return g.countPaths("you", "out")
}

// part2 returns the number of paths from the device "svr" to the device "out"
// that go through both "dac" and "fft". On the part 2 example, it returns 2.
//
// The graph has no cycles, because otherwise the number of paths would have no
// limit. So a path cannot visit "fft" after "dac" and also "dac" after "fft".
// Each path goes through them in one order, and we add the paths for the two
// orders. The paths for one order are the product of the counts for its three
// parts, because any path of each part can join to any path of the next part.
func part2(in string) any {
	g := parse(in)
	if _, ok := g["svr"]; !ok {
		return nil
	}

	dacFirst := g.countPaths("svr", "dac") * g.countPaths("dac", "fft") * g.countPaths("fft", "out")
	fftFirst := g.countPaths("svr", "fft") * g.countPaths("fft", "dac") * g.countPaths("dac", "out")

	return dacFirst + fftFirst
}

// graph maps each device to the devices that its outputs go to.
type graph map[string][]string

// parse reads lines of the form "aaa: bbb ccc".
func parse(in string) graph {
	g := graph{}

	for _, line := range input.Lines(in) {
		// strings.Cut splits at the first ":" and returns false if there is none.
		name, outputs, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		g[strings.TrimSpace(name)] = strings.Fields(outputs)
	}

	return g
}

// countPaths returns the number of paths from one device to another.
//
// The number of paths from a device is the sum of the numbers of paths from
// its outputs. Many paths share the same devices, so we keep each result in a
// memo and find it only once. Without the memo, the time would grow with the
// number of paths, which is very large on real inputs.
func (g graph) countPaths(from, to string) int {
	memo := map[string]int{}

	var count func(device string) int
	count = func(device string) int {
		if device == to {
			return 1
		}

		if n, ok := memo[device]; ok {
			return n
		}

		total := 0
		for _, next := range g[device] {
			total += count(next)
		}

		memo[device] = total

		return total
	}

	return count(from)
}
