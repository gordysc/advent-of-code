// Advent of Code 2019, day 12: The N-Body Problem.
// https://adventofcode.com/2019/day/12
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 12, part1, part2)
}

// steps is how long part 1 runs the simulation. The worked examples in the
// puzzle text stop after 10 steps (total energy 179) and after 100 steps
// (total energy 1940) instead, so part 1 gives a different number for them.
const steps = 1000

// axis holds one coordinate (x, y or z) of every moon: pos[i] and vel[i] are
// the position and velocity of moon i along that axis.
type axis struct {
	pos, vel []int
}

// part1 runs the simulation and adds up the total energy of every moon.
func part1(in string) any {
	axes := parse(in)

	for range steps {
		for i := range axes {
			axes[i].step()
		}
	}

	total := 0

	for m := range axes[0].pos {
		potential, kinetic := 0, 0

		for _, a := range axes {
			potential += mathx.Abs(a.pos[m])
			kinetic += mathx.Abs(a.vel[m])
		}

		total += potential * kinetic
	}

	return total
}

// part2 finds how many steps pass before the moons return to an earlier
// state. Gravity along x only depends on x positions, and so on, so each
// axis repeats on its own cycle. The whole system repeats when all three
// cycles line up, which is their least common multiple.
func part2(in string) any {
	axes := parse(in)
	cycles := make([]int, len(axes))

	for i := range axes {
		cycles[i] = axes[i].period()
	}

	return mathx.LCM(cycles...)
}

// parse reads one moon per line and splits the positions into three axes.
// Every moon starts with zero velocity.
func parse(in string) [3]axis {
	var axes [3]axis

	for _, line := range input.Lines(in) {
		nums := input.Ints(line)

		for i := range axes {
			axes[i].pos = append(axes[i].pos, nums[i])
			axes[i].vel = append(axes[i].vel, 0)
		}
	}

	return axes
}

// step moves the moons one time step along this axis. First every pair of
// moons pulls each other one unit closer, then each moon moves by its
// velocity. The method has a pointer receiver, but pos and vel are slices,
// so even a copy of the struct would change the same backing arrays.
func (a *axis) step() {
	for i := range a.pos {
		for j := i + 1; j < len(a.pos); j++ {
			pull := mathx.Sign(a.pos[j] - a.pos[i])
			a.vel[i] += pull
			a.vel[j] -= pull
		}
	}

	for i := range a.pos {
		a.pos[i] += a.vel[i]
	}
}

// period counts the steps until this axis returns to its starting state.
// Each step can be run backwards (the velocity change depends only on the
// positions), so two different states never lead to the same next state.
// That means the first state to repeat must be the starting state, and
// there is no need to remember any of the states in between.
func (a *axis) period() int {
	// Every moon starts at rest, so the starting velocities are all zero.
	start := slices.Clone(a.pos)
	still := make([]int, len(a.vel))
	n := 0

	for {
		a.step()
		n++

		if slices.Equal(a.pos, start) && slices.Equal(a.vel, still) {
			return n
		}
	}
}
