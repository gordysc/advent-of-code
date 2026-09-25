// Advent of Code 2019, day 16: Flawed Frequency Transmission.
// https://adventofcode.com/2019/day/16
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 16, part1, part2)
}

// phases is how many times both parts run the FFT. The first small example in
// the puzzle text (12345678) shows only 4 phases, but every example with an
// answer to check uses 100.
const phases = 100

// repeats is how many copies of the input make the real signal in part 2.
const repeats = 10000

// offsetDigits is how many leading digits of the input give the message
// offset in part 2.
const offsetDigits = 7

// messageLen is how many digits both parts return.
const messageLen = 8

// part1 runs 100 phases of the FFT and returns the first eight digits.
func part1(in string) any {
	signal := input.Digits(in)

	for range phases {
		signal = phase(signal)
	}

	return digitString(signal[:messageLen])
}

// part2 repeats the input 10000 times, runs 100 phases and returns the eight
// digits at the offset that the first seven digits give.
//
// The fast method below only works when the offset is in the second half of
// the signal, which is true for every real input. The offset of example.txt
// is past the end of its signal, so part 2 has no answer for it.
func part2(in string) any {
	digits := input.Digits(in)

	offset := 0
	for _, d := range digits[:offsetDigits] {
		offset = offset*10 + d
	}

	total := len(digits) * repeats
	if offset < total/2 || offset+messageLen > total {
		return nil
	}

	// Only the digits from the offset to the end matter. The copies wrap
	// around the input, so position p of the long signal is digit p mod len.
	tail := make([]int, total-offset)
	for i := range tail {
		tail[i] = digits[(offset+i)%len(digits)]
	}

	// In the second half of the signal, output digit i has a pattern of 0 for
	// every position before i and 1 for every position from i to the end. So
	// each new digit is the sum of the old digits from i to the end, mod 10.
	// Walking from the end backwards builds those suffix sums in place: when
	// we reach i, tail[i+1] already holds its new value, which is the sum of
	// everything after i.
	for range phases {
		for i := len(tail) - 2; i >= 0; i-- {
			tail[i] = (tail[i] + tail[i+1]) % 10
		}
	}

	return digitString(tail[:messageLen])
}

// phase runs one phase of the FFT and returns the new signal.
//
// Output digit i (counting from 0) uses the base pattern 0, 1, 0, -1 with each
// value repeated i+1 times, shifted left by one. So the input digits alternate
// in runs of n = i+1: a run that adds, a run of zeros, a run that subtracts,
// and a run of zeros. The first adding run starts at n-1 and the pattern
// repeats every 4n digits.
//
// A prefix sum gives the sum of any run in one step, so digit i costs about
// len/(4n) steps. Added over all i, one phase costs about len*log(len)
// instead of len*len.
func phase(signal []int) []int {
	// prefix[k] is the sum of the first k digits, so the sum of the digits
	// from a up to (but not including) b is prefix[b] - prefix[a].
	prefix := make([]int, len(signal)+1)
	for k, d := range signal {
		prefix[k+1] = prefix[k] + d
	}

	// sum adds the digits from a up to b, clamped to the end of the signal.
	sum := func(a, b int) int {
		a, b = min(a, len(signal)), min(b, len(signal))

		return prefix[b] - prefix[a]
	}

	out := make([]int, len(signal))
	for i := range out {
		n := i + 1

		total := 0
		for start := n - 1; start < len(signal); start += 4 * n {
			total += sum(start, start+n)
			total -= sum(start+2*n, start+3*n)
		}

		out[i] = mathx.Abs(total) % 10
	}

	return out
}

// digitString joins digits into a string. The answer is a string and not a
// number, so leading zeros stay.
func digitString(digits []int) string {
	var sb strings.Builder
	for _, d := range digits {
		sb.WriteByte(byte('0' + d))
	}

	return sb.String()
}
