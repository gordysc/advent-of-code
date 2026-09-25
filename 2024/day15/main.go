// Advent of Code 2024, day 15: Warehouse Woes.
// https://adventofcode.com/2024/day/15
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 15, part1, part2)
}

// part1 moves the robot through the warehouse and adds the GPS coordinates
// of the boxes.
//
// The puzzle has a small and a larger example. The file example.txt holds
// the larger one. On it, part 1 gives 10092 and part 2 gives 9021. On the
// small example, part 1 gives 2028.
func part1(in string) any {
	return simulate(in, false)
}

// part2 does the same in a warehouse where everything except the robot is
// twice as wide.
func part2(in string) any {
	return simulate(in, true)
}

// simulate reads the warehouse and the moves, runs every move, and returns
// the sum of the GPS coordinates of the boxes. When wide is true, it first
// makes the map twice as wide.
func simulate(in string, wide bool) int {
	blocks := input.Blocks(in)
	if len(blocks) < 2 {
		return 0
	}

	rows := blocks[0]
	if wide {
		rows = widen(rows)
	}

	g := grid.Parse(strings.Join(rows, "\n"))

	robot, ok := g.FindByte('@')
	if !ok {
		return 0
	}

	// The moves span many lines, but the line breaks have no meaning.
	for _, r := range strings.Join(blocks[1], "") {
		d, ok := grid.DirFromRune[r]
		if !ok {
			continue
		}

		if push(g, robot, d) {
			robot = robot.Add(d)
		}
	}

	total := 0
	for p, c := range g.All() {
		// A wide box measures from its left half, "[".
		if c == 'O' || c == '[' {
			total += 100*p.Y + p.X
		}
	}

	return total
}

// widen doubles each tile of the map, as part 2 tells us to.
func widen(rows []string) []string {
	r := strings.NewReplacer("#", "##", "O", "[]", ".", "..", "@", "@.")

	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = r.Replace(row)
	}

	return out
}

// push tries to move the robot at start one step in direction d. It moves
// the robot and every box in the way, and returns true. If a wall stops any
// of them, nothing moves and it returns false.
//
// A breadth-first search collects every tile that must move. Each tile pushes
// the tile in front of it. A wide box half also pulls its other half, so one
// wide box can push two boxes when it moves up or down. For a sideways move,
// the other half is already in the line of tiles, so adding it has no effect.
func push(g grid.Grid[byte], start, d grid.Point) bool {
	moving := []grid.Point{start}
	seen := map[grid.Point]bool{start: true}

	add := func(p grid.Point) {
		if !seen[p] {
			seen[p] = true
			moving = append(moving, p)
		}
	}

	// moving grows while we loop over it, so use an index, not range.
	for i := 0; i < len(moving); i++ {
		next := moving[i].Add(d)

		switch g.At(next) {
		case '#':
			return false
		case 'O':
			add(next)
		case '[':
			add(next)
			add(next.Add(grid.Right))
		case ']':
			add(next)
			add(next.Add(grid.Left))
		}
	}

	// Read every tile before we write any. A tile can land on the old place
	// of another tile, so writes in place could overwrite a tile we still
	// need to move.
	tiles := make([]byte, len(moving))
	for i, p := range moving {
		tiles[i] = g.At(p)
		g.Set(p, '.')
	}

	for i, p := range moving {
		g.Set(p.Add(d), tiles[i])
	}

	return true
}
