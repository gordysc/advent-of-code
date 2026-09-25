// Advent of Code 2020, day 13: Shuttle Search.
// https://adventofcode.com/2020/day/13
package main

import (
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 13, part1, part2)
}

// part1 finds the first bus that departs at or after the earliest time. A bus
// with ID id departs at each multiple of id, so the wait for it is the
// distance from the earliest time up to the next multiple. The answer is the
// bus ID times the wait.
func part1(in string) any {
	earliest, buses := parse(in)

	bestID, bestWait := 0, -1
	for _, b := range buses {
		wait := mathx.Mod(-earliest, b.id)
		if bestWait < 0 || wait < bestWait {
			bestID, bestWait = b.id, wait
		}
	}

	return bestID * bestWait
}

// part2 finds the first time t where each bus departs at t plus its offset in
// the list. This is a system of congruences, and a sieve solves it. The
// function adds one bus at a time. It steps t by the LCM of the IDs that
// already match, so those buses stay matched, until the new bus also
// matches. The step then grows to include the new ID.
func part2(in string) any {
	_, buses := parse(in)

	t, step := 0, 1
	for _, b := range buses {
		for (t+b.offset)%b.id != 0 {
			t += step
		}

		step = mathx.LCM(step, b.id)
	}

	return t
}

// bus is one bus ID and its position in the schedule list.
type bus struct {
	id     int
	offset int
}

// parse reads the earliest departure time from the first line and the buses
// from the second line. It skips each "x", but it keeps the list position as
// the offset of each bus.
func parse(in string) (int, []bus) {
	lines := input.Lines(in)

	earliest, err := strconv.Atoi(lines[0])
	if err != nil {
		panic(err)
	}

	var buses []bus
	for i, field := range strings.Split(lines[1], ",") {
		if field == "x" {
			continue
		}

		id, err := strconv.Atoi(field)
		if err != nil {
			panic(err)
		}

		buses = append(buses, bus{id, i})
	}

	return earliest, buses
}
