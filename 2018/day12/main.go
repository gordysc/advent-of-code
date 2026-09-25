// Advent of Code 2018, day 12: Subterranean Sustainability.
// https://adventofcode.com/2018/day/12
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 12, part1, part2)
}

// shortRun and longRun are the generation counts for the two parts.
const (
	shortRun = 20
	longRun  = 50_000_000_000
)

// rules tells, for each pattern of five pots, whether the middle pot has a
// plant in the next generation. The index is the pattern read as a 5-bit
// number, with the leftmost pot as the high bit and a plant as a 1.
type rules [32]bool

// row is the run of pots from the leftmost plant to the rightmost plant.
// first is the pot number of pots[0].
type row struct {
	pots  string
	first int
}

// part1 sums the numbers of the pots that have a plant after 20 generations.
func part1(in string) any {
	r, rs := parse(in)

	for range shortRun {
		r = r.next(rs)
	}

	return r.sum()
}

// part2 sums the plant pot numbers after fifty billion generations. That is
// too many to simulate. After some time the pattern stops changing and only
// moves along the row by the same distance each generation. From then on the
// sum grows by the same amount each generation, so we can extrapolate.
// The example gives 999999999374.
func part2(in string) any {
	r, rs := parse(in)

	for gen := 1; gen <= longRun; gen++ {
		next := r.next(rs)

		if next.pots == r.pots {
			shift := next.first - r.first
			plants := strings.Count(next.pots, "#")

			return next.sum() + (longRun-gen)*shift*plants
		}

		r = next
	}

	return r.sum()
}

// parse reads the initial state and the rules. A rule that is not in the
// input leaves the middle pot empty.
func parse(in string) (row, rules) {
	blocks := input.Blocks(in)
	start := row{pots: strings.TrimPrefix(blocks[0][0], "initial state: ")}

	var rs rules
	for _, line := range blocks[1] {
		pattern, result, _ := strings.Cut(line, " => ")
		rs[mask(pattern)] = result == "#"
	}

	return start.trim(), rs
}

// mask reads a pattern of pots as a number, with a plant as a 1 bit.
func mask(pattern string) int {
	m := 0

	for i := range len(pattern) {
		m <<= 1
		if pattern[i] == '#' {
			m |= 1
		}
	}

	return m
}

// next returns the row one generation later. A new plant can grow at most
// two pots past the current ends, so the row is padded with four empty pots
// on each side. Then every pot that can change sits in the middle of a full
// window of five. This assumes the rule for five empty pots gives no plant,
// which is true for every real input.
func (r row) next(rs rules) row {
	padded := "...." + r.pots + "...."
	out := make([]byte, len(padded)-4)

	window := 0
	for i := range len(padded) {
		window = window<<1&0b11110 | bit(padded[i])

		if i >= 4 {
			out[i-4] = '.'
			if rs[window] {
				out[i-4] = '#'
			}
		}
	}

	// out[0] is the middle of the first window, which is padded[2], and that
	// is two pots left of the old first pot.
	return row{pots: string(out), first: r.first - 2}.trim()
}

// bit returns 1 for a plant and 0 for an empty pot.
func bit(c byte) int {
	if c == '#' {
		return 1
	}

	return 0
}

// trim removes the empty pots at both ends and moves first to match.
func (r row) trim() row {
	left := strings.IndexByte(r.pots, '#')
	if left < 0 {
		return row{}
	}

	right := strings.LastIndexByte(r.pots, '#')

	return row{pots: r.pots[left : right+1], first: r.first + left}
}

// sum adds up the pot numbers of the pots that have a plant.
func (r row) sum() int {
	total := 0

	for i := range len(r.pots) {
		if r.pots[i] == '#' {
			total += r.first + i
		}
	}

	return total
}
