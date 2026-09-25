// Advent of Code 2025, day 5: Cafeteria.
// https://adventofcode.com/2025/day/5
package main

import (
	"cmp"
	"slices"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 5, part1, part2)
}

// part1 counts the available ingredient IDs that are in a fresh range.
//
// The ranges are merged first, so they are sorted and do not overlap. Then a
// binary search finds the one range that can hold each ID.
func part1(in string) any {
	blocks := input.Blocks(in)
	if len(blocks) < 2 {
		return nil
	}

	ranges := merge(parseRanges(blocks[0]))
	fresh := 0

	for _, line := range blocks[1] {
		ids := input.UInts(line)
		if len(ids) == 0 {
			continue
		}

		if contains(ranges, ids[0]) {
			fresh++
		}
	}

	return fresh
}

// part2 counts all the IDs that the fresh ranges cover. The list of available
// IDs is not used.
//
// After the merge, no two ranges overlap, so the sum of their lengths counts
// each ID one time only.
func part2(in string) any {
	blocks := input.Blocks(in)
	total := 0

	for _, r := range merge(parseRanges(blocks[0])) {
		total += r.hi - r.lo + 1
	}

	return total
}

// span is an inclusive range of ingredient IDs.
type span struct {
	lo, hi int
}

// parseRanges reads lines such as "3-5". It uses UInts because Ints would read
// the dash as a minus sign.
func parseRanges(lines []string) []span {
	var ranges []span

	for _, line := range lines {
		n := input.UInts(line)
		if len(n) < 2 {
			continue
		}

		ranges = append(ranges, span{n[0], n[1]})
	}

	return ranges
}

// merge sorts the ranges by start and joins each range that overlaps or
// touches the one before it. The result is sorted and has no overlaps.
func merge(ranges []span) []span {
	// slices.SortFunc sorts in place. cmp.Compare returns -1, 0 or +1, which is
	// the result that SortFunc expects.
	slices.SortFunc(ranges, func(a, b span) int {
		return cmp.Compare(a.lo, b.lo)
	})

	var merged []span

	for _, r := range ranges {
		last := len(merged) - 1

		// The ranges are inclusive, so 3-5 and 6-8 touch and become 3-8.
		if last >= 0 && r.lo <= merged[last].hi+1 {
			merged[last].hi = max(merged[last].hi, r.hi)
			continue
		}

		merged = append(merged, r)
	}

	return merged
}

// contains reports whether id is in one of the merged ranges.
func contains(ranges []span, id int) bool {
	// slices.BinarySearchFunc returns the index of the first range whose end is
	// not below id. Only that range can hold id, because the ranges are sorted
	// and do not overlap.
	i, _ := slices.BinarySearchFunc(ranges, id, func(r span, id int) int {
		return cmp.Compare(r.hi, id)
	})

	return i < len(ranges) && ranges[i].lo <= id
}
