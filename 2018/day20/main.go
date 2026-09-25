// Advent of Code 2018, day 20: A Regular Map.
// https://adventofcode.com/2018/day/20
package main

import (
	"aoc/lib/ds"
	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 20, part1, part2)
}

// farDoors is how many doors away a room must be for part 2. Every room in
// the examples is much closer, so part 2 gives 0 for them.
const farDoors = 1000

// part1 returns the fewest doors needed to reach the farthest room.
func part1(in string) any {
	return maxValue(distances(in))
}

// part2 counts the rooms that need at least farDoors doors to reach.
func part2(in string) any {
	count := 0

	for _, d := range distances(in) {
		if d >= farDoors {
			count++
		}
	}

	return count
}

// distances walks the route regex and returns the fewest doors from the start
// to every room.
//
// Each letter moves one room in that direction. A "(" saves the current room
// on a stack, a "|" goes back to the saved room to try the next branch, and a
// ")" takes the saved room off the stack and goes on from there. That last
// rule is a shortcut: in these inputs, a group that has more route after it
// is a detour like (NEWS|), which ends back in the room it started from. Every
// step records the door count to the new room, and it keeps the smaller count
// when a room is reached by more than one path.
func distances(in string) map[grid.Point]int {
	here := grid.P(0, 0)
	dist := map[grid.Point]int{here: 0}

	var saved ds.Stack[grid.Point]

	for _, c := range in {
		switch c {
		case '(':
			saved.Push(here)
		case '|':
			here = saved.Peek()
		case ')':
			here = saved.Pop()
		case 'N', 'E', 'S', 'W':
			next := here.Add(grid.DirFromRune[c])

			if d, ok := dist[next]; !ok || dist[here]+1 < d {
				dist[next] = dist[here] + 1
			}

			here = next
		}
	}

	return dist
}

// maxValue returns the largest value in the map, or 0 for an empty map.
func maxValue(m map[grid.Point]int) int {
	best := 0

	for _, v := range m {
		best = max(best, v)
	}

	return best
}
