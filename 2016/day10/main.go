// Advent of Code 2016, day 10: Balance Bots.
// https://adventofcode.com/2016/day/10
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 10, part1, part2)
}

// lowChip and highChip are the pair of chips part 1 asks about. The worked
// example in the puzzle text asks about chips 2 and 5 instead, so part 1 has
// no answer for the example.
const (
	lowChip  = 17
	highChip = 61
)

// target is where a bot sends a chip: another bot, or an output bin.
type target struct {
	output bool
	id     int
}

// rule is the pair of targets a bot sends its low and high chips to.
type rule struct {
	low, high target
}

// factory is the state after every chip has moved as far as it can.
type factory struct {
	compared map[[2]int]int // bot that compared each (low, high) pair of chips
	outputs  map[int]int    // chip value in each output bin
}

// part1 finds the bot that compares the two chips from the puzzle.
func part1(in string) any {
	bot, ok := run(in).compared[[2]int{lowChip, highChip}]
	if !ok {
		return nil
	}

	return bot
}

// part2 multiplies the chips that end up in outputs 0, 1 and 2.
func part2(in string) any {
	outputs := run(in).outputs

	return outputs[0] * outputs[1] * outputs[2]
}

// run reads the instructions and moves the chips until no bot holds two.
// A bot acts only when it holds two chips, so the bots that are ready wait in
// a queue. Each one hands off its chips, which can make other bots ready.
func run(in string) factory {
	hands := map[int][]int{}
	rules := map[int]rule{}

	for _, line := range input.Lines(in) {
		f := strings.Fields(line)

		if f[0] == "value" {
			bot := input.Int(f[5])
			hands[bot] = append(hands[bot], input.Int(f[1]))
			continue
		}

		rules[input.Int(f[1])] = rule{
			low:  target{output: f[5] == "output", id: input.Int(f[6])},
			high: target{output: f[10] == "output", id: input.Int(f[11])},
		}
	}

	var ready []int
	for bot, chips := range hands {
		if len(chips) == 2 {
			ready = append(ready, bot)
		}
	}

	out := factory{compared: map[[2]int]int{}, outputs: map[int]int{}}

	// give puts a chip in an output bin, or in a bot's hands. A bot that now
	// holds two chips joins the queue.
	give := func(to target, chip int) {
		if to.output {
			out.outputs[to.id] = chip
			return
		}

		hands[to.id] = append(hands[to.id], chip)
		if len(hands[to.id]) == 2 {
			ready = append(ready, to.id)
		}
	}

	for len(ready) > 0 {
		bot := ready[0]
		ready = ready[1:]

		a, b := hands[bot][0], hands[bot][1]
		low, high := min(a, b), max(a, b)
		hands[bot] = nil
		out.compared[[2]int{low, high}] = bot

		give(rules[bot].low, low)
		give(rules[bot].high, high)
	}

	return out
}
