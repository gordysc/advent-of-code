// Advent of Code 2016, day 21: Scrambled Letters and Hash.
// https://adventofcode.com/2016/day/21
package main

import (
	"bytes"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 21, part1, part2)
}

// password is the text part 1 scrambles, and scrambled is the text part 2
// unscrambles. The worked example in the puzzle text scrambles the 5-letter
// "abcde" into "decab" instead. On these 8-letter passwords the example steps
// give different strings, so the example output does not match "decab".
const (
	password  = "abcdefgh"
	scrambled = "fbgdceah"
)

// step is one scrambling operation. kind is the first two words of the line,
// and x and y are its two operands. Letters are stored as their byte value.
type step struct {
	kind string
	x, y int
}

// part1 scrambles the password with every step in order.
func part1(in string) any {
	s := []byte(password)

	for _, st := range parse(in) {
		apply(s, st)
	}

	return string(s)
}

// part2 finds the password that scrambles to the given text. It runs the
// steps in reverse order and undoes each one.
func part2(in string) any {
	s := []byte(scrambled)
	steps := parse(in)

	for i := len(steps) - 1; i >= 0; i-- {
		undo(s, steps[i])
	}

	return string(s)
}

// parse reads one step from every line. The operands sit at fixed word
// positions for each kind of step.
func parse(in string) []step {
	var steps []step

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)
		kind := f[0] + " " + f[1]

		var st step
		switch kind {
		case "swap position", "move position":
			st = step{kind, input.Int(f[2]), input.Int(f[5])}
		case "swap letter":
			st = step{kind, int(f[2][0]), int(f[5][0])}
		case "reverse positions":
			st = step{kind, input.Int(f[2]), input.Int(f[4])}
		case "rotate left", "rotate right":
			st = step{kind: kind, x: input.Int(f[2])}
		case "rotate based":
			st = step{kind: kind, x: int(f[6][0])}
		default:
			panic("unknown step: " + line)
		}

		steps = append(steps, st)
	}

	return steps
}

// apply runs one step on s in place.
func apply(s []byte, st step) {
	switch st.kind {
	case "swap position":
		s[st.x], s[st.y] = s[st.y], s[st.x]
	case "swap letter":
		i, j := bytes.IndexByte(s, byte(st.x)), bytes.IndexByte(s, byte(st.y))
		s[i], s[j] = s[j], s[i]
	case "reverse positions":
		slices.Reverse(s[st.x : st.y+1])
	case "rotate left":
		rotateLeft(s, st.x)
	case "rotate right":
		rotateLeft(s, -st.x)
	case "rotate based":
		i := bytes.IndexByte(s, byte(st.x))
		n := 1 + i
		if i >= 4 {
			n++
		}

		rotateLeft(s, -n)
	case "move position":
		move(s, st.x, st.y)
	}
}

// undo reverses one step on s in place. Swaps and reversals undo themselves,
// and rotations and moves undo by going the other way. The rotation based on
// a letter has no simple inverse, so undo tries every left rotation and keeps
// the one that the forward step turns back into s.
func undo(s []byte, st step) {
	switch st.kind {
	case "rotate left":
		rotateLeft(s, -st.x)
	case "rotate right":
		rotateLeft(s, st.x)
	case "move position":
		move(s, st.y, st.x)
	case "rotate based":
		want := slices.Clone(s)
		try := make([]byte, len(s))

		for n := range len(s) {
			copy(try, want)
			rotateLeft(try, n)

			before := slices.Clone(try)
			apply(try, st)

			if bytes.Equal(try, want) {
				copy(s, before)
				return
			}
		}

		panic("no rotation undoes the step")
	default:
		apply(s, st)
	}
}

// rotateLeft rotates s left by n places in place. A negative n rotates right.
func rotateLeft(s []byte, n int) {
	n = mathx.Mod(n, len(s))
	rotated := append(slices.Clone(s[n:]), s[:n]...)
	copy(s, rotated)
}

// move takes the letter at position from out of s and puts it back in at
// position to.
func move(s []byte, from, to int) {
	c := s[from]
	rest := slices.Delete(slices.Clone(s), from, from+1)
	copy(s, slices.Insert(rest, to, c))
}
