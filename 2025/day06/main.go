// Advent of Code 2025, day 6: Trash Compactor.
// https://adventofcode.com/2025/day/6
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 6, part1, part2)
}

// part1 adds the answers of all problems, where each row of a problem holds
// one number.
func part1(in string) any {
	return grandTotal(in, rowNumbers)
}

// part2 adds the answers of all problems, where each column of a problem holds
// one number, with the most significant digit at the top.
//
// The puzzle reads the columns from right to left. Addition and
// multiplication give the same result in any order, so the order is not
// important here.
func part2(in string) any {
	return grandTotal(in, columnNumbers)
}

// grandTotal splits the worksheet into problems, reads the numbers of each
// problem with read, and adds the answers.
//
// A problem is a group of columns. A column that is all spaces separates two
// problems. The last line holds the operator of each problem.
func grandTotal(in string, read func(rows []string, lo, hi int) []int) any {
	// Do not trim the lines: the position of each digit is important in part 2.
	// input.Lines only removes the trailing newlines.
	rows := pad(input.Lines(in))
	if len(rows) < 2 {
		return nil
	}

	ops := rows[len(rows)-1]
	nums := rows[:len(rows)-1]
	width := len(ops)
	total := 0

	for lo := 0; lo < width; {
		if blankColumn(rows, lo) {
			lo++
			continue
		}

		// Find the end of this problem: the next blank column, or the edge.
		hi := lo
		for hi < width && !blankColumn(rows, hi) {
			hi++
		}

		total += solve(ops[lo:hi], read(nums, lo, hi))
		lo = hi
	}

	return total
}

// pad makes all lines the same length. An editor can remove trailing spaces,
// and then the short lines would cause an index out of range.
func pad(lines []string) []string {
	width := 0
	for _, line := range lines {
		width = max(width, len(line))
	}

	padded := make([]string, len(lines))
	for i, line := range lines {
		padded[i] = line + strings.Repeat(" ", width-len(line))
	}

	return padded
}

// blankColumn reports whether column x is a space in every row.
func blankColumn(rows []string, x int) bool {
	for _, row := range rows {
		if row[x] != ' ' {
			return false
		}
	}

	return true
}

// solve finds the operator in the operator row of a problem, and applies it
// to all the numbers.
func solve(op string, nums []int) int {
	if strings.Contains(op, "*") {
		return slicesx.Product(nums)
	}

	return slicesx.Sum(nums)
}

// rowNumbers reads one number from each row, in columns lo to hi-1. The spaces
// around a number only align it, so they are removed.
func rowNumbers(rows []string, lo, hi int) []int {
	var nums []int

	for _, row := range rows {
		text := strings.TrimSpace(row[lo:hi])
		if text == "" {
			continue
		}

		nums = append(nums, input.Int(text))
	}

	return nums
}

// columnNumbers reads one number from each column lo to hi-1. It joins the
// digits of the column from top to bottom and skips the spaces.
func columnNumbers(rows []string, lo, hi int) []int {
	var nums []int

	for x := lo; x < hi; x++ {
		n, digits := 0, 0

		for _, row := range rows {
			c := row[x]
			if c < '0' || c > '9' {
				continue
			}

			// A byte minus '0' gives the value of a digit character.
			n = n*10 + int(c-'0')
			digits++
		}

		if digits > 0 {
			nums = append(nums, n)
		}
	}

	return nums
}
