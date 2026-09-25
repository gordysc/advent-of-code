// Advent of Code 2024, day 21: Keypad Conundrum.
// https://adventofcode.com/2024/day/21
package main

import (
	"math"
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 21, part1, part2)
}

// part1 adds the complexities of the codes when two robots with directional
// keypads stand between you and the robot at the door.
//
// You press the third directional keypad yourself, so the numeric keypad has
// three directional keypads above it.
func part1(in string) any {
	return complexity(in, 3)
}

// part2 adds the complexities of the codes when 25 robots with directional
// keypads stand between you and the robot at the door.
//
// The puzzle gives no answer for the example here. With 25 robots, the
// example codes give 154115708116294.
func part2(in string) any {
	return complexity(in, 26)
}

// keypad maps each button to its position. It also records the gap, because
// a robot arm must never point at the gap.
type keypad struct {
	keys map[byte]grid.Point
	gap  grid.Point
}

// newKeypad reads a keypad layout, one string per row. A space marks the gap.
func newKeypad(rows ...string) *keypad {
	k := &keypad{keys: map[byte]grid.Point{}}

	for y, row := range rows {
		for x := range len(row) {
			if row[x] == ' ' {
				k.gap = grid.P(x, y)
				continue
			}

			k.keys[row[x]] = grid.P(x, y)
		}
	}

	return k
}

// numeric is the keypad on the door. directional is the keypad that each
// robot, and you, use to control the next robot.
var (
	numeric     = newKeypad("789", "456", "123", " 0A")
	directional = newKeypad(" ^A", "<v>")
)

// paths returns the button presses that move an arm from one key to another
// and then press "A".
//
// Only two paths can be the shortest after the presses go up the chain: all
// horizontal moves first, or all vertical moves first. A path that switches
// direction more than once makes the robot above it travel more. A path is
// skipped when its corner is the gap.
func (k *keypad) paths(from, to byte) []string {
	a, b := k.keys[from], k.keys[to]
	d := b.Sub(a)

	horizontal := strings.Repeat(">", max(d.X, 0)) + strings.Repeat("<", max(-d.X, 0))
	vertical := strings.Repeat("v", max(d.Y, 0)) + strings.Repeat("^", max(-d.Y, 0))

	var out []string

	if grid.P(b.X, a.Y) != k.gap {
		out = append(out, horizontal+vertical+"A")
	}

	if grid.P(a.X, b.Y) != k.gap {
		out = append(out, vertical+horizontal+"A")
	}

	return out
}

// moveKey identifies one arm move on one keypad, with the number of
// directional keypads above that keypad.
type moveKey struct {
	pad      *keypad
	from, to byte
	above    int
}

// solver remembers the cost of each arm move. The same few moves come up
// again and again, so this keeps part 2 fast.
type solver struct {
	memo map[moveKey]int
}

// complexity adds, for each code, the number of presses you make times the
// numeric part of the code. above is the number of directional keypads above
// the numeric keypad, yours included.
func complexity(in string, above int) int {
	s := solver{memo: map[moveKey]int{}}
	total := 0

	for _, code := range input.Lines(in) {
		presses := s.typeCost(numeric, code, above)
		total += presses * input.UInts(code)[0]
	}

	return total
}

// typeCost returns the number of presses you make so that keypad k gets the
// buttons in seq pressed in order. Every arm starts on "A".
//
// When no keypad is above k, you press k yourself, one press per button.
func (s *solver) typeCost(k *keypad, seq string, above int) int {
	if above == 0 {
		return len(seq)
	}

	total := 0
	from := byte('A')

	for i := range len(seq) {
		total += s.moveCost(k, from, seq[i], above)
		from = seq[i]
	}

	return total
}

// moveCost returns the number of presses you make so that the arm over
// keypad k moves from one key to another and presses it.
//
// Every move ends with "A" on the keypad above. So all arms above are back on
// "A" after each move, and each move costs the same no matter what came
// before it.
func (s *solver) moveCost(k *keypad, from, to byte, above int) int {
	key := moveKey{k, from, to, above}
	if c, ok := s.memo[key]; ok {
		return c
	}

	best := math.MaxInt
	for _, p := range k.paths(from, to) {
		best = min(best, s.typeCost(directional, p, above-1))
	}

	s.memo[key] = best

	return best
}
