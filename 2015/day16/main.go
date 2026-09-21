// Advent of Code 2015, day 16: Aunt Sue.
// https://adventofcode.com/2015/day/16
//
// The puzzle text has no worked example. The example.txt file here is a small
// handmade list: Sue 2 is the answer to part 1 and Sue 3 is the answer to part 2.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// tickerTape is the output of the My First Crime Scene Analysis Machine. It is
// part of the puzzle text and is the same for every input.
var tickerTape = map[string]int{
	"children":    3,
	"cats":        7,
	"samoyeds":    2,
	"pomeranians": 3,
	"akitas":      0,
	"vizslas":     0,
	"goldfish":    5,
	"trees":       3,
	"cars":        2,
	"perfumes":    1,
}

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 16, part1, part2)
}

// part1 finds the Sue whose known things all have exactly the counts on the
// ticker tape.
func part1(in string) any {
	return findSue(parseSues(in), exactMatch)
}

// part2 finds the Sue again after the correction to the instructions: the
// readings for some of the things are ranges, not exact counts.
func part2(in string) any {
	return findSue(parseSues(in), rangeMatch)
}

// sue holds what you remember about one Aunt Sue. A thing that is not in the
// map is a thing you do not remember, not a count of zero.
type sue struct {
	number int
	things map[string]int
}

// parseSues reads lines like
// "Sue 1: goldfish: 9, cars: 0, samoyeds: 9"
// into a slice of sues.
func parseSues(in string) []sue {
	var sues []sue

	for _, line := range input.Lines(in) {
		// Only the first colon divides the label from the list of things. The
		// list has colons of its own, so Cut is safer here than Split.
		label, list, _ := strings.Cut(line, ": ")
		s := sue{number: input.Int(strings.TrimPrefix(label, "Sue ")), things: map[string]int{}}

		for _, item := range strings.Split(list, ", ") {
			thing, count, _ := strings.Cut(item, ": ")
			s.things[thing] = input.Int(count)
		}

		sues = append(sues, s)
	}

	return sues
}

// findSue returns the number of the first Sue for which every remembered thing
// agrees with the ticker tape. The matches function decides what "agrees"
// means. The result is nil when no Sue fits, which the runner shows as not
// solved.
func findSue(sues []sue, matches func(thing string, count, reading int) bool) any {
	for _, s := range sues {
		if s.fits(matches) {
			return s.number
		}
	}

	return nil
}

// fits reports whether every thing you remember about this Sue agrees with the
// ticker tape. Things you do not remember cannot rule her out.
func (s sue) fits(matches func(thing string, count, reading int) bool) bool {
	for thing, count := range s.things {
		if !matches(thing, count, tickerTape[thing]) {
			return false
		}
	}

	return true
}

// exactMatch is the rule for part 1: the count must equal the reading.
func exactMatch(thing string, count, reading int) bool {
	return count == reading
}

// rangeMatch is the rule for part 2. The cats and trees readings mean "more
// than this many". The pomeranians and goldfish readings mean "fewer than this
// many". All other readings are still exact.
func rangeMatch(thing string, count, reading int) bool {
	switch thing {
	case "cats", "trees":
		return count > reading
	case "pomeranians", "goldfish":
		return count < reading
	default:
		return count == reading
	}
}
