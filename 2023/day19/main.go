// Advent of Code 2023, day 19: Aplenty.
// https://adventofcode.com/2023/day/19
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 19, part1, part2)
}

// rule is one step of a workflow. When cat is -1, the rule has no condition
// and always sends the part to target. Otherwise the rule sends the part to
// target when rating cat is less than (op '<') or more than (op '>') value.
type rule struct {
	cat    int
	op     byte
	value  int
	target string
}

// part holds the four ratings x, m, a and s, in that order.
type part [4]int

// span is a range of ratings from lo to hi. Both ends are included.
type span struct {
	lo, hi int
}

// box is a set of parts: one span for each of the four categories.
type box [4]span

// part1 adds all the ratings of the parts that the workflows accept.
func part1(in string) any {
	flows, parts := parse(in)
	total := 0

	for _, p := range parts {
		if accepts(flows, p) {
			total += p[0] + p[1] + p[2] + p[3]
		}
	}

	return total
}

// part2 counts the rating combinations from 1 to 4000 that the workflows
// accept.
//
// We do not test each of the 4000^4 parts. Instead, we send a 4D box of
// ratings through the workflows. Each rule cuts the box into the part that
// matches the rule and the part that does not. The two parts continue on
// different paths. The boxes that get to "A" never overlap, so we add their
// sizes.
func part2(in string) any {
	flows, _ := parse(in)
	full := box{{1, 4000}, {1, 4000}, {1, 4000}, {1, 4000}}

	return countAccepted(flows, "in", full)
}

// accepts follows one part through the workflows, from "in" to "A" or "R".
func accepts(flows map[string][]rule, p part) bool {
	name := "in"

	for name != "A" && name != "R" {
		for _, r := range flows[name] {
			if r.cat < 0 || matches(r, p[r.cat]) {
				name = r.target
				break
			}
		}
	}

	return name == "A"
}

// matches tells if a rating satisfies the condition of a rule.
func matches(r rule, v int) bool {
	if r.op == '<' {
		return v < r.value
	}

	return v > r.value
}

// countAccepted counts the parts in b that the workflow name (and the
// workflows after it) accept.
func countAccepted(flows map[string][]rule, name string, b box) int {
	if name == "R" {
		return 0
	}

	if name == "A" {
		return b.size()
	}

	total := 0

	for _, r := range flows[name] {
		if r.cat < 0 {
			total += countAccepted(flows, r.target, b)
			break
		}

		// Cut the box on this rule. The matching part goes to the target.
		// The rest goes on to the next rule in this workflow.
		match, rest := b, b
		s := b[r.cat]

		if r.op == '<' {
			match[r.cat] = span{s.lo, min(s.hi, r.value-1)}
			rest[r.cat] = span{max(s.lo, r.value), s.hi}
		} else {
			match[r.cat] = span{max(s.lo, r.value+1), s.hi}
			rest[r.cat] = span{s.lo, min(s.hi, r.value)}
		}

		if match[r.cat].lo <= match[r.cat].hi {
			total += countAccepted(flows, r.target, match)
		}

		if rest[r.cat].lo > rest[r.cat].hi {
			break
		}

		b = rest
	}

	return total
}

// size is the number of parts in the box.
func (b box) size() int {
	n := 1

	for _, s := range b {
		n *= s.hi - s.lo + 1
	}

	return n
}

// parse reads the workflows and the parts. A workflow line looks like
// "px{a<2006:qkq,m>2090:A,rfg}". A part line looks like
// "{x=787,m=2655,a=1222,s=2876}".
func parse(in string) (map[string][]rule, []part) {
	blocks := input.Blocks(in)
	flows := map[string][]rule{}

	for _, line := range blocks[0] {
		name, body, _ := strings.Cut(line, "{")
		body = strings.TrimSuffix(body, "}")

		var rules []rule

		for _, text := range strings.Split(body, ",") {
			cond, target, found := strings.Cut(text, ":")
			if !found {
				rules = append(rules, rule{cat: -1, target: text})
				continue
			}

			rules = append(rules, rule{
				cat:    strings.IndexByte("xmas", cond[0]),
				op:     cond[1],
				value:  input.Int(cond[2:]),
				target: target,
			})
		}

		flows[name] = rules
	}

	var parts []part

	for _, line := range blocks[1] {
		nums := input.Ints(line)
		parts = append(parts, part{nums[0], nums[1], nums[2], nums[3]})
	}

	return flows, parts
}
