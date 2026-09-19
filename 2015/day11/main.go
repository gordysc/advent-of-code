// Advent of Code 2015, day 11: Corporate Policy.
// https://adventofcode.com/2015/day/11
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 11, part1, part2)
}

// part1 finds the next password after the input that passes every rule.
func part1(in string) any {
	return next(strings.TrimSpace(in))
}

// part2 finds the password after the one from part 1.
func part2(in string) any {
	return next(next(strings.TrimSpace(in)))
}

// next returns the first valid password that comes after old. It works on a
// byte slice so each step can change letters in place instead of building new
// strings.
func next(old string) string {
	pw := []byte(old)

	skipForbidden(pw)

	for {
		increment(pw)

		if valid(pw) {
			return string(pw)
		}
	}
}

// forbidden reports whether c is one of the letters the policy bans because
// they look like other characters.
func forbidden(c byte) bool {
	return c == 'i' || c == 'o' || c == 'l'
}

// skipForbidden jumps past every password that contains a banned letter.
//
// If position i holds a banned letter then every password with the same prefix
// up to i is also invalid, so the fastest move is to bump that letter and reset
// everything after it to 'a'. A banned letter is never the last letter of the
// alphabet, so the bump never carries.
func skipForbidden(pw []byte) {
	for i, c := range pw {
		if !forbidden(c) {
			continue
		}

		pw[i]++

		for j := i + 1; j < len(pw); j++ {
			pw[j] = 'a'
		}

		return
	}
}

// increment moves pw to the next candidate, like the odometer in a car: the
// last letter goes up by one, and a 'z' wraps to 'a' and carries into the
// letter to its left.
//
// When a letter lands on a banned value the function moves straight past it,
// which skips 26^k candidates at once for a letter k places from the end.
func increment(pw []byte) {
	for i := len(pw) - 1; i >= 0; i-- {
		if pw[i] == 'z' {
			pw[i] = 'a'
			continue
		}

		pw[i]++

		if forbidden(pw[i]) {
			pw[i]++
		}

		return
	}
}

// valid reports whether pw meets the two rules that increment does not already
// guarantee: a straight of three consecutive letters, and two different pairs.
// The banned-letter rule is enforced by skipForbidden and increment, so it does
// not need a check here.
func valid(pw []byte) bool {
	return hasStraight(pw) && hasTwoPairs(pw)
}

// hasStraight looks for three letters in a row that go up by one each, such as
// "abc" or "xyz".
func hasStraight(pw []byte) bool {
	for i := 0; i+2 < len(pw); i++ {
		if pw[i+1] == pw[i]+1 && pw[i+2] == pw[i]+2 {
			return true
		}
	}

	return false
}

// hasTwoPairs looks for two doubled letters that use different letters, such
// as "aa" and "zz". A run like "aaa" counts as one pair, so after a match the
// scan skips the second letter of the pair.
func hasTwoPairs(pw []byte) bool {
	var first byte
	found := 0

	for i := 0; i+1 < len(pw); i++ {
		if pw[i] != pw[i+1] || pw[i] == first {
			continue
		}

		found++
		if found == 2 {
			return true
		}

		first = pw[i]
		i++
	}

	return false
}
