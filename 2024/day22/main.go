// Advent of Code 2024, day 22: Monkey Market.
// https://adventofcode.com/2024/day/22
package main

import (
	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 22, part1, part2)
}

// rounds is the number of new secret numbers each buyer makes in a day.
const rounds = 2000

// part1 adds the 2000th new secret number of each buyer.
//
// The puzzle has two examples. The file example.txt holds the part 1 example
// (1, 10, 100, 2024). On this example, part 1 gives 37327623 and part 2
// gives 24. On the part 2 example (1, 2, 3, 2024), part 1 gives 37990510 and
// part 2 gives 23.
func part1(in string) any {
	total := 0

	for _, s := range input.IntLines(in) {
		for range rounds {
			s = next(s)
		}

		total += s
	}

	return total
}

// part2 finds the sequence of four price changes that gets the most bananas
// over all buyers, and returns that number of bananas.
//
// Each change is between -9 and 9, so it has 19 possible values. A sequence of
// four changes is then a four-digit number in base 19. We use that number as
// an index into a flat slice. This is much faster than a map with the
// sequence as the key.
func part2(in string) any {
	const size = 19 * 19 * 19 * 19

	bananas := make([]int, size)

	// seen[i] holds the last buyer (plus one) that sold on sequence i. A
	// monkey sells at the first match only, so later matches of the same
	// buyer must not count. Buyer numbers make this work without a reset of
	// the slice for each buyer.
	seen := make([]int, size)

	for buyer, s := range input.IntLines(in) {
		idx := 0
		price := s % 10

		for i := range rounds {
			s = next(s)
			newPrice := s % 10

			// Shift in the new change and drop the oldest one, which is the
			// highest base-19 digit.
			idx = (idx*19 + newPrice - price + 9) % size
			price = newPrice

			// The index holds a full sequence only after four changes.
			if i < 3 || seen[idx] == buyer+1 {
				continue
			}

			seen[idx] = buyer + 1
			bananas[idx] += price
		}
	}

	_, best := slicesx.MinMax(bananas)

	return best
}

// next returns the secret number that follows s.
//
// The puzzle steps use multiply by 64, divide by 32, multiply by 2048, and
// modulo 16777216. All are powers of two, so shifts and a bit mask do the
// same work. The mask 0xFFFFFF keeps the low 24 bits, which is modulo 2^24.
func next(s int) int {
	s = (s ^ s<<6) & 0xFFFFFF

	// A right shift makes the value smaller, so it needs no mask.
	s ^= s >> 5

	return (s ^ s<<11) & 0xFFFFFF
}
