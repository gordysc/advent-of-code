// Advent of Code 2016, day 7: Internet Protocol Version 7.
// https://adventofcode.com/2016/day/7
package main

import (
	"iter"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 7, part1, part2)
}

// address is an IPv7 address split into the text outside square brackets
// (supernet sequences) and the text inside them (hypernet sequences).
type address struct {
	supernets []string
	hypernets []string
}

// part1 counts the addresses that support TLS.
func part1(in string) any {
	return slicesx.Count(parse(in), address.supportsTLS)
}

// part2 counts the addresses that support SSL.
func part2(in string) any {
	return slicesx.Count(parse(in), address.supportsSSL)
}

// parse splits each line into its supernet and hypernet sequences.
func parse(in string) []address {
	var out []address

	for _, line := range input.Lines(in) {
		parts := strings.Split(line, "[")
		addr := address{supernets: []string{parts[0]}}

		for _, part := range parts[1:] {
			hyper, super, _ := strings.Cut(part, "]")

			addr.hypernets = append(addr.hypernets, hyper)
			addr.supernets = append(addr.supernets, super)
		}

		out = append(out, addr)
	}

	return out
}

// supportsTLS reports whether an ABBA appears outside the brackets and none
// appears inside them.
func (a address) supportsTLS() bool {
	return slices.ContainsFunc(a.supernets, hasABBA) && !slices.ContainsFunc(a.hypernets, hasABBA)
}

// hasABBA reports whether s holds four characters in the form xyyx, where x
// and y are different.
func hasABBA(s string) bool {
	for i := 0; i+4 <= len(s); i++ {
		if s[i] != s[i+1] && s[i] == s[i+3] && s[i+1] == s[i+2] {
			return true
		}
	}

	return false
}

// supportsSSL reports whether an ABA outside the brackets has its matching
// BAB inside them.
func (a address) supportsSSL() bool {
	abas := set.New[[2]byte]()

	for _, s := range a.supernets {
		for x, y := range triples(s) {
			abas.Add([2]byte{x, y})
		}
	}

	for _, s := range a.hypernets {
		for x, y := range triples(s) {
			if abas.Has([2]byte{y, x}) {
				return true
			}
		}
	}

	return false
}

// triples yields x and y for every three characters in the form xyx, where x
// and y are different.
func triples(s string) iter.Seq2[byte, byte] {
	return func(yield func(x, y byte) bool) {
		for i := 0; i+3 <= len(s); i++ {
			if s[i] != s[i+1] && s[i] == s[i+2] && !yield(s[i], s[i+1]) {
				return
			}
		}
	}
}
