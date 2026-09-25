// Advent of Code 2019, day 10: Monitoring Station.
// https://adventofcode.com/2019/day/10
package main

import (
	"cmp"
	"math"
	"slices"

	"aoc/lib/grid"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 10, part1, part2)
}

// betNumber is the asteroid that part 2 asks about: the 200th to be
// vaporized. The smaller examples in the puzzle text have fewer than 200
// other asteroids, so part 2 has no answer for them. example.txt is the large
// example, where the 200th asteroid is at 8,2.
const betNumber = 200

// asteroid is the character that marks an asteroid on the map.
const asteroid = '#'

// part1 counts the asteroids that the best station can see.
func part1(in string) any {
	_, groups := bestStation(in)

	return len(groups)
}

// part2 turns the laser clockwise from straight up, starting at the best
// station, and returns 100*x+y for the 200th asteroid it vaporizes.
//
// The laser hits only the nearest asteroid in each direction on one turn, so
// the asteroid that is k-th nearest in its direction goes on turn k. Sorting
// by turn, and then by angle inside a turn, gives the order they are hit.
func part2(in string) any {
	station, groups := bestStation(in)

	type target struct {
		p     grid.Point
		turn  int
		angle float64
	}

	var targets []target

	for dir, group := range groups {
		// All asteroids in a group are on the same ray, so the Manhattan
		// distance puts them in the same order as the true distance.
		slices.SortFunc(group, func(a, b grid.Point) int {
			return cmp.Compare(station.Manhattan(a), station.Manhattan(b))
		})

		angle := clockwise(dir)
		for turn, p := range group {
			targets = append(targets, target{p: p, turn: turn, angle: angle})
		}
	}

	if len(targets) < betNumber {
		return nil
	}

	slices.SortFunc(targets, func(a, b target) int {
		return cmp.Or(cmp.Compare(a.turn, b.turn), cmp.Compare(a.angle, b.angle))
	})

	p := targets[betNumber-1].p

	return 100*p.X + p.Y
}

// bestStation finds the asteroid that can see the most other asteroids. It
// returns that asteroid and the other asteroids grouped by the direction
// they are in from it. Each group holds the asteroids on one line of sight,
// so the number of groups is the number of asteroids the station can see.
func bestStation(in string) (grid.Point, map[grid.Point][]grid.Point) {
	var asteroids []grid.Point
	for p, b := range grid.Parse(in).All() {
		if b == asteroid {
			asteroids = append(asteroids, p)
		}
	}

	var station grid.Point
	var best map[grid.Point][]grid.Point

	for _, s := range asteroids {
		groups := sightLines(s, asteroids)
		if len(groups) > len(best) {
			station, best = s, groups
		}
	}

	return station, best
}

// sightLines groups the asteroids by their direction from station. Dividing
// the offset by the GCD of its two parts gives the smallest step in that
// direction, so (2, 4) and (3, 6) both become (1, 2) and share one group.
func sightLines(station grid.Point, asteroids []grid.Point) map[grid.Point][]grid.Point {
	groups := map[grid.Point][]grid.Point{}

	for _, p := range asteroids {
		if p == station {
			continue
		}

		d := p.Sub(station)
		g := mathx.GCD(d.X, d.Y)
		dir := grid.P(d.X/g, d.Y/g)
		groups[dir] = append(groups[dir], p)
	}

	return groups
}

// clockwise returns the angle of a direction in radians, measured clockwise
// from straight up, from 0 up to but not including 2π.
//
// math.Atan2(y, x) measures counterclockwise from the positive x axis. On
// the map Y grows downward, so up is (0, -1). Passing dx as the first
// argument and -dy as the second swaps the axes so that up gives 0 and right
// gives π/2. Directions to the left give negative angles, so adding 2π moves
// them after the rest of the turn.
func clockwise(d grid.Point) float64 {
	angle := math.Atan2(float64(d.X), float64(-d.Y))
	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle
}
