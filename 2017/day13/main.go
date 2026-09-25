// Advent of Code 2017, day 13: Packet Scanners.
// https://adventofcode.com/2017/day/13
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 13, part1, part2)
}

// layer is one firewall layer that has a scanner.
type layer struct {
	depth, rng int
}

// part1 adds up depth times range for every layer that catches the packet when
// it leaves at once.
func part1(in string) any {
	severity := 0

	for _, l := range parse(in) {
		if caught(l, 0) {
			severity += l.depth * l.rng
		}
	}

	return severity
}

// part2 finds the smallest delay that lets the packet pass without being
// caught. Each delay stops at the first layer that catches the packet, so
// most delays are rejected after only a few checks.
func part2(in string) any {
	layers := parse(in)

	for delay := 0; ; delay++ {
		safe := true
		for _, l := range layers {
			if caught(l, delay) {
				safe = false
				break
			}
		}

		if safe {
			return delay
		}
	}
}

// caught reports whether the scanner in l is at the top when a packet that
// waited delay picoseconds reaches it. A scanner with range r goes down and
// back up again, so it is at the top every 2*(r-1) picoseconds.
func caught(l layer, delay int) bool {
	return (l.depth+delay)%(2*(l.rng-1)) == 0
}

// parse reads each "depth: range" line into a layer.
func parse(in string) []layer {
	var layers []layer

	for _, line := range input.Lines(in) {
		nums := input.Ints(line)
		layers = append(layers, layer{depth: nums[0], rng: nums[1]})
	}

	return layers
}
