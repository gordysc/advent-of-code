// Advent of Code 2023, day 15: Lens Library.
// https://adventofcode.com/2023/day/15
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 15, part1, part2)
}

// part1 adds the HASH values of all steps in the initialization sequence.
func part1(in string) any {
	total := 0

	for _, step := range steps(in) {
		total += hash(step)
	}

	return total
}

// lens is one lens in a box.
type lens struct {
	label string
	focal int
}

// part2 does the steps in the 256 boxes and gives the focusing power of all
// lenses at the end.
//
// The HASH of the label selects the box. The step "label-" removes the lens
// with that label from its box. The step "label=N" replaces the focal length
// of the lens with that label, or adds a new lens at the back of the box.
func part2(in string) any {
	var boxes [256][]lens

	for _, step := range steps(in) {
		if label, ok := strings.CutSuffix(step, "-"); ok {
			box := &boxes[hash(label)]

			// slices.DeleteFunc removes the matching lenses and keeps the
			// order of the others. It returns the shorter slice.
			*box = slices.DeleteFunc(*box, func(l lens) bool { return l.label == label })

			continue
		}

		label, value, _ := strings.Cut(step, "=")
		focal := input.Int(value)
		box := &boxes[hash(label)]

		i := slices.IndexFunc(*box, func(l lens) bool { return l.label == label })
		if i >= 0 {
			(*box)[i].focal = focal
		} else {
			*box = append(*box, lens{label, focal})
		}
	}

	power := 0

	for b, box := range boxes {
		for slot, l := range box {
			power += (b + 1) * (slot + 1) * l.focal
		}
	}

	return power
}

// steps splits the initialization sequence at the commas. Newlines in the
// input are ignored.
func steps(in string) []string {
	return strings.Split(strings.ReplaceAll(in, "\n", ""), ",")
}

// hash is the HASH algorithm: for each character, add its ASCII code,
// multiply by 17 and keep the remainder after division by 256.
func hash(s string) int {
	h := 0

	for i := range len(s) {
		h = (h + int(s[i])) * 17 % 256
	}

	return h
}
