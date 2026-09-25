// Advent of Code 2023, day 5: If You Give A Seed A Fertilizer.
// https://adventofcode.com/2023/day/5
package main

import (
	"cmp"
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 5, part1, part2)
}

// part1 gives the lowest location of the seeds in the list. Each seed is a
// range that holds only that one seed.
func part1(in string) any {
	seeds, maps := parse(in)

	var spans []span

	for _, s := range seeds {
		spans = append(spans, span{s, s + 1})
	}

	return lowestLocation(spans, maps)
}

// part2 gives the lowest location when the seed list holds pairs of a start
// and a length.
//
// The real ranges hold billions of seeds, so the part cannot map each seed.
// Instead it maps whole ranges. A rule moves a range as one block, so each
// map only has to cut a range where the rules start and stop. The number of
// ranges stays small, and the lowest location is the lowest start of the
// ranges at the end.
func part2(in string) any {
	seeds, maps := parse(in)

	var spans []span

	for i := 0; i+1 < len(seeds); i += 2 {
		spans = append(spans, span{seeds[i], seeds[i] + seeds[i+1]})
	}

	return lowestLocation(spans, maps)
}

// span is a half-open range of numbers: it holds start and excludes end.
type span struct {
	start, end int
}

// rule moves the numbers from src up to (but not including) src+length by
// adding shift to them.
type rule struct {
	src, length, shift int
}

// parse reads the seed numbers and the maps. Each map is a list of rules,
// sorted by src so that a range can walk through them from low to high.
func parse(in string) ([]int, [][]rule) {
	blocks := input.Blocks(in)
	seeds := input.UInts(blocks[0][0])

	var maps [][]rule

	for _, block := range blocks[1:] {
		var rules []rule

		// The first line is the name of the map, such as "seed-to-soil map:".
		for _, line := range block[1:] {
			n := input.UInts(line)
			dst, src, length := n[0], n[1], n[2]
			rules = append(rules, rule{src, length, dst - src})
		}

		slices.SortFunc(rules, func(a, b rule) int { return cmp.Compare(a.src, b.src) })
		maps = append(maps, rules)
	}

	return seeds, maps
}

// lowestLocation sends the ranges through all maps in order and gives the
// lowest number that comes out at the end.
func lowestLocation(spans []span, maps [][]rule) int {
	for _, rules := range maps {
		var next []span

		for _, s := range spans {
			next = append(next, applyRules(s, rules)...)
		}

		spans = next
	}

	lowest := spans[0].start

	for _, s := range spans[1:] {
		lowest = min(lowest, s.start)
	}

	return lowest
}

// applyRules maps one range through one map and gives the resulting ranges.
//
// The rules are sorted and do not overlap. The function walks a cursor from
// the start of the range to its end. A gap before a rule keeps its numbers,
// and the part of the range inside a rule moves by the rule's shift.
func applyRules(s span, rules []rule) []span {
	var out []span
	cur := s.start

	for _, r := range rules {
		if cur >= s.end {
			break
		}

		ruleEnd := r.src + r.length
		if ruleEnd <= cur {
			continue
		}

		// The numbers before the rule starts are not in any rule, so they
		// keep their values.
		if r.src > cur {
			gapEnd := min(r.src, s.end)
			out = append(out, span{cur, gapEnd})
			cur = gapEnd
		}

		// The numbers that are inside the rule move by its shift.
		if cur < s.end && cur < ruleEnd {
			overlapEnd := min(ruleEnd, s.end)
			out = append(out, span{cur + r.shift, overlapEnd + r.shift})
			cur = overlapEnd
		}
	}

	// The numbers after the last rule also keep their values.
	if cur < s.end {
		out = append(out, span{cur, s.end})
	}

	return out
}
