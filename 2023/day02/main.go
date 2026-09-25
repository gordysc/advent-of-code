// Advent of Code 2023, day 2: Cube Conundrum.
// https://adventofcode.com/2023/day/2
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 2, part1, part2)
}

// game holds the ID of one game and the largest number of cubes of each
// colour that the elf showed in one handful.
//
// Both parts need only these largest counts. A bag can give a handful only
// when it holds at least that many cubes of each colour, so the largest
// handful of each colour sets the limit.
type game struct {
	id               int
	red, green, blue int
}

// part1 adds the IDs of the games that are possible with a bag of 12 red,
// 13 green and 14 blue cubes.
func part1(in string) any {
	total := 0

	for _, g := range parseGames(in) {
		if g.red <= 12 && g.green <= 13 && g.blue <= 14 {
			total += g.id
		}
	}

	return total
}

// part2 adds the power of each game. The power is the product of the
// smallest number of cubes of each colour that makes the game possible.
func part2(in string) any {
	total := 0

	for _, g := range parseGames(in) {
		total += g.red * g.green * g.blue
	}

	return total
}

// parseGames reads each line, such as "Game 3: 8 green, 6 blue; 1 red", into
// a game with the largest count of each colour.
//
// The ";" between handfuls does not matter here, because the largest count
// over all handfuls is the same as the largest count over all "N colour"
// pairs. So the parser splits on both "," and ";".
func parseGames(in string) []game {
	var games []game

	for _, line := range input.Lines(in) {
		header, draws, _ := strings.Cut(line, ": ")
		g := game{id: input.Int(strings.TrimPrefix(header, "Game "))}

		// FieldsFunc splits on every rune that the function accepts, so it
		// handles both separators in one pass.
		pairs := strings.FieldsFunc(draws, func(r rune) bool { return r == ',' || r == ';' })

		for _, pair := range pairs {
			count, colour, _ := strings.Cut(strings.TrimSpace(pair), " ")
			n := input.Int(count)

			switch colour {
			case "red":
				g.red = max(g.red, n)
			case "green":
				g.green = max(g.green, n)
			case "blue":
				g.blue = max(g.blue, n)
			}
		}

		games = append(games, g)
	}

	return games
}
