// Advent of Code 2022, day 11: Monkey in the Middle.
// https://adventofcode.com/2022/day/11
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 11, part1, part2)
}

// monkey holds the notes about one monkey.
type monkey struct {
	items   []int // worry levels of the items that the monkey holds
	mul     bool  // true when the operation multiplies, false when it adds
	operand int   // the other value of the operation; 0 means "old"
	divisor int   // the monkey tests if the worry level divides by this
	ifTrue  int   // the monkey that gets the item when the test passes
	ifFalse int   // the monkey that gets the item when the test fails
}

// part1 finds the level of monkey business after 20 rounds. After each
// inspection, the worry level is divided by 3.
func part1(in string) any {
	return monkeyBusiness(parse(in), 20, func(w int) int { return w / 3 })
}

// part2 finds the level of monkey business after 10000 rounds, without the
// division by 3.
//
// The worry levels become too large for an int. But each test only asks if
// a level divides by a monkey's divisor. Taking the level modulo the product
// of all divisors keeps the result of every test the same, and keeps the
// numbers small.
func part2(in string) any {
	monkeys := parse(in)

	product := 1
	for _, m := range monkeys {
		product *= m.divisor
	}

	return monkeyBusiness(monkeys, 10000, func(w int) int { return w % product })
}

// monkeyBusiness plays the given number of rounds and returns the product of
// the two highest inspection counts. The relief function changes the worry
// level after each inspection.
func monkeyBusiness(monkeys []monkey, rounds int, relief func(int) int) int {
	counts := make([]int, len(monkeys))

	for range rounds {
		// Use the index, because a range copy of the struct would not see the
		// items that other monkeys throw to it during this round.
		for i := range monkeys {
			m := &monkeys[i]
			counts[i] += len(m.items)

			for _, worry := range m.items {
				worry = relief(m.inspect(worry))

				target := m.ifFalse
				if worry%m.divisor == 0 {
					target = m.ifTrue
				}

				monkeys[target].items = append(monkeys[target].items, worry)
			}

			// Keep the backing array, so later rounds can reuse its space.
			m.items = m.items[:0]
		}
	}

	slices.Sort(counts)
	n := len(counts)

	return counts[n-1] * counts[n-2]
}

// inspect applies the monkey's operation to a worry level.
func (m *monkey) inspect(worry int) int {
	operand := m.operand
	if operand == 0 {
		operand = worry
	}

	if m.mul {
		return worry * operand
	}

	return worry + operand
}

// parse reads the notes for each monkey. The monkeys come in order, so the
// index in the slice is the monkey number.
func parse(in string) []monkey {
	var monkeys []monkey

	for _, block := range input.Blocks(in) {
		m := monkey{
			items:   input.Ints(block[1]),
			mul:     strings.Contains(block[2], "*"),
			divisor: input.Ints(block[3])[0],
			ifTrue:  input.Ints(block[4])[0],
			ifFalse: input.Ints(block[5])[0],
		}

		// "new = old * old" has no number, so the operand stays 0.
		if nums := input.Ints(block[2]); len(nums) > 0 {
			m.operand = nums[0]
		}

		monkeys = append(monkeys, m)
	}

	return monkeys
}
