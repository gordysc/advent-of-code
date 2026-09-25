// Advent of Code 2017, day 16: Permutation Promenade.
// https://adventofcode.com/2017/day/16
package main

import (
	"bytes"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 16, part1, part2)
}

// programs is the number of dancers in the puzzle, named a to p. The worked
// example in the puzzle text uses only five programs, a to e. Set this to 5 to
// get the example answers (baedc after one dance, ceadb after two).
const programs = 16

// dances is the number of full dances part 2 asks about.
const dances = 1_000_000_000

// move is one parsed dance move. The kind is 's' (spin), 'x' (exchange by
// position) or 'p' (partner, exchange by name). For a spin, a is the count.
// For an exchange, a and b are positions. For a partner, they are names.
type move struct {
	kind byte
	a, b int
}

// part1 returns the order of the programs after one dance.
func part1(in string) any {
	line := start()
	dance(line, parse(in))

	return string(line)
}

// part2 returns the order after one billion dances. The dance always maps the
// same order to the same next order, and it can be undone, so the orders
// repeat in a cycle that comes back to the start. Only the dances left after
// the last full cycle change the result.
func part2(in string) any {
	moves := parse(in)
	line := start()
	seen := []string{string(line)}

	for {
		dance(line, moves)
		if string(line) == seen[0] {
			break
		}

		seen = append(seen, string(line))
	}

	return seen[dances%len(seen)]
}

// start returns the programs in their first order: a, b, c, and so on.
func start() []byte {
	line := make([]byte, programs)
	for i := range line {
		line[i] = 'a' + byte(i)
	}

	return line
}

// parse turns the comma-separated moves into a list, so each dance does no
// string work.
func parse(in string) []move {
	var moves []move

	for _, s := range strings.Split(strings.TrimSpace(in), ",") {
		m := move{kind: s[0]}

		switch m.kind {
		case 's':
			m.a = input.Int(s[1:])
		case 'x':
			a, b, _ := strings.Cut(s[1:], "/")
			m.a, m.b = input.Int(a), input.Int(b)
		case 'p':
			m.a, m.b = int(s[1]), int(s[3])
		}

		moves = append(moves, m)
	}

	return moves
}

// dance does every move once, in place. A spin rotates the line to the right.
func dance(line []byte, moves []move) {
	n := len(line)
	spun := make([]byte, n)

	for _, m := range moves {
		switch m.kind {
		case 's':
			copy(spun, line[n-m.a:])
			copy(spun[m.a:], line[:n-m.a])
			copy(line, spun)
		case 'x':
			line[m.a], line[m.b] = line[m.b], line[m.a]
		case 'p':
			i := bytes.IndexByte(line, byte(m.a))
			j := bytes.IndexByte(line, byte(m.b))
			line[i], line[j] = line[j], line[i]
		}
	}
}
