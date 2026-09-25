// Advent of Code 2018, day 14: Chocolate Charts.
// https://adventofcode.com/2018/day/14
package main

import (
	"bytes"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 14, part1, part2)
}

// scoreCount is how many scores part 1 reads after the given recipes.
const scoreCount = 10

// capacity is how many scores the board has room for at the start. Real
// inputs need about 21 million recipes in part 2, so this avoids most of the
// slice growth.
const capacity = 32 << 20

// board is the scoreboard and the recipes the two elves are on.
type board struct {
	scores []byte
	a, b   int
}

// newBoard returns the board at the start: scores 3 and 7, one elf on each.
func newBoard() *board {
	scores := make([]byte, 2, capacity)
	scores[0], scores[1] = 3, 7

	return &board{scores: scores, a: 0, b: 1}
}

// part1 returns the ten scores that come right after the number of recipes
// in the input.
func part1(in string) any {
	n := input.Int(in)
	bd := newBoard()

	for len(bd.scores) < n+scoreCount {
		bd.step()
	}

	var sb strings.Builder
	for _, s := range bd.scores[n : n+scoreCount] {
		sb.WriteByte('0' + s)
	}

	return sb.String()
}

// part2 reads the input as a sequence of digits, leading zeros included, and
// counts the recipes to the left of where that sequence first appears. One
// step can add two scores, so the sequence can end at either of them.
// The puzzle gives no part 2 answer for the example 2018. The sequence
// 2,0,1,8 first appears after 86764 recipes.
func part2(in string) any {
	var want []byte
	for _, d := range strings.TrimSpace(in) {
		want = append(want, byte(d-'0'))
	}

	bd := newBoard()
	checked := 0

	for {
		bd.step()

		// Check each new end position that has room for the whole sequence.
		for ; checked+len(want) <= len(bd.scores); checked++ {
			if bytes.Equal(bd.scores[checked:checked+len(want)], want) {
				return checked
			}
		}
	}
}

// step makes the new recipes from the sum of the two current scores, then
// moves each elf forward by one plus the score of its recipe.
func (bd *board) step() {
	sum := bd.scores[bd.a] + bd.scores[bd.b]
	if sum >= 10 {
		bd.scores = append(bd.scores, 1, sum-10)
	} else {
		bd.scores = append(bd.scores, sum)
	}

	n := len(bd.scores)
	bd.a = (bd.a + 1 + int(bd.scores[bd.a])) % n
	bd.b = (bd.b + 1 + int(bd.scores[bd.b])) % n
}
