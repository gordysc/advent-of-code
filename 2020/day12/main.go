// Advent of Code 2020, day 12: Rain Risk.
// https://adventofcode.com/2020/day/12
package main

import (
	"strconv"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 12, part1, part2)
}

// part1 moves the ship itself. N, S, E and W move the ship, L and R turn its
// heading, and F moves it along the heading. The ship starts facing east. The
// answer is the Manhattan distance from the start.
func part1(in string) any {
	var ship grid.Point
	heading := grid.Right

	for _, s := range parse(in) {
		switch s.action {
		case 'L', 'R':
			heading = turn(heading, s.action, s.value)
		case 'F':
			ship = ship.Add(heading.Scale(s.value))
		default:
			ship = ship.Add(compass(s.action).Scale(s.value))
		}
	}

	return ship.Manhattan(grid.Point{})
}

// part2 moves a waypoint that is relative to the ship. N, S, E and W move the
// waypoint, L and R rotate it around the ship, and F moves the ship to the
// waypoint value times. The waypoint starts 10 east and 1 north of the ship.
func part2(in string) any {
	var ship grid.Point
	waypoint := grid.P(10, -1)

	for _, s := range parse(in) {
		switch s.action {
		case 'L', 'R':
			waypoint = turn(waypoint, s.action, s.value)
		case 'F':
			ship = ship.Add(waypoint.Scale(s.value))
		default:
			waypoint = waypoint.Add(compass(s.action).Scale(s.value))
		}
	}

	return ship.Manhattan(grid.Point{})
}

// step is one navigation instruction: an action letter and its value.
type step struct {
	action byte
	value  int
}

// parse reads one instruction from each line, such as "F10" or "R90".
func parse(in string) []step {
	var steps []step
	for _, line := range input.Lines(in) {
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}

		steps = append(steps, step{line[0], n})
	}

	return steps
}

// turn rotates a vector around the origin in 90 degree steps. The same
// rotation works for the heading in part 1 and the waypoint in part 2.
func turn(v grid.Point, action byte, degrees int) grid.Point {
	for range degrees / 90 {
		if action == 'L' {
			v = v.TurnLeft()
		} else {
			v = v.TurnRight()
		}
	}

	return v
}

// compass gives the unit vector for a N, S, E or W action. Y grows downward
// in the grid package, so north is Up.
func compass(action byte) grid.Point {
	switch action {
	case 'N':
		return grid.Up
	case 'S':
		return grid.Down
	case 'E':
		return grid.Right
	case 'W':
		return grid.Left
	}

	panic("unknown action " + string(action))
}
