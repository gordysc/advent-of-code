// Advent of Code 2017, day 4: High-Entropy Passphrases.
// https://adventofcode.com/2017/day/4
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 4, part1, part2)
}

// part1 counts the passphrases in which no word occurs twice.
func part1(in string) any {
	return countValid(in, func(word string) string { return word })
}

// part2 counts the passphrases in which no two words are anagrams of each
// other. Two words are anagrams when their sorted letters are the same.
func part2(in string) any {
	return countValid(in, func(word string) string {
		letters := []byte(word)
		slices.Sort(letters)

		return string(letters)
	})
}

// countValid counts the lines where no two words have the same key.
func countValid(in string, key func(string) string) int {
	count := 0

	for _, line := range input.Lines(in) {
		words := input.Fields(line)
		seen := set.New[string]()

		for _, word := range words {
			seen.Add(key(word))
		}

		if seen.Len() == len(words) {
			count++
		}
	}

	return count
}
