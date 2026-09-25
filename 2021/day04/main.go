// Advent of Code 2021, day 4: Giant Squid.
// https://adventofcode.com/2021/day/4
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// size is the width and height of a bingo board.
const size = 5

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 4, part1, part2)
}

// part1 finds the score of the first board to win.
func part1(in string) any {
	scores := play(in)
	if len(scores) == 0 {
		return nil
	}

	return scores[0]
}

// part2 finds the score of the last board to win.
func part2(in string) any {
	scores := play(in)
	if len(scores) == 0 {
		return nil
	}

	return scores[len(scores)-1]
}

// play draws all the numbers and returns the score of each board in the
// order that the boards win. A board wins when all numbers in one row or one
// column are marked. Its score is the sum of its unmarked numbers multiplied
// by the number that made it win.
//
// Each board keeps a count of the marked numbers in each row and column, so
// a win check is two counter updates. A map from a number to its places
// finds the cells to mark without a scan of every board.
func play(in string) []int {
	blocks := input.Blocks(in)
	draws := input.Ints(blocks[0][0])

	type place struct{ board, row, col int }

	var boards [][]int
	where := map[int][]place{}

	for b, block := range blocks[1:] {
		nums := input.Ints(strings.Join(block, " "))
		boards = append(boards, nums)

		for i, n := range nums {
			where[n] = append(where[n], place{b, i / size, i % size})
		}
	}

	rows := make([][size]int, len(boards))
	cols := make([][size]int, len(boards))
	unmarked := make([]int, len(boards))
	won := make([]bool, len(boards))

	for b, nums := range boards {
		for _, n := range nums {
			unmarked[b] += n
		}
	}

	var scores []int

	for _, d := range draws {
		for _, p := range where[d] {
			if won[p.board] {
				continue
			}

			unmarked[p.board] -= d
			rows[p.board][p.row]++
			cols[p.board][p.col]++

			if rows[p.board][p.row] == size || cols[p.board][p.col] == size {
				won[p.board] = true
				scores = append(scores, unmarked[p.board]*d)
			}
		}
	}

	return scores
}
