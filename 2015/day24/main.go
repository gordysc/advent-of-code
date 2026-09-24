// Advent of Code 2015, day 24: It Hangs in the Balance.
// https://adventofcode.com/2015/day/24
package main

import (
	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 24, part1, part2)
}

// part1 splits the packages into three groups of equal weight and returns the
// quantum entanglement of the best passenger compartment group.
func part1(in string) any {
	return bestEntanglement(input.IntLines(in), 3)
}

// part2 does the same with four groups, because the trunk now holds one too.
func part2(in string) any {
	return bestEntanglement(input.IntLines(in), 4)
}

// bestEntanglement finds the group for the passenger compartment. That group
// has the fewest packages of any valid split, and among groups of that size,
// the smallest product of weights (the "quantum entanglement"). The result is
// that product, or nil when the weights cannot be split into equal groups.
func bestEntanglement(weights []int, groups int) any {
	total := slicesx.Sum(weights)

	if total%groups != 0 {
		return nil
	}

	target := total / groups

	// Try the smallest group sizes first. The first size that gives any valid
	// split wins, because a smaller front group always beats a smaller product.
	for size := 1; size <= len(weights); size++ {
		best := 0

		for chosen := range slicesx.Combinations(weights, size) {
			if slicesx.Sum(chosen) != target {
				continue
			}

			// A group is only valid when the packages left over can still fill
			// the other groups evenly. Check the product first: it is cheap, and
			// a group that cannot beat the best so far never needs the check.
			product := slicesx.Product(chosen)

			if best != 0 && product >= best {
				continue
			}

			if canSplit(without(weights, chosen), groups-1, target) {
				best = product
			}
		}

		if best != 0 {
			return best
		}
	}

	return nil
}

// canSplit reports whether weights can be divided into exactly groups groups
// that each weigh target. The caller guarantees the weights add up to
// groups*target, so one remaining group always works.
func canSplit(weights []int, groups, target int) bool {
	if groups == 1 {
		return true
	}

	// Depth-first search over the packages, deciding for each one whether it
	// joins the current group. Every time the group reaches the target, the
	// leftover packages must in turn split into the remaining groups.
	inGroup := make([]bool, len(weights))

	var search func(i, remaining int) bool

	search = func(i, remaining int) bool {
		if remaining == 0 {
			rest := make([]int, 0, len(weights))

			for j, w := range weights {
				if !inGroup[j] {
					rest = append(rest, w)
				}
			}

			return canSplit(rest, groups-1, target)
		}

		if i == len(weights) {
			return false
		}

		// Take this package if it fits, then try leaving it out. The inGroup
		// flag is reset before the second branch so the slice stays correct.
		if weights[i] <= remaining {
			inGroup[i] = true

			if search(i+1, remaining-weights[i]) {
				return true
			}

			inGroup[i] = false
		}

		return search(i+1, remaining)
	}

	return search(0, target)
}

// without returns a copy of weights with one occurrence of each value in
// chosen removed. Chosen is a subset of weights, so every value is found.
func without(weights, chosen []int) []int {
	rest := make([]int, 0, len(weights)-len(chosen))
	skip := slicesx.Counts(chosen)

	for _, w := range weights {
		if skip[w] > 0 {
			skip[w]--
			continue
		}

		rest = append(rest, w)
	}

	return rest
}
