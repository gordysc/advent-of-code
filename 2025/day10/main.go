// Advent of Code 2025, day 10: Factory.
// https://adventofcode.com/2025/day/10
package main

import (
	"encoding/binary"
	"math"
	"math/bits"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2025, 10, part1, part2)
}

// part1 returns the fewest button presses that set the indicator lights of
// every machine to their pattern.
//
// A press toggles lights, so two presses of the same button cancel. Thus each
// button is pressed zero times or one time, and we can try every subset of the
// buttons. A machine has few buttons, so this is fast.
func part1(in string) any {
	machines := parse(in)
	if len(machines) == 0 {
		return nil
	}

	total := 0

	for _, m := range machines {
		best := math.MaxInt

		for subset := range 1 << len(m.buttons) {
			if m.lightsAfter(subset) == m.lights {
				// bits.OnesCount gives the number of 1 bits, which is the
				// number of buttons in the subset.
				best = min(best, bits.OnesCount(uint(subset)))
			}
		}

		if best == math.MaxInt {
			return nil
		}

		total += best
	}

	return total
}

// part2 returns the fewest button presses that set the joltage counters of
// every machine to their required values. Each press adds 1 to each counter
// that the button is wired to.
//
// This is an integer linear program. We solve it exactly with a halving
// recursion. Look at any solution, and let odd be the set of buttons that it
// presses an odd number of times. Press each button in odd once. The rest of
// the presses come in pairs, so the remaining targets are all even, and half
// of those targets is again a machine to solve. We do not know odd, but the
// parity of the targets tells which subsets can be odd. So we try each of them
// and keep the minimum. The targets halve at each step, so the recursion is
// shallow.
func part2(in string) any {
	machines := parse(in)
	if len(machines) == 0 {
		return nil
	}

	total := 0

	for _, m := range machines {
		presses := newSolver(m).solve(m.joltage)
		if presses == unsolvable {
			return nil
		}

		total += presses
	}

	return total
}

// machine is one line of the manual.
type machine struct {
	// lights has bit i set when light i must be on.
	lights int
	// buttons lists, for each button, the indexes of the lights and counters
	// that it is wired to.
	buttons [][]int
	// joltage gives the required value of each counter.
	joltage []int
}

// parse reads one machine per line. The fields have the forms "[.##.]",
// "(0,3)" and "{3,5,4,7}".
func parse(in string) []machine {
	var machines []machine

	for _, line := range input.Lines(in) {
		var m machine

		for _, field := range strings.Fields(line) {
			switch field[0] {
			case '[':
				for i, c := range field[1 : len(field)-1] {
					if c == '#' {
						m.lights |= 1 << i
					}
				}
			case '(':
				m.buttons = append(m.buttons, input.Ints(field))
			case '{':
				m.joltage = input.Ints(field)
			}
		}

		machines = append(machines, m)
	}

	return machines
}

// lightsAfter returns the lights that are on after one press of each button in
// subset. Bit i of subset selects button i.
func (m machine) lightsAfter(subset int) int {
	lights := 0

	for i, button := range m.buttons {
		if subset&(1<<i) == 0 {
			continue
		}

		for _, light := range button {
			lights ^= 1 << light
		}
	}

	return lights
}

// unsolvable is the result for targets that no presses can reach.
const unsolvable = math.MaxInt

// press is one subset of buttons, each pressed once.
type press struct {
	// count is the number of buttons in the subset.
	count int
	// adds gives how much the subset adds to each counter.
	adds []int
}

// solver finds the fewest presses for the counters of one machine.
type solver struct {
	// byParity groups the button subsets by the counters that they change by
	// an odd amount. The key has bit i set when counter i changes by an odd
	// amount.
	byParity map[int][]press
	// memo keeps the result for each target that we already solved.
	memo map[string]int
}

// newSolver finds the effect of every button subset of m once, so that the
// recursion only has to look them up.
func newSolver(m machine) solver {
	s := solver{byParity: map[int][]press{}, memo: map[string]int{}}

	for subset := range 1 << len(m.buttons) {
		p := press{adds: make([]int, len(m.joltage))}

		for i, button := range m.buttons {
			if subset&(1<<i) == 0 {
				continue
			}

			p.count++
			for _, counter := range button {
				p.adds[counter]++
			}
		}

		s.byParity[parity(p.adds)] = append(s.byParity[parity(p.adds)], p)
	}

	return s
}

// parity returns a mask with bit i set when values[i] is odd.
func parity(values []int) int {
	mask := 0

	for i, v := range values {
		mask |= (v & 1) << i
	}

	return mask
}

// solve returns the fewest presses that add exactly target to the counters,
// or unsolvable.
func (s solver) solve(target []int) int {
	if isZero(target) {
		return 0
	}

	key := keyOf(target)
	if n, ok := s.memo[key]; ok {
		return n
	}

	best := unsolvable
	rest := make([]int, len(target))

	// Only subsets with the same parity as the target leave even targets.
	for _, p := range s.byParity[parity(target)] {
		if !subtractHalf(target, p.adds, rest) {
			continue
		}

		// Each press in the half problem stands for two presses here.
		half := s.solve(rest)
		if half == unsolvable {
			continue
		}

		best = min(best, p.count+2*half)
	}

	s.memo[key] = best

	return best
}

// subtractHalf writes (target - adds) / 2 into rest. It returns false when
// adds is larger than target for some counter. The caller makes sure that
// target - adds is even.
func subtractHalf(target, adds, rest []int) bool {
	for i := range target {
		d := target[i] - adds[i]
		if d < 0 {
			return false
		}

		rest[i] = d / 2
	}

	return true
}

// isZero reports whether every value is zero.
func isZero(values []int) bool {
	for _, v := range values {
		if v != 0 {
			return false
		}
	}

	return true
}

// keyOf packs values into a string, so that the targets can be map keys.
// A slice cannot be a map key in Go, but a string can.
func keyOf(values []int) string {
	var b []byte

	// AppendUvarint writes each number in as few bytes as it needs.
	for _, v := range values {
		b = binary.AppendUvarint(b, uint64(v))
	}

	return string(b)
}
