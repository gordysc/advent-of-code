// Advent of Code 2021, day 19: Beacon Scanner.
// https://adventofcode.com/2021/day/19
package main

import (
	"aoc/lib/ds"
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 19, part1, part2)
}

// part1 counts the beacons after all scanners are put on one map.
//
// Scanner 0 sets the coordinate system. Each scanner gets a fingerprint: the
// squared distances between all pairs of its beacons. Distances do not change
// when a scanner turns or moves. So two scanners that see the same 12 beacons
// share at least 66 distances (12 * 11 / 2). Only such pairs of scanners go to
// the slower alignment step, which tries all 24 orientations.
func part1(in string) any {
	beacons, _ := locate(parse(in))

	return beacons.Len()
}

// part2 finds the largest Manhattan distance between two scanners.
func part2(in string) any {
	_, positions := locate(parse(in))
	best := 0

	for i, a := range positions {
		for _, b := range positions[i+1:] {
			best = max(best, a.sub(b).manhattan())
		}
	}

	return best
}

// minOverlap is the number of beacons that two scanners must share.
const minOverlap = 12

// vec is a 3D position or offset.
type vec [3]int

// sub returns v - o.
func (v vec) sub(o vec) vec {
	return vec{v[0] - o[0], v[1] - o[1], v[2] - o[2]}
}

// add returns v + o.
func (v vec) add(o vec) vec {
	return vec{v[0] + o[0], v[1] + o[1], v[2] + o[2]}
}

// manhattan returns the sum of the absolute components.
func (v vec) manhattan() int {
	return mathx.Abs(v[0]) + mathx.Abs(v[1]) + mathx.Abs(v[2])
}

// matrix is a 3x3 rotation matrix. Each row has one 1 or -1 and two zeros.
type matrix [3][3]int

// apply rotates v with the matrix.
func (m matrix) apply(v vec) vec {
	var out vec

	for i := range 3 {
		out[i] = m[i][0]*v[0] + m[i][1]*v[1] + m[i][2]*v[2]
	}

	return out
}

// rotations holds the 24 orientations a scanner can have.
var rotations = makeRotations()

// makeRotations builds the 24 rotation matrices.
//
// Each output axis takes one input axis, with a sign. That gives 6 axis
// orders times 8 sign choices, which is 48 matrices. Half of them are
// mirror images (determinant -1), which a real rotation cannot make.
func makeRotations() []matrix {
	perms := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	var result []matrix

	for _, p := range perms {
		for signs := range 8 {
			var m matrix
			for row := range 3 {
				m[row][p[row]] = 1 - 2*(signs>>row&1)
			}

			if determinant(m) == 1 {
				result = append(result, m)
			}
		}
	}

	return result
}

// determinant calculates the determinant of a 3x3 matrix.
func determinant(m matrix) int {
	return m[0][0]*(m[1][1]*m[2][2]-m[1][2]*m[2][1]) -
		m[0][1]*(m[1][0]*m[2][2]-m[1][2]*m[2][0]) +
		m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0])
}

// scanner is the beacons that one scanner sees, relative to itself.
type scanner struct {
	beacons     []vec
	fingerprint map[int]int
}

// parse reads the scanner reports and makes the fingerprint of each scanner.
func parse(in string) []scanner {
	var scanners []scanner

	for _, block := range input.Blocks(in) {
		var s scanner

		// The first line of each block is the "--- scanner N ---" title.
		for _, line := range block[1:] {
			n := input.Ints(line)
			s.beacons = append(s.beacons, vec{n[0], n[1], n[2]})
		}

		s.fingerprint = fingerprint(s.beacons)
		scanners = append(scanners, s)
	}

	return scanners
}

// fingerprint counts the squared distances between all pairs of beacons.
// Squared distances stay integers, so they are safe to use as map keys.
func fingerprint(beacons []vec) map[int]int {
	counts := make(map[int]int)

	for i, a := range beacons {
		for _, b := range beacons[i+1:] {
			d := a.sub(b)
			counts[d[0]*d[0]+d[1]*d[1]+d[2]*d[2]]++
		}
	}

	return counts
}

// sharedDistances counts the distances that two fingerprints have in common.
func sharedDistances(a, b map[int]int) int {
	shared := 0

	for d, n := range a {
		shared += min(n, b[d])
	}

	return shared
}

// locate puts all scanners on the map of scanner 0. It returns the set of
// all beacons and the position of each scanner.
//
// A breadth-first search starts from scanner 0. Each scanner that has a
// known place is compared with the scanners that do not have one yet.
func locate(scanners []scanner) (set.Set[vec], []vec) {
	// placed[i] holds the beacons of scanner i in scanner 0 coordinates, or
	// nil if scanner i does not have a known place yet.
	placed := make([][]vec, len(scanners))
	positions := make([]vec, len(scanners))

	placed[0] = scanners[0].beacons
	queue := ds.NewQueue(0)
	minShared := minOverlap * (minOverlap - 1) / 2

	for !queue.Empty() {
		known := queue.Pop()

		for j, s := range scanners {
			if placed[j] != nil {
				continue
			}

			if sharedDistances(scanners[known].fingerprint, s.fingerprint) < minShared {
				continue
			}

			beacons, pos, ok := align(placed[known], s.beacons)
			if !ok {
				continue
			}

			placed[j] = beacons
			positions[j] = pos
			queue.Push(j)
		}
	}

	all := set.New[vec]()
	for _, beacons := range placed {
		all.Add(beacons...)
	}

	return all, positions
}

// align tries to put the beacons of an unknown scanner on the map of known
// beacons. It returns the moved beacons and the scanner position.
//
// For each orientation, each pair of a known and an unknown beacon suggests
// one offset between the two scanners. If one offset gets at least 12 votes,
// then 12 beacons agree and that orientation and offset are correct.
func align(known, unknown []vec) ([]vec, vec, bool) {
	for _, m := range rotations {
		rotated := make([]vec, len(unknown))
		for i, b := range unknown {
			rotated[i] = m.apply(b)
		}

		votes := make(map[vec]int)

		for _, a := range known {
			for _, b := range rotated {
				offset := a.sub(b)
				votes[offset]++

				if votes[offset] < minOverlap {
					continue
				}

				for i := range rotated {
					rotated[i] = rotated[i].add(offset)
				}

				return rotated, offset, true
			}
		}
	}

	return nil, vec{}, false
}
