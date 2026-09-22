// Advent of Code 2015, day 20: Infinite Elves and Infinite Houses.
// https://adventofcode.com/2015/day/20
//
// The puzzle text has no example input, only a table of the presents at the
// first nine houses. The example.txt file here asks for 130 presents. From the
// table, part 1 gives house 8 (150 presents) and part 2 gives house 6 (132).
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 20, part1, part2)
}

// part1 finds the lowest house that gets at least the target number of
// presents when every elf visits every multiple of its number forever.
func part1(in string) any {
	target := input.Int(in)

	// Every house n gets at least 10n presents from elf n alone, so the search
	// never has to go past house target/10.
	return firstHouse(target, mathx.DivCeil(target, 10), 10, 0)
}

// part2 is the same search, but every elf stops after 50 houses and leaves 11
// presents at each one.
func part2(in string) any {
	target := input.Int(in)

	// Elf n still visits house n first, so house target/11 is enough.
	return firstHouse(target, mathx.DivCeil(target, 11), 11, 50)
}

// firstHouse returns the lowest house with at least target presents. It uses
// a sieve: each elf adds its presents to every house it visits, up to the
// house at limit. When maxVisits is 0, an elf never stops.
func firstHouse(target, limit, perVisit, maxVisits int) int {
	presents := make([]int, limit+1)

	for elf := 1; elf <= limit; elf++ {
		last := limit
		if maxVisits > 0 {
			last = min(limit, elf*maxVisits)
		}

		for house := elf; house <= last; house += elf {
			presents[house] += elf * perVisit
		}
	}

	for house := 1; house <= limit; house++ {
		if presents[house] >= target {
			return house
		}
	}

	panic("no house reached the target")
}
