// Advent of Code 2017, day 25: The Halting Problem.
// https://adventofcode.com/2017/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 25, part1, nil)
}

// action is what the machine does for one state and one current value: the
// value it writes, the way it moves (-1 for left, 1 for right), and the state
// it goes to next.
type action struct {
	write int
	move  int
	next  int
}

// machine is the Turing machine from the blueprint. States are numbered by
// their letter, so state A is 0. Each state has one action for a current
// value of 0 and one for a current value of 1.
type machine struct {
	start  int
	steps  int
	states [][2]action
}

// part1 runs the machine for the number of steps in the blueprint and
// returns the checksum: the number of 1s on the tape.
func part1(in string) any {
	m := parse(in)

	// The tape is a slice that grows when the cursor walks off either end.
	// pos is the cursor's index in the slice.
	tape := make([]int, 1)
	pos := 0
	state := m.start

	for range m.steps {
		act := m.states[state][tape[pos]]
		tape[pos] = act.write
		pos += act.move
		state = act.next

		switch {
		case pos < 0:
			// Double the tape and put the old cells in the right half.
			grown := make([]int, 2*len(tape))
			copy(grown[len(tape):], tape)
			pos += len(tape)
			tape = grown
		case pos == len(tape):
			tape = append(tape, 0)
		}
	}

	checksum := 0
	for _, v := range tape {
		checksum += v
	}

	return checksum
}

// parse reads the blueprint. The first block names the start state and the
// step count. Every later block is one state: a header line that names it,
// then two groups of four lines, one for each current value. Each line ends
// with the word that matters, so parse only reads the last word of a line.
func parse(in string) machine {
	blocks := input.Blocks(in)
	head := blocks[0]

	m := machine{
		start:  letter(lastWord(head[0])),
		steps:  input.Ints(head[1])[0],
		states: make([][2]action, len(blocks)-1),
	}

	for _, block := range blocks[1:] {
		var acts [2]action

		for v := range 2 {
			lines := block[2+4*v : 5+4*v]

			move := 1
			if lastWord(lines[1]) == "left" {
				move = -1
			}

			acts[v] = action{
				write: input.Int(lastWord(lines[0])),
				move:  move,
				next:  letter(lastWord(lines[2])),
			}
		}

		m.states[letter(lastWord(block[0]))] = acts
	}

	return m
}

// lastWord returns the last word of a line without its closing '.' or ':'.
func lastWord(line string) string {
	f := strings.Fields(line)

	return strings.TrimRight(f[len(f)-1], ".:")
}

// letter turns a state name such as "B" into its number.
func letter(name string) int {
	return int(name[0] - 'A')
}
