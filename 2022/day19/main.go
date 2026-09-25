// Advent of Code 2022, day 19: Not Enough Minerals.
// https://adventofcode.com/2022/day/19
//
// The puzzle text writes each example blueprint across several lines. The
// real input has one blueprint on each line, so example.txt uses that layout.
package main

import (
	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 19, part1, part2)
}

// The four materials. Each one is also the index of its robot type.
const (
	ore = iota
	clay
	obsidian
	geode
)

// blueprint holds the costs to build each robot type. cost[r][m] is the
// quantity of material m that robot type r needs.
type blueprint struct {
	id   int
	cost [4][3]int
	// limit[m] is the largest quantity of material m that one robot costs.
	// We can spend only one robot's cost each minute, so more robots than
	// this for material m give no benefit.
	limit [3]int
}

// part1 adds the quality level (id times most geodes in 24 minutes) of each
// blueprint.
func part1(in string) any {
	total := 0

	for _, bp := range parse(in) {
		total += bp.id * maxGeodes(bp, 24)
	}

	return total
}

// part2 multiplies the most geodes in 32 minutes for the first three
// blueprints. The example has only two blueprints, so it uses both.
func part2(in string) any {
	bps := parse(in)
	product := 1

	for _, bp := range bps[:min(3, len(bps))] {
		product *= maxGeodes(bp, 32)
	}

	return product
}

// state is one moment in the search. The geode count is the total that the
// built geode robots will crack by the end, so geode robots are not stored.
type state struct {
	timeLeft int
	robots   [3]int
	stock    [3]int
	geodes   int
}

// maxGeodes finds the most geodes the blueprint can crack in the given time.
//
// The search does not step one minute at a time. From each state, it picks
// the next robot to build and jumps ahead to the minute that robot is ready.
// This removes all the states where the robots only wait. Three rules prune
// the search:
//   - Do not build more robots of a material than the limit in the blueprint.
//   - When a geode robot is built, add all the geodes it will crack before
//     the end at once.
//   - Stop a branch when an optimistic estimate cannot beat the best result.
func maxGeodes(bp blueprint, minutes int) int {
	best := 0

	// dfs is a closure, so it can read bp and update best without passing
	// them on each call. A closure must be declared before it can call
	// itself, which is why it uses "var" first.
	var dfs func(s state)
	dfs = func(s state) {
		best = max(best, s.geodes)

		if upperBound(bp, s) <= best {
			return
		}

		// Try the geode robot first. Good results early make the prune
		// stronger.
		for r := geode; r >= ore; r-- {
			if r != geode && s.robots[r] >= bp.limit[r] {
				continue
			}

			wait, ok := waitFor(bp, s, r)
			if !ok {
				continue
			}

			// The robot is ready after the wait plus one minute to build it.
			// It must have at least one minute left to be useful.
			left := s.timeLeft - wait - 1
			if left <= 0 {
				continue
			}

			next := s
			next.timeLeft = left

			for m := range 3 {
				next.stock[m] += s.robots[m]*(wait+1) - bp.cost[r][m]
			}

			if r == geode {
				next.geodes += left
			} else {
				next.robots[r]++
			}

			dfs(next)
		}
	}

	dfs(state{timeLeft: minutes, robots: [3]int{ore: 1}})

	return best
}

// waitFor returns how many minutes to collect enough material to start robot
// r. It returns false when a needed material has no robot to collect it.
func waitFor(bp blueprint, s state, r int) (int, bool) {
	wait := 0

	for m := range 3 {
		need := bp.cost[r][m] - s.stock[m]
		if need <= 0 {
			continue
		}

		if s.robots[m] == 0 {
			return 0, false
		}

		wait = max(wait, mathx.DivCeil(need, s.robots[m]))
	}

	return wait, true
}

// upperBound returns a number of geodes that the state cannot do better than.
//
// It acts as if ore and clay are free. Then each minute we can build an
// obsidian robot, and also a geode robot when there is enough obsidian.
func upperBound(bp blueprint, s state) int {
	geodes := s.geodes
	obs := s.stock[obsidian]
	obsRobots := s.robots[obsidian]

	for t := s.timeLeft; t > 0; t-- {
		if obs >= bp.cost[geode][obsidian] {
			obs -= bp.cost[geode][obsidian]
			geodes += t - 1
		}

		obs += obsRobots
		obsRobots++
	}

	return geodes
}

// parse reads one blueprint from each line. The seven numbers on a line are
// the id and then the costs, in the order the text gives them.
func parse(in string) []blueprint {
	var bps []blueprint

	for _, line := range input.Lines(in) {
		n := input.Ints(line)

		bp := blueprint{id: n[0]}
		bp.cost[ore][ore] = n[1]
		bp.cost[clay][ore] = n[2]
		bp.cost[obsidian][ore] = n[3]
		bp.cost[obsidian][clay] = n[4]
		bp.cost[geode][ore] = n[5]
		bp.cost[geode][obsidian] = n[6]

		for r := range 4 {
			for m := range 3 {
				bp.limit[m] = max(bp.limit[m], bp.cost[r][m])
			}
		}

		bps = append(bps, bp)
	}

	return bps
}
