// Advent of Code 2017, day 20: Particle Swarm.
// https://adventofcode.com/2017/day/20
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 20, part1, part2)
}

// ticks is how long both parts simulate. After this many ticks, the particle
// with the smallest acceleration is the closest one, and every collision has
// happened. The puzzle inputs settle well before that.
const ticks = 1000

// vec is a 3D position, velocity or acceleration.
type vec [3]int

// add returns the sum of two vectors.
func (v vec) add(w vec) vec {
	return vec{v[0] + w[0], v[1] + w[1], v[2] + w[2]}
}

// dist returns the Manhattan distance from v to the origin.
func (v vec) dist() int {
	return mathx.Abs(v[0]) + mathx.Abs(v[1]) + mathx.Abs(v[2])
}

// particle is one particle: its position, velocity and acceleration.
type particle struct {
	p, v, a vec
}

// move runs one tick: the acceleration changes the velocity first, and then
// the velocity changes the position.
func (pt *particle) move() {
	pt.v = pt.v.add(pt.a)
	pt.p = pt.p.add(pt.v)
}

// part1 returns the number of the particle that stays closest to the origin
// in the long run. It moves every particle for a long time and then picks the
// closest one.
func part1(in string) any {
	ps := parse(in)

	for range ticks {
		for i := range ps {
			ps[i].move()
		}
	}

	best := 0
	for i, pt := range ps {
		if pt.p.dist() < ps[best].p.dist() {
			best = i
		}
	}

	return best
}

// part2 counts the particles left after all collisions. After each tick, all
// particles that share a position are removed together.
//
// The example.txt file holds the part 1 example. The puzzle text has a
// separate example for part 2, which gives 1.
func part2(in string) any {
	ps := parse(in)

	for range ticks {
		count := map[vec]int{}
		for i := range ps {
			ps[i].move()
			count[ps[i].p]++
		}

		// Keep only the particles that are alone at their position. The new
		// slice reuses the memory of the old one, which is safe because it
		// never gets ahead of the loop that reads from it.
		left := ps[:0]
		for _, pt := range ps {
			if count[pt.p] == 1 {
				left = append(left, pt)
			}
		}

		ps = left
	}

	return len(ps)
}

// parse reads one particle per line. Each line has nine numbers, which can be
// negative: the position, the velocity and the acceleration.
func parse(in string) []particle {
	var ps []particle

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		ps = append(ps, particle{
			p: vec{n[0], n[1], n[2]},
			v: vec{n[3], n[4], n[5]},
			a: vec{n[6], n[7], n[8]},
		})
	}

	return ps
}
