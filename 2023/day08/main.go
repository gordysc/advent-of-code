// Advent of Code 2023, day 8: Haunted Wasteland.
// https://adventofcode.com/2023/day/8
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 8, part1, part2)
}

// part1 counts the steps from AAA to ZZZ. It gives nil when the map has no
// node AAA.
//
// example.txt holds the part 2 example (LR, with the nodes 11A and 22A). It
// has no node AAA, so part 1 gives nil for it.
func part1(in string) any {
	m := parse(in)

	if _, ok := m.nodes["AAA"]; !ok {
		return nil
	}

	return m.steps("AAA", func(node string) bool { return node == "ZZZ" })
}

// part2 counts the steps until all ghosts stand on a node that ends in Z at
// the same time. Each ghost starts on a node that ends in A.
//
// A step-by-step walk of all ghosts together takes too long. In the real
// input, each ghost goes into a cycle: it gets to its Z node after n steps,
// and then again after each n more steps. So all ghosts are on a Z node
// together after the least common multiple of their n values. The puzzle
// does not promise this, but all real inputs have it.
//
// For example.txt, part 2 gives 6.
func part2(in string) any {
	m := parse(in)

	var cycles []int
	for node := range m.nodes {
		if strings.HasSuffix(node, "A") {
			cycles = append(cycles, m.steps(node, func(node string) bool {
				return strings.HasSuffix(node, "Z")
			}))
		}
	}

	return mathx.LCM(cycles...)
}

// network holds the left/right instructions and the node map.
type network struct {
	turns string
	nodes map[string][2]string // nodes gives the left and right node of each node.
}

// parse reads the instruction line and the node lines.
func parse(in string) network {
	blocks := input.Blocks(in)
	m := network{turns: blocks[0][0], nodes: map[string][2]string{}}

	// A node line looks like "AAA = (BBB, CCC)". Replacing the punctuation
	// with spaces leaves three names that Fields can split.
	clean := strings.NewReplacer("=", " ", "(", " ", ")", " ", ",", " ")
	for _, line := range blocks[1] {
		f := strings.Fields(clean.Replace(line))
		m.nodes[f[0]] = [2]string{f[1], f[2]}
	}

	return m
}

// steps follows the instructions from start, repeating them as necessary. It
// counts the steps until it gets to a node where done is true.
func (m network) steps(start string, done func(string) bool) int {
	node := start
	count := 0

	for !done(node) {
		side := 0
		if m.turns[count%len(m.turns)] == 'R' {
			side = 1
		}

		node = m.nodes[node][side]
		count++
	}

	return count
}
