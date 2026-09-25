// Advent of Code 2019, day 7: Amplification Circuit.
// https://adventofcode.com/2019/day/7
package main

import (
	"aoc/lib/intcode"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 7, part1, part2)
}

// part1 finds the highest signal the five amplifiers can send to the
// thrusters when they run once in series with phase settings 0 to 4.
//
// example.txt is the first part 2 example. Its amplifiers are made for the
// feedback loop, so they wait for more input instead of halting, and part 1
// has no answer for it.
func part1(in string) any {
	program := intcode.Parse(in)

	best := 0
	for phases := range slicesx.Permutations([]int{0, 1, 2, 3, 4}) {
		signal, ok := series(program, phases)
		if !ok {
			return nil
		}

		best = max(best, signal)
	}

	return best
}

// part2 finds the highest signal the amplifiers can send when they run in a
// feedback loop with phase settings 5 to 9. The part 1 examples are made to
// run in series, so part 2 has no answer for them.
func part2(in string) any {
	program := intcode.Parse(in)

	best := 0
	for phases := range slicesx.Permutations([]int{5, 6, 7, 8, 9}) {
		signal, ok := feedback(program, phases)
		if !ok {
			return nil
		}

		best = max(best, signal)
	}

	return best
}

// series runs one amplifier per phase setting. Each one reads its phase and
// the previous amplifier's output (0 for the first one), and its output goes
// to the next. It returns the last output, and false if an amplifier did not
// halt, because then the program was not made to run in series.
func series(program []int, phases []int) (int, bool) {
	signal := 0

	for _, phase := range phases {
		m := intcode.New(program)
		m.Send(phase, signal)

		out := m.Run()
		if !m.Halted() || len(out) == 0 {
			return 0, false
		}

		signal = out[len(out)-1]
	}

	return signal, true
}

// feedback connects the amplifiers in a loop: E's output goes back to A. The
// machines keep their state between turns, because Run stops when a machine
// needs input and the next Run carries on from there. The loop goes around
// until E halts, and E's last output is the signal for the thrusters.
//
// A program made for the loop sends signals around it more than once. If E
// halts after the first pass, the program was made to run in series, so
// feedback returns false.
func feedback(program []int, phases []int) (int, bool) {
	amps := make([]*intcode.Machine, len(phases))
	for i, phase := range phases {
		amps[i] = intcode.New(program)
		amps[i].Send(phase)
	}

	signals := []int{0}
	last, passes := 0, 0

	for !amps[len(amps)-1].Halted() {
		for _, amp := range amps {
			amp.Send(signals...)
			signals = amp.Run()
		}

		if len(signals) > 0 {
			last = signals[len(signals)-1]
		}

		passes++
	}

	return last, passes > 1
}
