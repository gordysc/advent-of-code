// Advent of Code 2024, day 9: Disk Fragmenter.
// https://adventofcode.com/2024/day/9
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 9, part1, part2)
}

// part1 moves file blocks one at a time, from the end of the disk to the
// leftmost free block, until no gap is left. It returns the checksum.
func part1(in string) any {
	files, _ := parse(in)

	// Expand the disk map to one entry per block. -1 marks a free block.
	var disk []int
	for id, f := range files {
		for len(disk) < f.start {
			disk = append(disk, -1)
		}

		// Since Go 1.22, range over an int n runs the loop n times.
		for range f.size {
			disk = append(disk, id)
		}
	}

	// Two pointers: left looks for the next free block, right for the last
	// file block. Stop when they meet.
	left, right := 0, len(disk)-1
	for {
		for left < right && disk[left] != -1 {
			left++
		}

		for left < right && disk[right] == -1 {
			right--
		}

		if left >= right {
			break
		}

		disk[left], disk[right] = disk[right], -1
	}

	sum := 0

	for pos, id := range disk {
		if id != -1 {
			sum += pos * id
		}
	}

	return sum
}

// part2 moves whole files, one time each, in order of decreasing file ID.
// Each file goes to the leftmost free span that fits it, but only if that
// span is to the left of the file. It returns the checksum.
//
// When a file moves, its old blocks become free. We do not add them back as
// free spans: every file that moves later has a lower ID, so it is to the
// left of this file and can never move to the right.
func part2(in string) any {
	files, spans := parse(in)

	for id := len(files) - 1; id >= 0; id-- {
		// Take pointers to the slice elements, so that the changes below go
		// into the slices and not into copies.
		f := &files[id]

		// spans is in disk order, so the first span that fits is the leftmost.
		for i := range spans {
			s := &spans[i]
			if s.start >= f.start {
				break
			}

			if s.size >= f.size {
				f.start = s.start
				s.start += f.size
				s.size -= f.size

				break
			}
		}
	}

	sum := 0
	for id, f := range files {
		for pos := f.start; pos < f.start+f.size; pos++ {
			sum += pos * id
		}
	}

	return sum
}

// span is a run of blocks on the disk: a file, or a run of free blocks.
type span struct {
	start, size int
}

// parse reads the disk map. The digits alternate between a file size and a
// free space size. It returns the files, indexed by file ID, and the free
// spans in disk order. Both are in increasing order of start position.
func parse(in string) ([]span, []span) {
	var files, spans []span
	pos := 0

	// Ranging over a string gives byte offsets and runes. Each digit is one
	// byte, so the offset i is also the digit's index.
	for i, c := range strings.TrimSpace(in) {
		size := int(c - '0')

		if i%2 == 0 {
			files = append(files, span{pos, size})
		} else if size > 0 {
			spans = append(spans, span{pos, size})
		}

		pos += size
	}

	return files, spans
}
