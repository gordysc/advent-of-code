// Advent of Code 2020, day 7: Handy Haversacks.
// https://adventofcode.com/2020/day/7
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// target is the colour of the bag that both parts ask about.
const target = "shiny gold"

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 7, part1, part2)
}

// part1 counts the bag colours that can hold a shiny gold bag at some depth.
// It inverts the rules into a map from each colour to the colours that hold
// it directly, then walks that map outwards from shiny gold.
//
// example.txt is the first example. It gives 4 for part 1 and 32 for part 2.
func part1(in string) any {
	rules := parse(in)

	holders := map[string][]string{}
	for outer, inner := range rules {
		for colour := range inner {
			holders[colour] = append(holders[colour], outer)
		}
	}

	seen := map[string]bool{}
	stack := []string{target}
	for len(stack) > 0 {
		colour := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, outer := range holders[colour] {
			if !seen[outer] {
				seen[outer] = true
				stack = append(stack, outer)
			}
		}
	}

	return len(seen)
}

// part2 counts the bags that one shiny gold bag must hold, at all depths.
func part2(in string) any {
	rules := parse(in)

	return inside(target, rules, map[string]int{})
}

// inside counts the bags inside one bag of the given colour. Each child bag
// adds itself plus its own contents. The memo keeps each colour to one
// calculation, because many colours share the same children.
func inside(colour string, rules map[string]map[string]int, memo map[string]int) int {
	if n, ok := memo[colour]; ok {
		return n
	}

	total := 0
	for child, n := range rules[colour] {
		total += n * (1 + inside(child, rules, memo))
	}

	memo[colour] = total

	return total
}

// parse reads each rule into a map from the outer colour to the colours it
// holds and how many of each. A rule looks like
// "light red bags contain 1 bright white bag, 2 muted yellow bags.".
// A colour is always two words, so each child is "<n> <adj> <colour> bag(s)".
func parse(in string) map[string]map[string]int {
	rules := map[string]map[string]int{}

	for _, line := range input.Lines(in) {
		outer, rest, _ := strings.Cut(line, " bags contain ")
		inner := map[string]int{}

		if rest != "no other bags." {
			for _, part := range strings.Split(rest, ", ") {
				f := strings.Fields(part)
				inner[f[1]+" "+f[2]] = input.Int(f[0])
			}
		}

		rules[outer] = inner
	}

	return rules
}
