// Advent of Code 2016, day 16: Dragon Checksum.
// https://adventofcode.com/2016/day/16
package main

import (
	"strings"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 16, part1, part2)
}

// diskSmall and diskLarge are the sizes of the two disks in the puzzle. The
// worked example in the puzzle text fills a disk of 20 instead, where "10000"
// gives the checksum "01100". On these disks the example gives other checksums.
const (
	diskSmall = 272
	diskLarge = 35651584
)

// part1 fills the first disk and returns its checksum.
func part1(in string) any {
	return checksum(fill(strings.TrimSpace(in), diskSmall))
}

// part2 fills the much larger second disk and returns its checksum.
func part2(in string) any {
	return checksum(fill(strings.TrimSpace(in), diskLarge))
}

// fill grows the state with the dragon curve until it covers the disk, then
// cuts it to the disk size. Each step appends a '0' and then the data so far,
// reversed and with every bit flipped. The digits are the bytes '0' and '1',
// so XOR with 1 flips a bit.
func fill(state string, size int) []byte {
	// One step at most doubles the data, so this capacity is always enough.
	data := make([]byte, 0, 2*size+1)
	data = append(data, state...)

	for len(data) < size {
		n := len(data)
		data = append(data, '0')

		for i := n - 1; i >= 0; i-- {
			data = append(data, data[i]^1)
		}
	}

	return data[:size]
}

// checksum halves the data while its length is even. Each pair of bits
// becomes '1' when the two match and '0' when they differ. The work happens
// in place, because the result of a pair is written behind the pairs that
// are still to be read.
func checksum(data []byte) string {
	for len(data)%2 == 0 {
		half := len(data) / 2

		for i := range half {
			if data[2*i] == data[2*i+1] {
				data[i] = '1'
			} else {
				data[i] = '0'
			}
		}

		data = data[:half]
	}

	return string(data)
}
