// Advent of Code 2020, day 4: Passport Processing.
// https://adventofcode.com/2020/day/4
package main

import (
	"regexp"
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// rules maps each required field to its part 2 check. The "cid" field is
// optional, so it has no rule and part 1 does not require it.
var rules = map[string]func(string) bool{
	"byr": func(v string) bool { return inRange(v, 1920, 2002) },
	"iyr": func(v string) bool { return inRange(v, 2010, 2020) },
	"eyr": func(v string) bool { return inRange(v, 2020, 2030) },
	"hgt": validHeight,
	"hcl": regexp.MustCompile(`^#[0-9a-f]{6}$`).MatchString,
	"ecl": regexp.MustCompile(`^(amb|blu|brn|gry|grn|hzl|oth)$`).MatchString,
	"pid": regexp.MustCompile(`^[0-9]{9}$`).MatchString,
}

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 4, part1, part2)
}

// part1 counts the passports that have all the required fields.
//
// example.txt holds the part 2 examples: four invalid passports, then four
// valid passports. All eight have the required fields, so part 1 gives 8 for
// it and not the 2 from the part 1 example.
func part1(in string) any {
	count := 0
	for _, p := range parse(in) {
		if hasFields(p) {
			count++
		}
	}

	return count
}

// part2 counts the passports that have all the required fields and where
// each field value passes its rule.
func part2(in string) any {
	count := 0
	for _, p := range parse(in) {
		if hasFields(p) && validValues(p) {
			count++
		}
	}

	return count
}

// hasFields reports whether the passport has every field in rules.
func hasFields(p map[string]string) bool {
	for field := range rules {
		if _, ok := p[field]; !ok {
			return false
		}
	}

	return true
}

// validValues reports whether every required field value passes its rule.
// It expects that the passport has all the required fields.
func validValues(p map[string]string) bool {
	for field, rule := range rules {
		if !rule(p[field]) {
			return false
		}
	}

	return true
}

// inRange reports whether v is a four-digit number from lo to hi, inclusive.
func inRange(v string, lo, hi int) bool {
	return len(v) == 4 && between(v, lo, hi)
}

// validHeight reports whether v is a number followed by "cm" (150 to 193)
// or "in" (59 to 76).
func validHeight(v string) bool {
	if num, ok := strings.CutSuffix(v, "cm"); ok {
		return between(num, 150, 193)
	}

	if num, ok := strings.CutSuffix(v, "in"); ok {
		return between(num, 59, 76)
	}

	return false
}

// between reports whether v is a number from lo to hi, inclusive.
func between(v string, lo, hi int) bool {
	n, err := strconv.Atoi(v)

	return err == nil && n >= lo && n <= hi
}

// parse reads each blank-line separated block into a map of field to value.
// A block can spread its "key:value" pairs over many lines.
func parse(in string) []map[string]string {
	var passports []map[string]string

	for _, block := range input.Blocks(in) {
		p := map[string]string{}

		for _, pair := range strings.Fields(strings.Join(block, " ")) {
			key, value, _ := strings.Cut(pair, ":")
			p[key] = value
		}

		passports = append(passports, p)
	}

	return passports
}
