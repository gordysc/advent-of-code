// Advent of Code 2016, day 17: Two Steps Forward.
// https://adventofcode.com/2016/day/17
package main

import (
	"crypto/md5"
	"strings"

	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 17, part1, part2)
}

// size is the width and height of the grid of rooms. The start is the top-left
// room and the vault is the bottom-right room.
const size = 4

// move is one door out of a room: the letter it adds to the path, and the
// direction it goes. The moves are in the same order as the first four hex
// characters of the hash.
type move struct {
	letter byte
	dir    grid.Point
}

// moves lists the doors in hash order: up, down, left, right.
var moves = [4]move{
	{'U', grid.Up},
	{'D', grid.Down},
	{'L', grid.Left},
	{'R', grid.Right},
}

// state is a room and the path that led to it. The same room with a different
// path has different doors open, so the path is part of the state.
type state struct {
	pos  grid.Point
	path string
}

// part1 finds the shortest path to the vault.
func part1(in string) any {
	shortest, _, ok := walk(strings.TrimSpace(in))
	if !ok {
		return nil
	}

	return shortest
}

// part2 finds the length of the longest path to the vault.
func part2(in string) any {
	_, longest, ok := walk(strings.TrimSpace(in))
	if !ok {
		return nil
	}

	return longest
}

// walk tries every path from the start, one step at a time for all paths
// together (a breadth-first search). A path stops when it reaches the vault,
// or when all its doors are locked. The first path to reach the vault is the
// shortest, and the last one is the longest. The result is false when no
// path reaches the vault.
func walk(passcode string) (shortest string, longest int, ok bool) {
	vault := grid.P(size-1, size-1)
	queue := []state{{pos: grid.P(0, 0)}}

	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]

		if s.pos == vault {
			if !ok {
				shortest, ok = s.path, true
			}

			longest = len(s.path)
			continue
		}

		sum := md5.Sum([]byte(passcode + s.path))

		for i, m := range moves {
			// Hex character i is the high nibble of byte i/2 when i is even,
			// and the low nibble when i is odd. A door is open for b to f.
			nibble := sum[i/2] >> (4 * (1 - i%2)) & 0xf
			if nibble < 0xb {
				continue
			}

			next := s.pos.Add(m.dir)
			if next.X < 0 || next.X >= size || next.Y < 0 || next.Y >= size {
				continue
			}

			queue = append(queue, state{pos: next, path: s.path + string(m.letter)})
		}
	}

	return shortest, longest, ok
}
