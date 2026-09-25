// Advent of Code 2020, day 23: Crab Cups.
// https://adventofcode.com/2020/day/23
package main

import (
	"strconv"
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 23, part1, part2)
}

// part1 plays 100 moves with the nine cups from the input. The answer is the
// labels of the cups after cup 1, read clockwise, without cup 1.
func part1(in string) any {
	labels := parse(in)
	next := play(labels, len(labels), 100)

	var sb strings.Builder
	for c := next[1]; c != 1; c = next[c] {
		sb.WriteString(strconv.Itoa(int(c)))
	}

	return sb.String()
}

// part2 fills the circle up to one million cups and plays ten million moves.
// The answer is the product of the two cups right after cup 1.
func part2(in string) any {
	next := play(parse(in), 1_000_000, 10_000_000)

	a := int(next[1])
	b := int(next[a])

	return a * b
}

// play runs the game and returns the final circle.
//
// The circle is a linked list stored in a flat slice: next[c] is the label of
// the cup clockwise from cup c. Index 0 is not used. With this layout a move
// changes three links and needs no search, so ten million moves are fast.
// int32 halves the memory of the slice, which helps the CPU cache.
//
// The first cups come from labels. The cups after them have the labels
// len(labels)+1 up to total, in order.
func play(labels []int32, total, moves int) []int32 {
	next := make([]int32, total+1)

	order := make([]int32, 0, total)
	order = append(order, labels...)
	for c := int32(len(labels)) + 1; c <= int32(total); c++ {
		order = append(order, c)
	}

	for i, c := range order {
		next[c] = order[(i+1)%total]
	}

	highest := int32(total)
	cur := order[0]

	for range moves {
		// Pick up the three cups after the current cup and close the gap.
		p1 := next[cur]
		p2 := next[p1]
		p3 := next[p2]
		next[cur] = next[p3]

		// The destination is the next lower label that is not picked up.
		// Below label 1 it wraps round to the highest label.
		dest := cur
		for {
			dest--
			if dest == 0 {
				dest = highest
			}

			if dest != p1 && dest != p2 && dest != p3 {
				break
			}
		}

		// Put the three cups back in, right after the destination.
		next[p3] = next[dest]
		next[dest] = p1

		cur = next[cur]
	}

	return next
}

// parse reads the digit string into cup labels, in clockwise order.
func parse(in string) []int32 {
	var labels []int32

	for _, r := range strings.TrimSpace(in) {
		labels = append(labels, r-'0')
	}

	return labels
}
