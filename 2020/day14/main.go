// Advent of Code 2020, day 14: Docking Data.
// https://adventofcode.com/2020/day/14
package main

import (
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 14, part1, part2)
}

// part1 applies the mask to each value before it goes into memory. A 1 or 0
// in the mask sets that bit of the value, and an X keeps the bit. The answer
// is the sum of all values in memory.
//
// example.txt is the part 2 example. The part 1 example has a mask with 34
// X bits, so part 2 would write to 2^34 addresses for it. Part 1 gives 51 for
// the part 2 example and not the 165 from the part 1 example.
func part1(in string) any {
	mem := map[uint64]uint64{}
	var m mask

	for _, line := range input.Lines(in) {
		if s, ok := strings.CutPrefix(line, "mask = "); ok {
			m = parseMask(s)
			continue
		}

		addr, value := parseWrite(line)
		mem[addr] = value&m.floating | m.ones
	}

	return sum(mem)
}

// part2 applies the mask to each address. A 1 sets that bit of the address,
// a 0 keeps the bit, and an X bit floats. A floating bit takes both values,
// so a write goes to 2^n addresses when the mask has n X bits. The answer is
// the sum of all values in memory.
func part2(in string) any {
	mem := map[uint64]uint64{}
	var m mask

	for _, line := range input.Lines(in) {
		if s, ok := strings.CutPrefix(line, "mask = "); ok {
			m = parseMask(s)
			continue
		}

		addr, value := parseWrite(line)
		base := (addr | m.ones) &^ m.floating

		// Count down through all subsets of the floating bits. The value
		// sub-1 & floating gives the next smaller subset, and the loop
		// stops after it writes the empty subset.
		for sub := m.floating; ; sub = (sub - 1) & m.floating {
			mem[base|sub] = value
			if sub == 0 {
				break
			}
		}
	}

	return sum(mem)
}

// mask holds the 36-bit mask as two bit sets: the bits that are 1 and the
// bits that are X. All other bits are 0.
type mask struct {
	ones     uint64
	floating uint64
}

// parseMask reads a mask string, most significant bit first.
func parseMask(s string) mask {
	var m mask
	for _, c := range s {
		m.ones <<= 1
		m.floating <<= 1

		switch c {
		case '1':
			m.ones |= 1
		case 'X':
			m.floating |= 1
		}
	}

	return m
}

// parseWrite reads a "mem[addr] = value" line into its address and value.
func parseWrite(line string) (uint64, uint64) {
	addr, value, ok := strings.Cut(strings.TrimPrefix(line, "mem["), "] = ")
	if !ok {
		panic("bad line " + line)
	}

	a, err := strconv.ParseUint(addr, 10, 64)
	if err != nil {
		panic(err)
	}

	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		panic(err)
	}

	return a, v
}

// sum adds all values in memory.
func sum(mem map[uint64]uint64) uint64 {
	var total uint64
	for _, v := range mem {
		total += v
	}

	return total
}
