// Advent of Code 2019, day 21: Springdroid Adventure.
// https://adventofcode.com/2019/day/21
//
// The puzzle input is an Intcode program. The answer comes from the
// springscript programs below, which were written by hand, so the code only
// sends them and reads the result.
//
// The puzzle text has no example. example.txt holds a small handmade Intcode
// program that reads characters until it reads a 'K' (the end of WALK) or a
// 'U' (the middle of RUN), and then outputs 1000 plus the number of
// characters it read. No springscript instruction or register has those two
// letters, so it stops at the end of the script. Part 1 gives 1050 and part
// 2 gives 1079.
package main

import (
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 21, part1, part2)
}

// walkScript is the springscript for part 1. The sensors A, B, C and D say
// if there is ground 1, 2, 3 and 4 tiles ahead. A jump always lands 4 tiles
// ahead, on D. The script jumps when there is a hole in A, B or C and D is
// ground:
//
//	J = (!A || !B || !C) && D
//
// Jumping as soon as a hole comes into view gives the droid the most room,
// and D must be ground or the droid lands in a hole. T is a scratch register.
var walkScript = []string{
	"NOT A J",
	"NOT B T",
	"OR T J",
	"NOT C T",
	"OR T J",
	"AND D J",
	"WALK",
}

// runScript is the springscript for part 2. RUN adds the sensors E to I, for
// 5 to 9 tiles ahead. The part 1 rule can jump onto an island where the
// droid is then stuck: it lands on D, the next tile E is a hole, and a new
// jump from D would land on H, which is a hole too. So the script also needs
// E or H to be ground:
//
//	J = (!A || !B || !C) && D && (E || H)
//
// Springscript has no way to copy a register, so NOT E T then NOT T T puts
// E into T before H is added.
var runScript = []string{
	"NOT A J",
	"NOT B T",
	"OR T J",
	"NOT C T",
	"OR T J",
	"AND D J",
	"NOT E T",
	"NOT T T",
	"OR H T",
	"AND T J",
	"RUN",
}

// part1 walks the hull and returns the hull damage the droid reports.
func part1(in string) any {
	return survey(in, walkScript)
}

// part2 runs the hull with the longer sensor range and returns the damage.
func part2(in string) any {
	return survey(in, runScript)
}

// survey sends a springscript program to the droid, one instruction per line,
// and runs it. The program prints its prompts and pictures as ASCII, so every
// character is at most 127. If the droid gets across, the last output is the
// hull damage, a larger number. If it falls into a hole, the output is only
// an ASCII picture of the fall, and survey returns nil.
func survey(in string, script []string) any {
	m := intcode.New(intcode.Parse(in))
	for _, line := range script {
		m.SendLine(line)
	}

	out := m.Run()
	if len(out) == 0 || out[len(out)-1] <= 127 {
		return nil
	}

	return out[len(out)-1]
}
