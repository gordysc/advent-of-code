// Advent of Code 2022, day 13: Distress Signal.
// https://adventofcode.com/2022/day/13
package main

import (
	"cmp"
	"encoding/json"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 13, part1, part2)
}

// part1 adds the indices (from 1) of the pairs that are in the right order.
func part1(in string) any {
	total := 0

	for i, pair := range input.Blocks(in) {
		if compare(parse(pair[0]), parse(pair[1])) < 0 {
			total += i + 1
		}
	}

	return total
}

// part2 finds the decoder key: the product of the positions (from 1) of the
// two divider packets [[2]] and [[6]] after a sort of all packets.
//
// A full sort is not necessary. The position of a divider is one more than
// the number of packets that come before it. [[2]] comes before [[6]], so
// the position of [[6]] also counts [[2]].
func part2(in string) any {
	first, second := parse("[[2]]"), parse("[[6]]")
	before1, before2 := 1, 2

	for _, line := range input.Lines(in) {
		if line == "" {
			continue
		}

		packet := parse(line)
		if compare(packet, first) < 0 {
			before1++
		}

		if compare(packet, second) < 0 {
			before2++
		}
	}

	return before1 * before2
}

// parse reads one packet. Packets use the JSON syntax, so the JSON decoder
// does the work. Into an "any" value, it decodes a list as []any and a
// number as float64.
func parse(line string) any {
	var packet any
	if err := json.Unmarshal([]byte(line), &packet); err != nil {
		panic(err)
	}

	return packet
}

// compare returns a negative number when a comes before b, a positive number
// when b comes before a, and 0 when the order is not decided.
//
// Two numbers compare by value. Two lists compare item by item, and when all
// shared items are equal, the shorter list comes first. When only one value
// is a number, the number becomes a list with that one number.
func compare(a, b any) int {
	// A type assertion with two results reports if the "any" value holds a
	// float64, without a panic when it holds something else.
	an, aNum := a.(float64)
	bn, bNum := b.(float64)

	switch {
	case aNum && bNum:
		return cmp.Compare(an, bn)
	case aNum:
		return compare([]any{an}, b)
	case bNum:
		return compare(a, []any{bn})
	}

	al, bl := a.([]any), b.([]any)

	for i := range min(len(al), len(bl)) {
		if c := compare(al[i], bl[i]); c != 0 {
			return c
		}
	}

	return cmp.Compare(len(al), len(bl))
}
