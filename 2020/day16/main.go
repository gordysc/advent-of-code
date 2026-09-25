// Advent of Code 2020, day 16: Ticket Translation.
// https://adventofcode.com/2020/day/16
package main

import (
	"math/bits"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 16, part1, part2)
}

// part1 adds up every value on a nearby ticket that no field rule accepts.
// Such a value makes its ticket invalid, and the sum is the error rate.
func part1(in string) any {
	notes := parse(in)

	rate := 0
	for _, ticket := range notes.nearby {
		for _, v := range ticket {
			if !notes.anyAccepts(v) {
				rate += v
			}
		}
	}

	return rate
}

// part2 finds which column holds which field, and multiplies the values of
// the "departure" fields on our own ticket.
//
// Only valid nearby tickets count. For each field we keep a bit mask of the
// columns that every valid ticket accepts. Then we fix the field that has
// exactly one candidate column, remove that column from all other fields,
// and repeat until all fields have a column.
//
// example.txt is the part 1 example. Neither puzzle example has a
// "departure" field, so part 2 has no answer for it and returns nil.
func part2(in string) any {
	notes := parse(in)
	n := len(notes.fields)

	valid := [][]int{notes.mine}
	for _, ticket := range notes.nearby {
		if notes.validTicket(ticket) {
			valid = append(valid, ticket)
		}
	}

	candidates := make([]uint64, n)
	for i, f := range notes.fields {
		for col := range n {
			ok := true
			for _, ticket := range valid {
				if !f.accepts(ticket[col]) {
					ok = false
					break
				}
			}

			if ok {
				candidates[i] |= 1 << col
			}
		}
	}

	column := resolve(candidates)

	product, found := 1, false
	for i, f := range notes.fields {
		if strings.HasPrefix(f.name, "departure") {
			product *= notes.mine[column[i]]
			found = true
		}
	}

	if !found {
		return nil
	}

	return product
}

// resolve turns the candidate column masks into one column per field. It
// repeatedly takes a field with a single candidate and removes that column
// from the other fields. The puzzle guarantees that this always makes
// progress.
func resolve(candidates []uint64) []int {
	column := make([]int, len(candidates))
	done := make([]bool, len(candidates))

	for range candidates {
		for i, mask := range candidates {
			if done[i] || bits.OnesCount64(mask) != 1 {
				continue
			}

			col := bits.TrailingZeros64(mask)
			column[i] = col
			done[i] = true

			for j := range candidates {
				if j != i {
					candidates[j] &^= mask
				}
			}

			break
		}
	}

	return column
}

// field is one rule from the notes: a name and two inclusive ranges.
type field struct {
	name               string
	lo1, hi1, lo2, hi2 int
}

// accepts tells if v is in one of the two ranges of the field.
func (f field) accepts(v int) bool {
	return (v >= f.lo1 && v <= f.hi1) || (v >= f.lo2 && v <= f.hi2)
}

// notes holds the three sections of the input.
type notes struct {
	fields []field
	mine   []int
	nearby [][]int
}

// anyAccepts tells if at least one field rule accepts v.
func (n notes) anyAccepts(v int) bool {
	for _, f := range n.fields {
		if f.accepts(v) {
			return true
		}
	}

	return false
}

// validTicket tells if some field rule accepts each value on the ticket.
func (n notes) validTicket(ticket []int) bool {
	for _, v := range ticket {
		if !n.anyAccepts(v) {
			return false
		}
	}

	return true
}

// parse reads the three blank-line separated sections: the field rules, our
// ticket and the nearby tickets. The ticket sections start with a header
// line, which we skip. The dash in "1-3" is a separator, not a minus sign, so
// the ranges use input.UInts.
func parse(in string) notes {
	blocks := input.Blocks(in)
	var n notes

	for _, line := range blocks[0] {
		name, rest, _ := strings.Cut(line, ": ")
		r := input.UInts(rest)
		n.fields = append(n.fields, field{name, r[0], r[1], r[2], r[3]})
	}

	n.mine = input.Ints(blocks[1][1])

	for _, line := range blocks[2][1:] {
		n.nearby = append(n.nearby, input.Ints(line))
	}

	return n
}
