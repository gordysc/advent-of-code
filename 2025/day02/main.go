// Advent of Code 2025, day 2: Gift Shop.
// https://adventofcode.com/2025/day/2
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 2, part1, part2)
}

// part1 adds the invalid IDs in all ranges, where an invalid ID is a block of
// digits repeated exactly two times, such as 6464.
func part1(in string) any {
	return invalidSum(in, func(repeats int) bool { return repeats == 2 })
}

// part2 adds the invalid IDs in all ranges, where an invalid ID is a block of
// digits repeated two or more times, such as 6464 or 123123123.
func part2(in string) any {
	return invalidSum(in, func(repeats int) bool { return repeats >= 2 })
}

// invalidSum adds the invalid IDs of every range. The allow function tells
// which repeat counts make an ID invalid.
//
// The ranges can be very wide, so we do not test each ID. We make each
// repeated number directly. A number of L digits that is a block of k digits
// repeated r times (L = k*r) is block * m, where m is 1 followed by r-1 groups
// of "0...01" (for k=3, r=3: m = 1001001). So the invalid IDs of that shape in
// a range are block * m for each block in a range that we can calculate.
//
// One ID can have more than one shape: 222222 is "222" two times, "22" three
// times and "2" six times. A set per range makes sure we add it only once.
func invalidSum(in string, allow func(repeats int) bool) int {
	nums := input.UInts(in)
	total := 0

	for i := 0; i+1 < len(nums); i += 2 {
		lo, hi := nums[i], nums[i+1]
		found := set.New[int]()

		for length := mathx.NumDigits(lo); length <= mathx.NumDigits(hi); length++ {
			for repeats := 2; repeats <= length; repeats++ {
				if length%repeats != 0 || !allow(repeats) {
					continue
				}

				addRepeated(found, lo, hi, length/repeats, repeats)
			}
		}

		for id := range found.All() {
			total += id
		}
	}

	return total
}

// addRepeated adds to found each number from lo to hi that is a block of k
// digits repeated r times.
func addRepeated(found set.Set[int], lo, hi, k, r int) {
	unit := mathx.Pow(10, k)

	m := 0
	for range r {
		m = m*unit + 1
	}

	// The block must have exactly k digits, because an ID has no leading zero.
	first := max(unit/10, mathx.DivCeil(lo, m))
	last := min(unit-1, hi/m)

	for block := first; block <= last; block++ {
		found.Add(block * m)
	}
}
