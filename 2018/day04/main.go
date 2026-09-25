// Advent of Code 2018, day 4: Repose Record.
// https://adventofcode.com/2018/day/4
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 4, part1, part2)
}

// minutes counts, for each minute of the midnight hour, how many times a guard
// was asleep during it.
type minutes [60]int

// part1 finds the guard who sleeps the most minutes in total, and multiplies
// their ID by the minute they are most often asleep.
func part1(in string) any {
	guards := sleep(in)
	best, bestTotal := 0, -1

	for id, m := range guards {
		total := 0
		for _, n := range m {
			total += n
		}

		if total > bestTotal {
			best, bestTotal = id, total
		}
	}

	minute, _ := sleepiest(guards[best])

	return best * minute
}

// part2 finds the guard who is asleep most often on the same minute, and
// multiplies their ID by that minute.
func part2(in string) any {
	best, bestMinute, bestCount := 0, 0, -1

	for id, m := range sleep(in) {
		minute, count := sleepiest(m)
		if count > bestCount {
			best, bestMinute, bestCount = id, minute, count
		}
	}

	return best * bestMinute
}

// sleep reads the log and counts the minutes each guard is asleep. The lines
// in the input are not in order, but each starts with a timestamp in a fixed
// format, so a plain string sort puts them in time order.
func sleep(in string) map[int]*minutes {
	lines := input.Lines(in)
	slices.Sort(lines)

	guards := map[int]*minutes{}
	var current *minutes
	asleep := 0

	for _, line := range lines {
		minute := input.Int(line[15:17])

		switch {
		case strings.Contains(line, "Guard"):
			id := input.UInts(line[18:])[0]
			if guards[id] == nil {
				guards[id] = &minutes{}
			}

			current = guards[id]

		case strings.HasSuffix(line, "falls asleep"):
			asleep = minute

		case strings.HasSuffix(line, "wakes up"):
			for m := asleep; m < minute; m++ {
				current[m]++
			}
		}
	}

	return guards
}

// sleepiest returns the minute a guard is most often asleep, and how many
// times they were asleep during it.
func sleepiest(m *minutes) (int, int) {
	best := 0

	for i, n := range m {
		if n > m[best] {
			best = i
		}
	}

	return best, m[best]
}
