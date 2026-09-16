// Package input turns raw puzzle text into lines, numbers and blocks.
//
// Every function panics on malformed data instead of returning an error. Puzzle
// inputs are trusted and fixed, so a panic with a clear message is faster to
// work with than error plumbing in every solution.
package input

import (
	"regexp"
	"strconv"
	"strings"
)

// Lines splits the input into lines. A trailing newline does not produce an
// empty last line.
func Lines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

// Blocks splits the input into groups of lines separated by blank lines.
// Many puzzles use a blank line to separate records or sections.
func Blocks(s string) [][]string {
	var blocks [][]string

	for _, block := range strings.Split(strings.TrimRight(s, "\n"), "\n\n") {
		blocks = append(blocks, strings.Split(block, "\n"))
	}

	return blocks
}

// Fields splits a line on any run of whitespace, like strings.Fields.
func Fields(s string) []string {
	return strings.Fields(s)
}

// intPattern matches an optional minus sign followed by digits.
var intPattern = regexp.MustCompile(`-?\d+`)

// Ints returns every integer found anywhere in the text, in order. A minus sign
// directly before digits is treated as part of the number, so "3-4" gives 3 and -4.
// Use [UInts] when a dash is a separator rather than a sign.
func Ints(s string) []int {
	return atoiAll(intPattern.FindAllString(s, -1))
}

// uintPattern matches only digits.
var uintPattern = regexp.MustCompile(`\d+`)

// UInts returns every unsigned integer in the text, in order. Minus signs are
// ignored, so "3-4" gives 3 and 4.
func UInts(s string) []int {
	return atoiAll(uintPattern.FindAllString(s, -1))
}

// Int parses one integer and panics if the text is not a number.
func Int(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic("input.Int: " + err.Error())
	}

	return n
}

// IntLines parses every line as one integer.
func IntLines(s string) []int {
	return atoiAll(Lines(s))
}

// Digits returns each decimal digit character in the text as an int. Other
// characters are skipped.
func Digits(s string) []int {
	var out []int

	for _, r := range s {
		if r >= '0' && r <= '9' {
			out = append(out, int(r-'0'))
		}
	}

	return out
}

// atoiAll converts every string to an int, with a panic on the first failure.
func atoiAll(parts []string) []int {
	out := make([]int, 0, len(parts))

	for _, p := range parts {
		out = append(out, Int(p))
	}

	return out
}
