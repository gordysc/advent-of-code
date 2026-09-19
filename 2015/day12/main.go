// Advent of Code 2015, day 12: JSAbacusFramework.io.
// https://adventofcode.com/2015/day/12
package main

import (
	"encoding/json"
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 12, part1, part2)
}

// part1 adds up every number in the JSON document.
func part1(in string) any {
	return sum(parse(in), false)
}

// part2 adds up every number, but ignores each object that has "red" as one
// of its property values.
func part2(in string) any {
	return sum(parse(in), true)
}

// parse decodes the JSON text into a tree of plain Go values.
//
// When the target of json.Unmarshal is an empty interface (any), the package
// picks a Go type for each JSON type: an object becomes map[string]any, an
// array becomes []any, a number becomes float64, and a string becomes string.
// That lets the code walk a document without a struct that describes it.
func parse(in string) any {
	var doc any

	// The puzzle input is always valid JSON, so a decode error means the wrong
	// file was given. A panic stops the run with a clear message.
	if err := json.Unmarshal([]byte(strings.TrimSpace(in)), &doc); err != nil {
		panic(err)
	}

	return doc
}

// sum returns the total of every number at or below node. When skipRed is
// true, an object that has the string "red" as a property value counts as 0,
// together with everything inside it.
//
// The function calls itself for each child, so it reaches every depth of the
// document. The type switch below asks "which concrete type is inside this
// any value?" and binds v to the value with that type in each case.
func sum(node any, skipRed bool) int {
	switch v := node.(type) {
	case float64:
		// The JSON decoder gives all numbers as float64. The puzzle only has
		// whole numbers, so the conversion to int loses nothing.
		return int(v)

	case []any:
		total := 0

		// A "red" string inside an array does not cancel anything. Only
		// objects have that rule, so arrays are always summed in full.
		for _, child := range v {
			total += sum(child, skipRed)
		}

		return total

	case map[string]any:
		if skipRed && hasRed(v) {
			return 0
		}

		total := 0

		for _, child := range v {
			total += sum(child, skipRed)
		}

		return total
	}

	// Strings hold no numbers. The input has no booleans or nulls, but they
	// would also land here and count as 0.
	return 0
}

// hasRed reports whether one of the direct property values of obj is the
// string "red". Property names do not count, and neither do deeper values.
func hasRed(obj map[string]any) bool {
	for _, value := range obj {
		// A comparison between an any value and a string is only true when
		// the any holds a string with the same text. Numbers, arrays and
		// objects compare as not equal, so no type check is needed first.
		if value == "red" {
			return true
		}
	}

	return false
}
