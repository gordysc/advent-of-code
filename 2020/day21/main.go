// Advent of Code 2020, day 21: Allergen Assessment.
// https://adventofcode.com/2020/day/21
package main

import (
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 21, part1, part2)
}

// food is one line of the list: its ingredients and the allergens it is
// known to contain.
type food struct {
	ingredients []string
	allergens   []string
}

// part1 counts how often the ingredients that cannot hold an allergen appear.
// An ingredient is safe when no allergen has it as a candidate.
func part1(in string) any {
	foods := parse(in)
	candidates := findCandidates(foods)

	unsafe := set.New[string]()
	for _, c := range candidates {
		for ing := range c {
			unsafe.Add(ing)
		}
	}

	count := 0
	for _, f := range foods {
		for _, ing := range f.ingredients {
			if !unsafe.Has(ing) {
				count++
			}
		}
	}

	return count
}

// part2 finds the one ingredient that holds each allergen. It repeatedly
// takes an allergen with one candidate left, fixes that match, and removes
// the ingredient from the other allergens. The answer lists the dangerous
// ingredients sorted by their allergen, joined with commas.
func part2(in string) any {
	candidates := findCandidates(parse(in))
	matched := map[string]string{}

	for len(matched) < len(candidates) {
		for allergen, c := range candidates {
			if _, done := matched[allergen]; done || c.Len() != 1 {
				continue
			}

			ing := c.Items()[0]
			matched[allergen] = ing

			for other, oc := range candidates {
				if other != allergen {
					oc.Remove(ing)
				}
			}
		}
	}

	allergens := make([]string, 0, len(matched))
	for allergen := range matched {
		allergens = append(allergens, allergen)
	}
	slices.Sort(allergens)

	dangerous := make([]string, len(allergens))
	for i, allergen := range allergens {
		dangerous[i] = matched[allergen]
	}

	return strings.Join(dangerous, ",")
}

// findCandidates maps each allergen to the ingredients that can hold it.
// An allergen is in exactly one ingredient, so that ingredient is in every
// food that lists the allergen. The candidates are the intersection of the
// ingredient lists of those foods.
func findCandidates(foods []food) map[string]set.Set[string] {
	candidates := map[string]set.Set[string]{}

	for _, f := range foods {
		ings := set.From(f.ingredients)

		for _, allergen := range f.allergens {
			if c, ok := candidates[allergen]; ok {
				candidates[allergen] = c.Intersect(ings)
			} else {
				candidates[allergen] = ings.Clone()
			}
		}
	}

	return candidates
}

// parse reads lines of the form "a b c (contains x, y)" into foods.
func parse(in string) []food {
	var foods []food

	for _, line := range input.Lines(in) {
		ings, allergens, _ := strings.Cut(line, " (contains ")
		allergens = strings.TrimSuffix(allergens, ")")

		foods = append(foods, food{
			ingredients: strings.Fields(ings),
			allergens:   strings.Split(allergens, ", "),
		})
	}

	return foods
}
