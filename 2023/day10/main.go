// Advent of Code 2023, day 10: Pipe Maze.
// https://adventofcode.com/2023/day/10
package main

import (
	"slices"

	"aoc/lib/grid"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 10, part1, part2)
}

// part1 finds the distance from S to the loop tile that is farthest from it.
// The two ways around the loop meet halfway, so the answer is half the loop
// length.
//
// example.txt holds the last part 2 example (the one that starts "FF7FSF7F7").
// The puzzle gives no part 1 answer for it. Part 1 gives 80.
func part1(in string) any {
	return len(findLoop(grid.Parse(in))) / 2
}

// part2 counts the tiles that the loop encloses.
//
// The shoelace formula gives the area A of the polygon through the centres of
// the loop tiles. Pick's theorem says A = i + b/2 - 1, where b is the number
// of tiles on the loop and i is the number of tiles inside it. So
// i = A - b/2 + 1.
//
// For example.txt, part 2 gives 10.
func part2(in string) any {
	loop := findLoop(grid.Parse(in))

	// The shoelace sum is twice the signed area. Its sign depends on the
	// direction of the walk, so take the absolute value.
	twice := 0
	for i, p := range loop {
		q := loop[(i+1)%len(loop)]
		twice += p.X*q.Y - q.X*p.Y
	}

	area := mathx.Abs(twice) / 2

	return area - len(loop)/2 + 1
}

// pipes gives the two directions that each pipe shape connects.
var pipes = map[byte][2]grid.Point{
	'|': {grid.Up, grid.Down},
	'-': {grid.Left, grid.Right},
	'L': {grid.Up, grid.Right},
	'J': {grid.Up, grid.Left},
	'7': {grid.Down, grid.Left},
	'F': {grid.Down, grid.Right},
}

// findLoop replaces S with its real pipe shape, then walks the loop from S.
// It returns the loop tiles in walk order, S first.
func findLoop(g grid.Grid[byte]) []grid.Point {
	start, _ := g.FindByte('S')
	g.Set(start, startShape(g, start))

	loop := []grid.Point{start}
	dir := pipes[g.At(start)][0]
	p := start.Add(dir)

	for p != start {
		loop = append(loop, p)

		// We came into this pipe through the end that points back along dir.
		// Leave it through the other end.
		ends := pipes[g.At(p)]
		if ends[0] == dir.Reverse() {
			dir = ends[1]
		} else {
			dir = ends[0]
		}

		p = p.Add(dir)
	}

	return loop
}

// startShape finds the pipe shape under S. S connects in the two directions
// whose neighbour pipe connects back to S.
func startShape(g grid.Grid[byte], start grid.Point) byte {
	var ends []grid.Point

	for _, d := range grid.Dirs4 {
		next, ok := g.Get(start.Add(d))
		if !ok {
			continue
		}

		// A missing key gives the zero value, so ground ('.') has no ends.
		back := pipes[next]
		if back[0] == d.Reverse() || back[1] == d.Reverse() {
			ends = append(ends, d)
		}
	}

	if len(ends) != 2 {
		panic("S does not connect to exactly two pipes")
	}

	for shape, pair := range pipes {
		if slices.Contains(ends, pair[0]) && slices.Contains(ends, pair[1]) {
			return shape
		}
	}

	panic("no pipe shape connects the two ends of S")
}
