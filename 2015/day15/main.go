// Advent of Code 2015, day 15: Science for Hungry People.
// https://adventofcode.com/2015/day/15
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// teaspoons is the total amount of ingredients that a cookie must contain.
const teaspoons = 100

// mealCalories is the calorie count that a meal-replacement cookie in part 2
// must have exactly.
const mealCalories = 500

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 15, part1, part2)
}

// part1 finds the highest score of any cookie made from 100 teaspoons of
// ingredients.
func part1(in string) any {
	return bestScore(parseIngredients(in), func(calories int) bool { return true })
}

// part2 finds the highest score again, but only for cookies with exactly 500
// calories.
func part2(in string) any {
	return bestScore(parseIngredients(in), func(calories int) bool { return calories == mealCalories })
}

// ingredient holds the effect of one teaspoon of an ingredient. The four scored
// properties are capacity, durability, flavor and texture, in that order.
// Calories are kept apart because they do not count towards the score.
type ingredient struct {
	name     string
	scored   [4]int
	calories int
}

// parseIngredients reads lines like
// "Sugar: capacity 3, durability 0, flavor 0, texture -3, calories 2"
// into a slice of ingredients.
func parseIngredients(in string) []ingredient {
	var pantry []ingredient

	for _, line := range input.Lines(in) {
		// The name is the text before the colon. The five numbers always come in
		// the same order, and they can be negative, so Ints is the right helper.
		name, _, _ := strings.Cut(line, ":")
		nums := input.Ints(line)

		pantry = append(pantry, ingredient{
			name:     name,
			scored:   [4]int{nums[0], nums[1], nums[2], nums[3]},
			calories: nums[4],
		})
	}

	return pantry
}

// bestScore tries every way to divide the teaspoons between the ingredients
// and returns the highest score of a cookie that the accept function allows.
// The accept function gets the calorie count of the cookie.
func bestScore(pantry []ingredient, accept func(calories int) bool) int {
	best := 0

	// mix gives amounts to the ingredients from index i to the end. The totals
	// and the calories describe the part of the cookie that is already decided.
	var mix func(i, remaining int, totals [4]int, calories int)

	mix = func(i, remaining int, totals [4]int, calories int) {
		// The last ingredient has no choice: it must use all of the teaspoons
		// that remain. This removes one full level from the search.
		if i == len(pantry)-1 {
			totals = add(totals, pantry[i], remaining)
			calories += pantry[i].calories * remaining

			if accept(calories) {
				best = max(best, score(totals))
			}

			return
		}

		for amount := 0; amount <= remaining; amount++ {
			mix(i+1, remaining-amount, add(totals, pantry[i], amount), calories+pantry[i].calories*amount)
		}
	}

	if len(pantry) > 0 {
		mix(0, teaspoons, [4]int{}, 0)
	}

	return best
}

// add returns the totals after the given amount of an ingredient goes in.
// Arrays are values in Go, so the totals of the caller do not change.
func add(totals [4]int, ing ingredient, amount int) [4]int {
	for p, value := range ing.scored {
		totals[p] += value * amount
	}

	return totals
}

// score multiplies the four property totals together. A negative total counts
// as zero, which makes the full score zero.
func score(totals [4]int) int {
	product := 1

	for _, total := range totals {
		product *= max(total, 0)
	}

	return product
}
