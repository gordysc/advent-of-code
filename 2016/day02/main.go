// Advent of Code 2016, day 2: Bathroom Security.
// https://adventofcode.com/2016/day/2
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 2, part1, part2)
}

// square is the ordinary 3x3 keypad.
const square = `123
456
789`

// diamond is the keypad from part 2. Spaces are gaps you cannot move onto.
const diamond = `  1  
 234 
56789
 ABC 
  D  `

// part1 finds the bathroom code on a square keypad.
func part1(in string) any {
	return code(in, square)
}

// part2 finds the bathroom code on the diamond-shaped keypad.
func part2(in string) any {
	return code(in, diamond)
}

// code follows every line of moves on the keypad and returns the keys pressed.
// Each line starts where the previous one ended, beginning on the 5 key. A
// move that would leave the keypad, or land on a gap, is ignored.
func code(in, keypad string) string {
	keys := grid.Parse(keypad)
	pos, _ := keys.FindByte('5')

	var out strings.Builder

	for _, line := range input.Lines(in) {
		for _, r := range line {
			next := pos.Add(grid.DirFromRune[r])

			if key, ok := keys.Get(next); ok && key != ' ' {
				pos = next
			}
		}

		out.WriteByte(keys.At(pos))
	}

	return out.String()
}
