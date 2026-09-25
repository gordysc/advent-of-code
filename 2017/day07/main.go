// Advent of Code 2017, day 7: Recursive Circus.
// https://adventofcode.com/2017/day/7
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 7, part1, part2)
}

// program is one disc in the tower: its own weight and the programs it holds up.
type program struct {
	weight int
	above  []string
}

// tower maps each program name to its program.
type tower map[string]program

// part1 finds the program at the bottom of the tower.
func part1(in string) any {
	return parse(in).bottom()
}

// part2 finds the weight that the one wrong program must have to balance
// the tower.
func part2(in string) any {
	t := parse(in)
	_, fix, _ := t.balance(t.bottom())

	return fix
}

// parse reads lines such as "fwft (72) -> ktlj, cntj, xhth".
func parse(in string) tower {
	t := tower{}

	for _, line := range input.Lines(in) {
		// strings.Cut splits at the first separator. A line with no " -> "
		// leaves list empty.
		name, rest, _ := strings.Cut(line, " ")
		weight, list, _ := strings.Cut(rest, " -> ")

		p := program{weight: input.Int(strings.Trim(weight, "()"))}
		if list != "" {
			p.above = strings.Split(list, ", ")
		}

		t[name] = p
	}

	return t
}

// bottom finds the one program that no other program holds up.
func (t tower) bottom() string {
	held := map[string]bool{}
	for _, p := range t {
		for _, name := range p.above {
			held[name] = true
		}
	}

	for name := range t {
		if !held[name] {
			return name
		}
	}

	return ""
}

// balance returns the total weight of the sub-tower on top of name. The
// search goes bottom-up, so the first unbalanced disc it finds is the deepest
// one, and one of its children is the wrong program. When it finds that
// program, it returns the corrected weight and true.
func (t tower) balance(name string) (int, int, bool) {
	p := t[name]

	totals := make([]int, len(p.above))
	for i, child := range p.above {
		total, fix, found := t.balance(child)
		if found {
			return 0, fix, true
		}

		totals[i] = total
	}

	// A child whose total differs from all the others is the wrong one. A
	// disc with an unbalanced child always holds at least three programs, so
	// the odd one out is clear.
	for i, total := range totals {
		other := totals[(i+1)%len(totals)]
		if total != other && other == totals[(i+2)%len(totals)] {
			return 0, t[p.above[i]].weight + other - total, true
		}
	}

	return p.weight + slicesx.Sum(totals), 0, false
}
