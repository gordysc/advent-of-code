// Advent of Code 2015, day 19: Medicine for Rudolph.
// https://adventofcode.com/2015/day/19
//
// The puzzle text has two examples with different rule sets. The example.txt
// file here uses the rules from the part 2 example and the molecule HOH. The
// two e rules never match HOH, so part 1 still gives 4, the answer in the text.
package main

import (
	"math/rand/v2"
	"strings"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// electron is the molecule that the machine starts from in part 2.
const electron = "e"

// replacement is one rule of the machine: from becomes to.
type replacement struct {
	from, to string
}

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 19, part1, part2)
}

// part1 counts the distinct molecules that one replacement can make.
func part1(in string) any {
	rules, molecule := parse(in)

	return distinctMolecules(rules, molecule).Len()
}

// part2 finds the fewest steps to make the molecule from a single electron.
func part2(in string) any {
	rules, molecule := parse(in)

	return fewestSteps(rules, molecule)
}

// parse splits the input into the replacement rules and the target molecule.
func parse(in string) ([]replacement, string) {
	blocks := input.Blocks(in)
	rules := make([]replacement, 0, len(blocks[0]))

	for _, line := range blocks[0] {
		from, to, ok := strings.Cut(line, " => ")
		if !ok {
			panic("bad rule: " + line)
		}

		rules = append(rules, replacement{from, to})
	}

	return rules, blocks[1][0]
}

// distinctMolecules applies every rule at every position of the molecule and
// collects the results in a set, which removes the duplicates.
func distinctMolecules(rules []replacement, molecule string) set.Set[string] {
	seen := set.New[string]()

	for _, rule := range rules {
		for start := 0; ; start++ {
			i := strings.Index(molecule[start:], rule.from)
			if i < 0 {
				break
			}

			start += i
			seen.Add(molecule[:start] + rule.to + molecule[start+len(rule.from):])
		}
	}

	return seen
}

// fewestSteps runs the machine backwards. It replaces products with their
// sources until only the electron is left. Each backwards step undoes one
// forwards step, so the count of backwards steps is the answer.
//
// A greedy reduction can reach a dead end: a molecule that is not the electron
// but that no rule can shrink. The puzzle grammar has exactly one way to make
// each molecule, so any order of reductions that finishes gives the same step
// count. When a run gets stuck, the rules are shuffled and the run starts over.
func fewestSteps(rules []replacement, molecule string) int {
	rules = append([]replacement(nil), rules...)

	for {
		steps, ok := reduce(rules, molecule)
		if ok {
			return steps
		}

		rand.Shuffle(len(rules), func(i, j int) { rules[i], rules[j] = rules[j], rules[i] })
	}
}

// reduce shrinks the molecule with the rules in the given order until it is
// the electron. It returns the number of steps and false when it gets stuck.
func reduce(rules []replacement, molecule string) (int, bool) {
	steps := 0

	for molecule != electron {
		before := molecule

		for _, rule := range rules {
			if !applies(rule, molecule) {
				continue
			}

			molecule = strings.Replace(molecule, rule.to, rule.from, 1)
			steps++
		}

		if molecule == before {
			return steps, false
		}
	}

	return steps, true
}

// applies reports whether the rule can run backwards on the molecule. A rule
// from the electron only applies to the whole molecule, because the electron
// is never part of a larger molecule.
func applies(rule replacement, molecule string) bool {
	if rule.from == electron {
		return molecule == rule.to
	}

	return strings.Contains(molecule, rule.to)
}
