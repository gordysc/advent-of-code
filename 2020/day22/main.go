// Advent of Code 2020, day 22: Crab Combat.
// https://adventofcode.com/2020/day/22
package main

import (
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 22, part1, part2)
}

// part1 plays normal Combat. In each round the higher card wins, and the
// winner puts both cards on the bottom of their deck, their own card first.
// The game ends when one deck is empty.
func part1(in string) any {
	a, b := parse(in)

	for len(a) > 0 && len(b) > 0 {
		// Re-slicing drops the top card without a copy. append grows the
		// backing array when it runs out of room, so the deck stays valid.
		x, y := a[0], b[0]
		a, b = a[1:], b[1:]

		if x > y {
			a = append(a, x, y)
		} else {
			b = append(b, y, x)
		}
	}

	if len(a) > 0 {
		return score(a)
	}

	return score(b)
}

// part2 plays Recursive Combat and scores the winning deck.
func part2(in string) any {
	a, b := parse(in)

	_, deck := recursive(a, b)

	return score(deck)
}

// recursive plays one game of Recursive Combat. It returns true when
// player 1 wins, and the winner's final deck.
//
// If a round starts with the same two decks as an earlier round of this
// game, player 1 wins the game. The seen set holds a key for each earlier
// round. If both players have at least as many cards as the value they drew,
// a sub-game with copies of that many cards decides the round. Otherwise the
// higher card wins.
func recursive(a, b []int) (bool, []int) {
	seen := map[string]bool{}

	for len(a) > 0 && len(b) > 0 {
		k := key(a, b)
		if seen[k] {
			return true, a
		}
		seen[k] = true

		x, y := a[0], b[0]
		a, b = a[1:], b[1:]

		var aWins bool
		if len(a) >= x && len(b) >= y {
			aWins, _ = recursive(copyN(a, x), copyN(b, y))
		} else {
			aWins = x > y
		}

		if aWins {
			a = append(a, x, y)
		} else {
			b = append(b, y, x)
		}
	}

	if len(a) > 0 {
		return true, a
	}

	return false, b
}

// copyN returns a new slice with the first n cards of deck. The copy keeps
// appends in a sub-game from writing into the parent's deck.
func copyN(deck []int, n int) []int {
	out := make([]int, n)
	copy(out, deck[:n])

	return out
}

// key encodes both decks as a string. Card values are below 256, so each
// card fits in one byte, and a separator byte splits the two decks.
func key(a, b []int) string {
	var sb strings.Builder
	sb.Grow(len(a) + len(b) + 1)

	for _, c := range a {
		sb.WriteByte(byte(c))
	}
	sb.WriteByte(255)
	for _, c := range b {
		sb.WriteByte(byte(c))
	}

	return sb.String()
}

// score multiplies each card by its position counted from the bottom of
// the deck, starting at 1, and adds the results.
func score(deck []int) int {
	total := 0
	for i, c := range deck {
		total += c * (len(deck) - i)
	}

	return total
}

// parse reads the two "Player N:" blocks into two decks, top card first.
func parse(in string) ([]int, []int) {
	blocks := input.Blocks(in)
	decks := make([][]int, 2)

	for i, block := range blocks[:2] {
		for _, line := range block[1:] {
			n, _ := strconv.Atoi(line)
			decks[i] = append(decks[i], n)
		}
	}

	return decks[0], decks[1]
}
