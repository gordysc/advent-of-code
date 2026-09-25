// Advent of Code 2022, day 15: Beacon Exclusion Zone.
// https://adventofcode.com/2022/day/15
package main

import (
	"slices"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 15, part1, part2)
}

// part1 counts the positions on one row where a beacon cannot be.
//
// Each sensor covers a diamond. The diamond cuts the row in one interval, or
// not at all. We merge the intervals and add their lengths. Then we remove the
// known beacons on the row, because a beacon is there.
func part1(in string) any {
	sensors := parse(in)
	row, _ := limits(sensors)

	covered := 0
	for _, iv := range mergeIntervals(rowIntervals(sensors, row)) {
		covered += iv.hi - iv.lo + 1
	}

	beacons := set.New[int]()
	for _, s := range sensors {
		if s.beacon.Y == row {
			beacons.Add(s.beacon.X)
		}
	}

	return covered - beacons.Len()
}

// part2 finds the one position in the search square that no sensor covers,
// and returns its tuning frequency.
//
// The position is unique. So each cell next to it is covered, and it lies
// just outside the edge of some sensor diamonds. The diamond edges are
// diagonal lines of two types: x+y = a and x-y = b. The position is on the
// crossing of one line of each type, or where one line meets the edge of the
// square, or at a corner of the square. There are only a few thousand such
// points, so we test each one against all the sensors.
func part2(in string) any {
	sensors := parse(in)
	_, bound := limits(sensors)

	for _, p := range candidates(sensors, bound) {
		if p.X < 0 || p.Y < 0 || p.X > bound || p.Y > bound {
			continue
		}

		if !isCovered(sensors, p) {
			return p.X*4000000 + p.Y
		}
	}

	return nil
}

// sensor is a sensor, its closest beacon, and the distance between them.
type sensor struct {
	pos, beacon grid.Point
	radius      int
}

// interval is a range of x values from lo to hi, with both ends included.
type interval struct {
	lo, hi int
}

// parse reads one sensor and its beacon from each line.
func parse(in string) []sensor {
	var sensors []sensor

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		pos, beacon := grid.P(n[0], n[1]), grid.P(n[2], n[3])

		sensors = append(sensors, sensor{pos: pos, beacon: beacon, radius: pos.Manhattan(beacon)})
	}

	return sensors
}

// limits returns the row for part 1 and the size of the search square for
// part 2.
//
// The example uses row 10 and a square up to 20. The real input uses row
// 2000000 and a square up to 4000000. All example coordinates are less than
// 100, and real coordinates are in the millions. So small coordinates select
// the example values.
func limits(sensors []sensor) (row, bound int) {
	largest := 0
	for _, s := range sensors {
		largest = max(largest, s.pos.X, s.pos.Y, s.beacon.X, s.beacon.Y)
	}

	if largest < 100 {
		return 10, 20
	}

	return 2000000, 4000000
}

// rowIntervals returns the x interval that each sensor covers on the row.
// Sensors whose diamond does not reach the row give no interval.
func rowIntervals(sensors []sensor, row int) []interval {
	var ivs []interval

	for _, s := range sensors {
		half := s.radius - mathx.Abs(s.pos.Y-row)
		if half < 0 {
			continue
		}

		ivs = append(ivs, interval{s.pos.X - half, s.pos.X + half})
	}

	return ivs
}

// mergeIntervals sorts the intervals and joins the ones that overlap or
// touch.
func mergeIntervals(ivs []interval) []interval {
	slices.SortFunc(ivs, func(a, b interval) int { return a.lo - b.lo })

	var merged []interval
	for _, iv := range ivs {
		last := len(merged) - 1

		if last >= 0 && iv.lo <= merged[last].hi+1 {
			merged[last].hi = max(merged[last].hi, iv.hi)
			continue
		}

		merged = append(merged, iv)
	}

	return merged
}

// candidates returns the points where the missing beacon can be: each
// crossing of two diagonal lines just outside the sensor diamonds, each point
// where one such line meets the edge of the square, and the four corners.
// Some points are outside the square; the caller removes them.
func candidates(sensors []sensor, bound int) []grid.Point {
	// sums holds the a values of lines x+y = a. diffs holds the b values of
	// lines x-y = b.
	var sums, diffs []int
	for _, s := range sensors {
		d := s.radius + 1
		sums = append(sums, s.pos.X+s.pos.Y-d, s.pos.X+s.pos.Y+d)
		diffs = append(diffs, s.pos.X-s.pos.Y-d, s.pos.X-s.pos.Y+d)
	}

	points := []grid.Point{grid.P(0, 0), grid.P(bound, 0), grid.P(0, bound), grid.P(bound, bound)}

	for _, a := range sums {
		for _, b := range diffs {
			// The lines cross at x = (a+b)/2, y = (a-b)/2. The crossing is
			// on a whole cell only when a+b is even.
			if (a+b)%2 == 0 {
				points = append(points, grid.P((a+b)/2, (a-b)/2))
			}
		}
	}

	// Where the lines meet the four edges x=0, x=bound, y=0 and y=bound.
	for _, a := range sums {
		points = append(points, grid.P(0, a), grid.P(bound, a-bound), grid.P(a, 0), grid.P(a-bound, bound))
	}
	for _, b := range diffs {
		points = append(points, grid.P(0, -b), grid.P(bound, bound-b), grid.P(b, 0), grid.P(b+bound, bound))
	}

	return points
}

// isCovered reports whether some sensor is at least as close to p as to its
// own beacon.
func isCovered(sensors []sensor, p grid.Point) bool {
	for _, s := range sensors {
		if s.pos.Manhattan(p) <= s.radius {
			return true
		}
	}

	return false
}
