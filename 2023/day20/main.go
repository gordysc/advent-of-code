// Advent of Code 2023, day 20: Pulse Propagation.
// https://adventofcode.com/2023/day/20
//
// The puzzle text has two examples. example.txt holds the second one, which
// gives 11687500 for part 1. No example has the module rx, so part 2 gives
// no answer for the example.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 20, part1, part2)
}

// module is one node in the network. kind is '%' for a flip-flop, '&' for a
// conjunction, 'b' for the broadcaster, and 0 for a module that only receives
// pulses (such as output or rx).
type module struct {
	kind  byte
	dests []int
	// on is the state of a flip-flop.
	on bool
	// last holds the most recent pulse from each input of a conjunction,
	// keyed by the index of the input module. true is a high pulse.
	last map[int]bool
	// highs counts the inputs of a conjunction whose last pulse was high.
	highs int
}

// network holds all the modules. names maps a module name to its index.
type network struct {
	// mods holds pointers, so press can change the state of a module in
	// place.
	mods        []*module
	names       map[string]int
	broadcaster int
}

// pulse is one pulse on its way from module from to module to.
type pulse struct {
	from, to int
	high     bool
}

// part1 pushes the button 1000 times and multiplies the number of low pulses
// by the number of high pulses.
func part1(in string) any {
	net := parse(in)
	lows, highs := 0, 0

	for range 1000 {
		net.press(func(p pulse) {
			if p.high {
				highs++
			} else {
				lows++
			}
		})
	}

	return lows * highs
}

// part2 finds the fewest button presses that send a low pulse to rx.
//
// A direct simulation takes far too long. This solution relies on the
// structure of the real input: one conjunction feeds rx, and several
// conjunctions feed that one. rx gets a low pulse only when all these
// feeders send a high pulse in the same press. Each feeder sends a high
// pulse on a fixed cycle, so we find the first press where each one sends
// a high pulse and take the LCM of these press counts.
//
// The examples have no rx module, so the result is nil for them.
func part2(in string) any {
	net := parse(in)

	rx, ok := net.names["rx"]
	if !ok {
		return nil
	}

	// Find the one conjunction that sends to rx, then all the modules that
	// send to that conjunction.
	hub := -1

	for i, m := range net.mods {
		for _, d := range m.dests {
			if d == rx {
				hub = i
			}
		}
	}

	if hub < 0 {
		return nil
	}

	cycles := map[int]int{}

	for from := range net.mods[hub].last {
		cycles[from] = 0
	}

	found := 0

	for presses := 1; found < len(cycles); presses++ {
		net.press(func(p pulse) {
			if p.to != hub || !p.high || cycles[p.from] != 0 {
				return
			}

			cycles[p.from] = presses
			found++
		})
	}

	var nums []int

	for _, n := range cycles {
		nums = append(nums, n)
	}

	return mathx.LCM(nums...)
}

// press pushes the button once and sends every pulse through the network.
// seen is called for each pulse, in the order the pulses are processed.
func (net *network) press(seen func(pulse)) {
	// The button sends a low pulse to the broadcaster. The button has no
	// module index, so from is -1.
	queue := []pulse{{from: -1, to: net.broadcaster}}

	// A slice with a moving head is a simple first-in, first-out queue.
	for head := 0; head < len(queue); head++ {
		p := queue[head]
		seen(p)

		m := net.mods[p.to]
		var out bool

		switch m.kind {
		case 'b':
			out = p.high
		case '%':
			if p.high {
				continue
			}

			m.on = !m.on
			out = m.on
		case '&':
			// Keep a count of the high inputs, so we do not scan all the
			// inputs on each pulse.
			if m.last[p.from] != p.high {
				if p.high {
					m.highs++
				} else {
					m.highs--
				}

				m.last[p.from] = p.high
			}

			out = m.highs != len(m.last)
		default:
			continue
		}

		for _, d := range m.dests {
			queue = append(queue, pulse{from: p.to, to: d, high: out})
		}
	}
}

// parse reads lines like "%a -> inv, con" into a network.
func parse(in string) *network {
	net := &network{names: map[string]int{}}

	// index gives the module index for a name, and adds a new module for a
	// name we have not seen yet.
	index := func(name string) int {
		if i, ok := net.names[name]; ok {
			return i
		}

		net.names[name] = len(net.mods)
		net.mods = append(net.mods, &module{last: map[int]bool{}})

		return len(net.mods) - 1
	}

	for _, line := range input.Lines(in) {
		left, right, _ := strings.Cut(line, " -> ")

		kind := byte('b')
		if left[0] == '%' || left[0] == '&' {
			kind = left[0]
			left = left[1:]
		}

		i := index(left)
		net.mods[i].kind = kind

		if kind == 'b' {
			net.broadcaster = i
		}

		for _, name := range strings.Split(right, ", ") {
			net.mods[i].dests = append(net.mods[i].dests, index(name))
		}
	}

	// A conjunction must know all of its inputs before the first pulse,
	// because it starts with a low pulse remembered for each one.
	for i, m := range net.mods {
		for _, d := range m.dests {
			net.mods[d].last[i] = false
		}
	}

	return net
}
