// Advent of Code 2015, day 14: Reindeer Olympics.
// https://adventofcode.com/2015/day/14
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// raceSeconds is the length of the race in the puzzle. The worked example in
// the puzzle text stops after 1000 seconds instead, so the example answers
// from this program differ from the ones in the text.
const raceSeconds = 2503

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 14, part1, part2)
}

// part1 finds the distance of the reindeer that is furthest ahead when the
// race ends.
func part1(in string) any {
	return winningDistance(parseReindeer(in), raceSeconds)
}

// part2 scores the race by points instead: after every second, each reindeer
// in the lead gets one point. The answer is the highest score.
func part2(in string) any {
	return winningPoints(parseReindeer(in), raceSeconds)
}

// reindeer describes how one reindeer moves: it flies at speed km/s for
// flySeconds, then stands still for restSeconds, and repeats.
type reindeer struct {
	name        string
	speed       int
	flySeconds  int
	restSeconds int
}

// distanceAfter returns how far the reindeer is after the given number of
// seconds. No simulation is needed: every full fly-and-rest cycle covers the
// same distance, and in the part of a cycle that remains the reindeer flies
// until its flying time or the clock runs out, whichever comes first.
func (r reindeer) distanceAfter(seconds int) int {
	cycle := r.flySeconds + r.restSeconds
	flying := (seconds/cycle)*r.flySeconds + min(seconds%cycle, r.flySeconds)

	return flying * r.speed
}

// parseReindeer reads lines like
// "Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds."
// into a slice of reindeer.
func parseReindeer(in string) []reindeer {
	var herd []reindeer

	for _, line := range input.Lines(in) {
		// The name is the first word. The three numbers always come in the same
		// order, so UInts can pick them out without matching the words around them.
		name, _, _ := strings.Cut(line, " ")
		nums := input.UInts(line)

		herd = append(herd, reindeer{
			name:        name,
			speed:       nums[0],
			flySeconds:  nums[1],
			restSeconds: nums[2],
		})
	}

	return herd
}

// winningDistance returns the largest distance any reindeer has covered after
// the given number of seconds.
func winningDistance(herd []reindeer, seconds int) int {
	best := 0

	for _, r := range herd {
		best = max(best, r.distanceAfter(seconds))
	}

	return best
}

// winningPoints runs the race one second at a time and returns the highest
// score. The lead can change from second to second, so unlike part 1 every
// second has to be checked.
func winningPoints(herd []reindeer, seconds int) int {
	// points[i] is the score of herd[i].
	points := make([]int, len(herd))

	for t := 1; t <= seconds; t++ {
		lead := winningDistance(herd, t)

		// A tie for the lead gives a point to every reindeer in the tie, so this
		// loop cannot stop at the first match.
		for i, r := range herd {
			if r.distanceAfter(t) == lead {
				points[i]++
			}
		}
	}

	best := 0

	for _, p := range points {
		best = max(best, p)
	}

	return best
}
