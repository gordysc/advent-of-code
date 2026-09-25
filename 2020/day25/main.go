// Advent of Code 2020, day 25: Combo Breaker.
// https://adventofcode.com/2020/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 25, part1, nil)
}

// modulus is the value the handshake divides by after each step, and
// subject is the subject number that makes the public keys.
const (
	modulus = 20201227
	subject = 7
)

// part1 finds the encryption key from the card and door public keys.
//
// A public key is subject^loop mod modulus. The loop size is small enough to
// find by trying each loop size in turn, so this finds the card's loop size.
// The encryption key is then the door's public key raised to that loop size.
func part1(in string) any {
	keys := input.Ints(in)
	cardKey, doorKey := keys[0], keys[1]

	loop := 0
	for v := 1; v != cardKey; v = v * subject % modulus {
		loop++
	}

	return modPow(doorKey, loop, modulus)
}

// modPow returns base^exp mod m. It squares the base for each bit of the
// exponent, so it needs only about log2(exp) multiplications.
func modPow(base, exp, m int) int {
	result := 1
	base %= m

	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % m
		}

		base = base * base % m
		exp >>= 1
	}

	return result
}
