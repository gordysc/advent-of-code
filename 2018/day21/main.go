// Advent of Code 2018, day 21: Chronal Conversion.
// https://adventofcode.com/2018/day/21
//
// The puzzle text has no example. example.txt holds a sample program with the
// same shape as the real inputs, which differ only in two constants: the seed
// and the multiplier of the hash.
//
// The program only stops when register 0 equals a value that it computes in a
// loop. That value comes from a small hash: the program seeds it, mixes in
// the bytes of a helper number one at a time, and compares the result with
// register 0. It then feeds the result back in as the next helper number.
//
// Running the program on the device is too slow for part 2, because it finds
// each byte of the helper number with a slow loop that counts up to a
// division by 256. So this solution reads the seed and the multiplier from
// the input and runs the same hash directly in Go.
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 21, part1, part2)
}

// The constants that every input shares. The helper number always starts as
// the last hash with bit 16 set, the hash keeps 24 bits, and the helper
// number gives up one byte at a time.
const (
	helperBit = 65536
	hashMask  = 16777215
	byteSize  = 256
)

// part1 returns the first value the program compares with register 0. With
// that value in register 0, the program stops after the fewest instructions.
func part1(in string) any {
	seed, mult := constants(in)

	return next(0, seed, mult)
}

// part2 returns the last new value the program compares with register 0.
// The hash only has 2^24 values, so the sequence comes back to a value it
// already made and then repeats. The last new value before that makes the
// program run for the most instructions and still stop.
func part2(in string) any {
	seed, mult := constants(in)
	seen := map[int]bool{}

	last := 0
	for x := next(0, seed, mult); !seen[x]; x = next(x, seed, mult) {
		seen[x] = true
		last = x
	}

	return last
}

// constants finds the seed and the multiplier in the program. The seed is
// set by the line right after "bori X 65536 Y", which builds the helper
// number. The multiplier is the "muli" that works on the hash register X,
// which is the register that "eqrr" compares with register 0.
func constants(in string) (seed, mult int) {
	lines := input.Lines(in)

	hash := ""
	for _, line := range lines {
		f := strings.Fields(line)
		if f[0] == "eqrr" {
			hash = f[1]
			if hash == "0" {
				hash = f[2]
			}
		}
	}

	for i, line := range lines {
		f := strings.Fields(line)

		switch {
		case f[0] == "bori" && f[2] == "65536":
			seed = input.Int(strings.Fields(lines[i+1])[1])
		case f[0] == "muli" && f[1] == hash:
			mult = input.Int(f[2])
		}
	}

	return seed, mult
}

// next runs one pass of the program's outer loop. It takes the last value
// the program compared with register 0 and returns the next one.
func next(last, seed, mult int) int {
	helper := last | helperBit
	x := seed

	for {
		x = (x + (helper & (byteSize - 1))) & hashMask
		x = (x * mult) & hashMask

		if helper < byteSize {
			return x
		}

		helper /= byteSize
	}
}
