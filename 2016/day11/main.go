// Advent of Code 2016, day 11: Radioisotope Thermoelectric Generators.
// https://adventofcode.com/2016/day/11
package main

import (
	"regexp"
	"slices"

	"aoc/lib/input"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 11, part1, part2)
}

// floors is the number of floors in the facility. The elevator starts on the
// bottom floor, and every item must end on the top floor.
const floors = 4

// pair holds the floors of one element's microchip and generator. Floors are
// counted from 0 at the bottom.
type pair struct {
	chip, gen int
}

// state is one arrangement of the facility: the elevator floor and the
// floors of every pair.
//
// Which element is which does not matter. Two arrangements that differ only
// by a swap of element names take the same number of moves to finish. So the
// pairs are kept sorted, and all of those arrangements become one state. This
// makes the search much smaller.
//
// Each pair is stored as one byte, chip*floors + gen, in a string. A string
// can be a map key, but a slice can not.
type state struct {
	elevator int
	pairs    string
}

// part1 finds the fewest elevator moves that bring every item to the top floor.
func part1(in string) any {
	return solve(parse(in))
}

// part2 adds an elerium pair and a dilithium pair, both on the first floor.
// In the example, the new generators fry the hydrogen and lithium chips on
// the first floor at once, so part 2 has no answer for the example.
func part2(in string) any {
	return solve(append(parse(in), pair{}, pair{}))
}

// itemPattern matches one item on a line, such as "hydrogen-compatible
// microchip" or "lithium generator". Group 1 is the element and group 2 is
// the kind of item.
var itemPattern = regexp.MustCompile(`(\w+)(?:-compatible)? (microchip|generator)`)

// parse reads where each element's chip and generator start. Line i
// describes floor i.
func parse(in string) []pair {
	index := map[string]int{}
	var pairs []pair

	for floor, line := range input.Lines(in) {
		for _, m := range itemPattern.FindAllStringSubmatch(line, -1) {
			i, ok := index[m[1]]
			if !ok {
				i = len(pairs)
				index[m[1]] = i
				pairs = append(pairs, pair{})
			}

			if m[2] == "microchip" {
				pairs[i].chip = floor
			} else {
				pairs[i].gen = floor
			}
		}
	}

	return pairs
}

// solve runs a breadth-first search from the start to the state where every
// item is on the top floor. It returns nil when a chip is already fried at
// the start, or when no sequence of moves can finish.
func solve(pairs []pair) any {
	if !safe(pairs) {
		return nil
	}

	steps, ok := search.BFS(encode(0, pairs), next, done)
	if !ok {
		return nil
	}

	return steps
}

// encode builds the state for an elevator floor and a list of pairs. Sorting
// the bytes sorts the pairs by chip floor, then by generator floor.
func encode(elevator int, pairs []pair) state {
	b := make([]byte, len(pairs))
	for i, p := range pairs {
		b[i] = byte(p.chip*floors + p.gen)
	}

	slices.Sort(b)

	return state{elevator: elevator, pairs: string(b)}
}

// decode turns the pairs of a state back into a slice that can be changed.
func decode(s state) []pair {
	pairs := make([]pair, len(s.pairs))
	for i := range len(s.pairs) {
		pairs[i] = pair{chip: int(s.pairs[i]) / floors, gen: int(s.pairs[i]) % floors}
	}

	return pairs
}

// done reports whether every chip and generator is on the top floor.
func done(s state) bool {
	const top = (floors-1)*floors + floors - 1

	for i := range len(s.pairs) {
		if s.pairs[i] != top {
			return false
		}
	}

	return true
}

// safe reports whether no chip is fried. A chip is fried when it is on a
// floor with another element's generator and its own generator is elsewhere.
func safe(pairs []pair) bool {
	var hasGen [floors]bool
	for _, p := range pairs {
		hasGen[p.gen] = true
	}

	for _, p := range pairs {
		if p.chip != p.gen && hasGen[p.chip] {
			return false
		}
	}

	return true
}

// next returns every safe state one elevator move away. The elevator goes up
// or down one floor and carries one or two items from its current floor.
func next(s state) []state {
	pairs := decode(s)

	// items points at the floor field of each item on the elevator's floor.
	// A move writes the new floor through the pointer, checks the result, and
	// then writes the old floor back.
	var items []*int
	for i := range pairs {
		if pairs[i].chip == s.elevator {
			items = append(items, &pairs[i].chip)
		}
		if pairs[i].gen == s.elevator {
			items = append(items, &pairs[i].gen)
		}
	}

	var out []state
	for _, to := range []int{s.elevator + 1, s.elevator - 1} {
		if to < 0 || to >= floors {
			continue
		}

		// b starts at a, so a == b is the move that carries one item.
		for a := range items {
			for b := a; b < len(items); b++ {
				*items[a], *items[b] = to, to

				if safe(pairs) {
					out = append(out, encode(to, pairs))
				}

				*items[a], *items[b] = s.elevator, s.elevator
			}
		}
	}

	return out
}
