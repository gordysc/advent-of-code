// Advent of Code 2019, day 14: Space Stoichiometry.
// https://adventofcode.com/2019/day/14
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 14, part1, part2)
}

// cargo is the amount of ORE in the cargo hold for part 2.
const cargo = 1_000_000_000_000

// term is an amount of one chemical, such as "7 A".
type term struct {
	qty  int
	chem string
}

// reaction is one line of the input: it makes qty units of its chemical from
// the inputs, and it can only run a whole number of times.
type reaction struct {
	qty    int
	inputs []term
}

// nanofactory holds every reaction by the chemical it makes, and an order of
// the chemicals where each one comes before all the chemicals it is made from.
type nanofactory struct {
	reactions map[string]reaction
	order     []string
}

// part1 finds the ORE needed to make one FUEL.
func part1(in string) any {
	return parse(in).oreFor(1)
}

// part2 finds the most FUEL the cargo of ORE can make. More FUEL never needs
// less ORE, so a binary search finds the first amount that needs too much.
// The amount just below it is the answer.
func part2(in string) any {
	f := parse(in)

	tooMuch := search.BinarySearch(1, cargo, func(fuel int) bool {
		return f.oreFor(fuel) > cargo
	})

	return tooMuch - 1
}

// parse reads the reactions and puts the chemicals in order.
func parse(in string) nanofactory {
	f := nanofactory{reactions: map[string]reaction{}}

	for _, line := range input.Lines(in) {
		left, right, _ := strings.Cut(line, " => ")
		out := parseTerm(right)

		r := reaction{qty: out.qty}
		for _, part := range strings.Split(left, ", ") {
			r.inputs = append(r.inputs, parseTerm(part))
		}

		f.reactions[out.chem] = r
	}

	f.order = f.sort()

	return f
}

// parseTerm reads an amount and a chemical, such as "7 A".
func parseTerm(s string) term {
	qty, chem, _ := strings.Cut(s, " ")

	return term{qty: input.Int(qty), chem: chem}
}

// sort returns the chemicals so that every chemical comes before the
// chemicals it is made from (a topological order), starting with FUEL.
// A depth-first walk adds each chemical only after all of its inputs, which
// gives the opposite order, so the result is reversed at the end.
func (f nanofactory) sort() []string {
	var order []string
	seen := map[string]bool{}

	// visit is a closure: it is a function value that can read and change
	// the order and seen variables above, and it calls itself to walk down
	// the inputs. It must be declared before it is assigned for that to work.
	var visit func(chem string)
	visit = func(chem string) {
		if seen[chem] {
			return
		}

		seen[chem] = true

		for _, in := range f.reactions[chem].inputs {
			visit(in.chem)
		}

		order = append(order, chem)
	}

	visit("FUEL")

	slices.Reverse(order)

	return order
}

// oreFor finds the ORE needed to make the given amount of FUEL. It walks the
// chemicals in order, so by the time it reaches a chemical, every reaction
// that uses it has already added to its need. So the reaction for it runs
// once, with just enough batches to cover the whole need, and the spare units
// from rounding up are never needed by anything later.
func (f nanofactory) oreFor(fuel int) int {
	need := map[string]int{"FUEL": fuel}

	for _, chem := range f.order {
		r, ok := f.reactions[chem]
		if !ok {
			continue // ORE has no reaction: it is only collected
		}

		runs := mathx.DivCeil(need[chem], r.qty)

		for _, in := range r.inputs {
			need[in.chem] += runs * in.qty
		}
	}

	return need["ORE"]
}
