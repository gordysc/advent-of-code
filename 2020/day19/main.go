// Advent of Code 2020, day 19: Monster Messages.
// https://adventofcode.com/2020/day/19
package main

import (
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 19, part1, part2)
}

// part1 counts the messages that match rule 0 completely.
//
// example.txt is the part 2 example, because only that one has rules 8 and
// 11 in a form where the part 2 change matters. The puzzle gives 3 as its
// part 1 answer.
func part1(in string) any {
	rules, messages := parse(in)

	return countMatches(rules, messages)
}

// part2 replaces rules 8 and 11 with looping versions and counts the matches
// again. The matcher returns every possible end position, so it handles the
// loops without special code. Each loop step consumes at least one character
// before it recurses, so the recursion always ends.
func part2(in string) any {
	rules, messages := parse(in)

	rules[8] = rule{alts: [][]int{{42}, {42, 8}}}
	rules[11] = rule{alts: [][]int{{42, 31}, {42, 11, 31}}}

	return countMatches(rules, messages)
}

// countMatches counts the messages for which rule 0 can end exactly at the
// end of the message.
func countMatches(rules map[int]rule, messages []string) int {
	count := 0
	for _, msg := range messages {
		for _, end := range match(rules, 0, msg, 0) {
			if end == len(msg) {
				count++
				break
			}
		}
	}

	return count
}

// rule is either one literal character, or a list of alternatives. Each
// alternative is a sequence of rule numbers that must match one after the
// other.
type rule struct {
	char byte
	alts [][]int
}

// match tries rule id at position pos of msg and returns every position
// where a match can end. An empty result means no match.
//
// A literal matches one character. For an alternative, we start with the set
// {pos} and push each end position through the next rule in the sequence.
// The results of all alternatives are joined.
func match(rules map[int]rule, id int, msg string, pos int) []int {
	r := rules[id]

	if r.char != 0 {
		if pos < len(msg) && msg[pos] == r.char {
			return []int{pos + 1}
		}

		return nil
	}

	var ends []int
	for _, seq := range r.alts {
		cur := []int{pos}
		for _, sub := range seq {
			var next []int
			for _, p := range cur {
				next = append(next, match(rules, sub, msg, p)...)
			}

			cur = next
			if len(cur) == 0 {
				break
			}
		}

		ends = append(ends, cur...)
	}

	return ends
}

// parse reads the rules block and the messages block. A rule line is either
// `N: "a"` or `N: 1 2 | 3 4`.
func parse(in string) (map[int]rule, []string) {
	blocks := input.Blocks(in)
	rules := map[int]rule{}

	for _, line := range blocks[0] {
		idText, body, _ := strings.Cut(line, ": ")
		id, _ := strconv.Atoi(idText)

		if body[0] == '"' {
			rules[id] = rule{char: body[1]}
			continue
		}

		var r rule
		for _, alt := range strings.Split(body, " | ") {
			r.alts = append(r.alts, input.Ints(alt))
		}

		rules[id] = r
	}

	return rules, blocks[1]
}
