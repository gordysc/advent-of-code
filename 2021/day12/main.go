// Advent of Code 2021, day 12: Passage Pathing.
// https://adventofcode.com/2021/day/12
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 12, part1, part2)
}

// part1 counts the paths from start to end that visit each small cave at most
// one time.
func part1(in string) any {
	return parse(in).paths(false)
}

// part2 counts the paths when one small cave can be visited two times.
// The start and end caves still get one visit only.
func part2(in string) any {
	return parse(in).paths(true)
}

// caveMap holds the caves as numbers so that a set of visited caves fits in
// the bits of one integer.
type caveMap struct {
	start, end int
	small      []bool
	links      [][]int
}

// parse reads the tunnel list. It gives each cave name a number the first time
// it sees that name.
func parse(in string) caveMap {
	var m caveMap
	ids := map[string]int{}

	id := func(name string) int {
		if n, ok := ids[name]; ok {
			return n
		}

		n := len(ids)
		ids[name] = n
		m.small = append(m.small, name == strings.ToLower(name))
		m.links = append(m.links, nil)

		return n
	}

	for _, line := range input.Lines(in) {
		a, b, _ := strings.Cut(line, "-")
		x, y := id(a), id(b)
		m.links[x] = append(m.links[x], y)
		m.links[y] = append(m.links[y], x)
	}

	m.start, m.end = ids["start"], ids["end"]

	return m
}

// paths counts the paths from start to end. When spare is true, the path can
// visit one small cave (not start) a second time.
func (m caveMap) paths(spare bool) int {
	return m.walk(m.start, 1<<m.start, spare)
}

// walk counts the paths from cave at to the end. seen has one bit for each
// small cave on the path so far. spare tells if the second visit is still
// available.
func (m caveMap) walk(at int, seen uint64, spare bool) int {
	if at == m.end {
		return 1
	}

	count := 0

	for _, next := range m.links[at] {
		switch {
		case next == m.start:
			continue
		case !m.small[next]:
			count += m.walk(next, seen, spare)
		case seen&(1<<next) == 0:
			count += m.walk(next, seen|1<<next, spare)
		case spare:
			count += m.walk(next, seen, false)
		}
	}

	return count
}
