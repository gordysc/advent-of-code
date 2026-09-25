// Advent of Code 2018, day 7: The Sum of Its Parts.
// https://adventofcode.com/2018/day/7
package main

import (
	"maps"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 7, part1, part2)
}

// workers is the number of elves that work on steps at the same time, and
// baseTime is the number of seconds every step takes before its letter is
// added. The worked example in the puzzle text uses 2 workers and a base time
// of 0 instead, so part 2 gives a different answer for the example.
const (
	workers  = 5
	baseTime = 60
)

// graph holds every step and the steps that each one waits for.
type graph struct {
	steps []byte          // every step, in alphabetical order
	needs map[byte][]byte // steps that must finish before each step can begin
}

// part1 does the steps one at a time. When more than one step is ready, the
// first one in alphabetical order goes first.
func part1(in string) any {
	g := parse(in)
	done := map[byte]bool{}
	var order []byte

	for len(order) < len(g.steps) {
		step := g.ready(done, done)[0]
		done[step] = true
		order = append(order, step)
	}

	return string(order)
}

// part2 finds how many seconds the workers need to do every step. Each free
// worker takes the first ready step in alphabetical order. Time then jumps to
// the moment the next step is finished, because nothing changes before that.
func part2(in string) any {
	g := parse(in)
	started := map[byte]bool{}
	done := map[byte]bool{}
	finish := map[byte]int{} // steps in progress and the second each one ends
	now := 0

	for len(done) < len(g.steps) {
		for _, step := range g.ready(done, started) {
			if len(finish) == workers {
				break
			}

			started[step] = true
			finish[step] = now + baseTime + int(step-'A') + 1
		}

		// maps.Values gives an iterator, so slices.Collect turns it into a
		// slice that slices.Min can read.
		now = slices.Min(slices.Collect(maps.Values(finish)))

		for step, end := range finish {
			if end == now {
				done[step] = true
				delete(finish, step)
			}
		}
	}

	return now
}

// parse reads the rules. Each rule says that one step must finish before
// another step can begin.
func parse(in string) graph {
	needs := map[byte][]byte{}

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)
		before, after := f[1][0], f[7][0]

		// Every step needs a key, also a step that waits for nothing, so that
		// the keys of needs are the full list of steps.
		needs[after] = append(needs[after], before)
		if _, ok := needs[before]; !ok {
			needs[before] = nil
		}
	}

	return graph{steps: slices.Sorted(maps.Keys(needs)), needs: needs}
}

// ready returns, in alphabetical order, the steps that have not started and
// whose prerequisites are all done.
func (g graph) ready(done, started map[byte]bool) []byte {
	var out []byte

	for _, step := range g.steps {
		if started[step] {
			continue
		}

		if !slices.ContainsFunc(g.needs[step], func(b byte) bool { return !done[b] }) {
			out = append(out, step)
		}
	}

	return out
}
