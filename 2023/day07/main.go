// Advent of Code 2023, day 7: Camel Cards.
// https://adventofcode.com/2023/day/7
package main

import (
	"cmp"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 7, part1, part2)
}

// part1 gives the total winnings when J is a jack.
func part1(in string) any {
	return winnings(in, false)
}

// part2 gives the total winnings when J is a joker. A joker is the weakest
// card when two hands of the same type are compared, but it can act as any
// card to make the strongest type.
func part2(in string) any {
	return winnings(in, true)
}

// hand is one hand with its bid and a score that orders hands by strength.
type hand struct {
	score, bid int
}

// winnings ranks all hands from weakest to strongest and adds each bid times
// its rank. The weakest hand has rank 1.
func winnings(in string, jokers bool) int {
	var hands []hand

	for _, line := range input.Lines(in) {
		cards, bid, _ := strings.Cut(line, " ")
		hands = append(hands, hand{score(cards, jokers), input.Int(bid)})
	}

	slices.SortFunc(hands, func(a, b hand) int { return cmp.Compare(a.score, b.score) })

	total := 0

	for i, h := range hands {
		total += (i + 1) * h.bid
	}

	return total
}

// score turns a hand into one number, so that a stronger hand always has a
// larger score.
//
// The number is written in base 16. Its first digit is the hand type and the
// next five digits are the card values in order. Each card value is less than
// 16, so a better type always wins, and cards only break ties between hands
// of the same type, from the first card to the last.
func score(cards string, jokers bool) int {
	order := "23456789TJQKA"
	if jokers {
		order = "J23456789TQKA"
	}

	s := handType(cards, jokers)

	for i := range len(cards) {
		s = s*16 + strings.IndexByte(order, cards[i])
	}

	return s
}

// handType gives the type of a hand as a number from 0 (high card) to 6
// (five of a kind).
//
// The type depends only on the two largest counts of equal cards. For
// example, counts of 3 and 2 make a full house. With jokers, the best choice
// is always to make every joker a copy of the most common other card. This
// makes the largest count as big as possible, and no type gets better from a
// larger second count instead.
func handType(cards string, jokers bool) int {
	counts := map[rune]int{}

	for _, c := range cards {
		counts[c]++
	}

	wild := 0
	if jokers {
		wild = counts['J']
		delete(counts, 'J')
	}

	sizes := make([]int, 0, len(counts))

	for _, n := range counts {
		sizes = append(sizes, n)
	}

	// Two empty groups make sure that a largest and a second-largest count
	// always exist. A hand of five jokers, for example, has no other group
	// for the jokers to join.
	sizes = append(sizes, 0, 0)
	slices.SortFunc(sizes, func(a, b int) int { return cmp.Compare(b, a) })

	first, second := sizes[0]+wild, sizes[1]

	switch {
	case first == 5:
		return 6
	case first == 4:
		return 5
	case first == 3 && second == 2:
		return 4
	case first == 3:
		return 3
	case first == 2 && second == 2:
		return 2
	case first == 2:
		return 1
	default:
		return 0
	}
}
