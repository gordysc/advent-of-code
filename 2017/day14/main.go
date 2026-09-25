// Advent of Code 2017, day 14: Disk Defragmentation.
// https://adventofcode.com/2017/day/14
package main

import (
	"fmt"
	"math/bits"
	"strings"

	"aoc/lib/grid"
	"aoc/lib/search"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 14, part1, part2)
}

// size is the width and height of the disk grid, in squares.
const size = 128

// part1 counts the used squares. Each row is a knot hash, so it is enough to
// count the one bits in every hash.
func part1(in string) any {
	used := 0

	for row := range size {
		for _, b := range knotHash(rowKey(in, row)) {
			used += bits.OnesCount8(b)
		}
	}

	return used
}

// part2 counts the regions of used squares that touch up, down, left or right.
// Each used square that is not yet in a region starts a new one, and a flood
// fill from it marks every square in that region.
func part2(in string) any {
	disk := build(in)
	next := func(p grid.Point) []grid.Point {
		var out []grid.Point
		for _, n := range disk.Neighbors4(p) {
			if disk.At(n) {
				out = append(out, n)
			}
		}

		return out
	}

	seen := map[grid.Point]bool{}
	regions := 0

	for p, used := range disk.All() {
		if !used || seen[p] {
			continue
		}

		regions++
		for q := range search.Flood(p, next) {
			seen[q] = true
		}
	}

	return regions
}

// build makes the disk grid, where true marks a used square. Bit 7 of the
// first hash byte is the leftmost square of the row.
func build(in string) grid.Grid[bool] {
	disk := grid.New[bool](size, size)

	for y := range size {
		hash := knotHash(rowKey(in, y))
		for x := range size {
			bit := hash[x/8] >> (7 - x%8) & 1
			disk.Set(grid.P(x, y), bit == 1)
		}
	}

	return disk
}

// rowKey is the hash input for one row: the key string, a dash, and the row
// number.
func rowKey(key string, row int) string {
	return fmt.Sprintf("%s-%d", strings.TrimSpace(key), row)
}

// knotHash is the full knot hash from day 10. The lengths are the bytes of s
// plus a fixed suffix. After 64 rounds of twists, each block of 16 numbers is
// XORed together to make the 16-byte dense hash.
func knotHash(s string) [16]byte {
	lengths := append([]byte(s), 17, 31, 73, 47, 23)

	var list [256]byte
	for i := range list {
		list[i] = byte(i)
	}

	pos, skip := 0, 0
	for range 64 {
		for _, l := range lengths {
			n := int(l)
			for i := range n / 2 {
				a, b := (pos+i)%256, (pos+n-1-i)%256
				list[a], list[b] = list[b], list[a]
			}
			pos = (pos + n + skip) % 256
			skip++
		}
	}

	var dense [16]byte
	for i, v := range list {
		dense[i/16] ^= v
	}

	return dense
}
