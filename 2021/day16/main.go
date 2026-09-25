// Advent of Code 2021, day 16: Packet Decoder.
// https://adventofcode.com/2021/day/16
package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 16, part1, part2)
}

// part1 adds up the version numbers of all the packets.
//
// The decoder reads the outer packet and all the packets inside it into a
// tree. Then the sum walks the tree.
func part1(in string) any {
	return decode(in).versionSum()
}

// part2 calculates the value of the outer packet.
//
// The type ID of each operator packet tells which operation to do on the
// values of its sub-packets.
func part2(in string) any {
	return decode(in).value()
}

// packet is one decoded packet. A literal packet has a value and no
// sub-packets. An operator packet has sub-packets and no literal value.
type packet struct {
	version int
	typeID  int
	literal int
	subs    []packet
}

// typeLiteral is the type ID of a packet that holds a literal value.
const typeLiteral = 4

// versionSum adds the version of this packet and of all the packets inside it.
func (p packet) versionSum() int {
	sum := p.version

	for _, sub := range p.subs {
		sum += sub.versionSum()
	}

	return sum
}

// value calculates the value of the packet from its type ID.
func (p packet) value() int {
	if p.typeID == typeLiteral {
		return p.literal
	}

	values := make([]int, len(p.subs))
	for i, sub := range p.subs {
		values[i] = sub.value()
	}

	switch p.typeID {
	case 0:
		return slicesx.Sum(values)
	case 1:
		return slicesx.Product(values)
	case 2:
		return slices.Min(values)
	case 3:
		return slices.Max(values)
	case 5:
		return boolInt(values[0] > values[1])
	case 6:
		return boolInt(values[0] < values[1])
	case 7:
		return boolInt(values[0] == values[1])
	}

	panic(fmt.Sprintf("unknown type ID %d", p.typeID))
}

// boolInt returns 1 for true and 0 for false.
func boolInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

// reader reads numbers from a string of bits, from left to right.
type reader struct {
	bits string
	pos  int
}

// read returns the next n bits as a number.
func (r *reader) read(n int) int {
	v, err := strconv.ParseInt(r.bits[r.pos:r.pos+n], 2, 64)
	if err != nil {
		panic(err)
	}

	r.pos += n

	return int(v)
}

// decode turns the hexadecimal transmission into the outer packet.
func decode(in string) packet {
	var bits strings.Builder

	for _, c := range strings.TrimSpace(in) {
		nibble, err := strconv.ParseUint(string(c), 16, 8)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(&bits, "%04b", nibble)
	}

	r := &reader{bits: bits.String()}

	return r.packet()
}

// packet reads one packet and all the packets inside it.
func (r *reader) packet() packet {
	p := packet{version: r.read(3), typeID: r.read(3)}

	if p.typeID == typeLiteral {
		// A literal is a series of 5-bit groups. The first bit of each group
		// is 1 for all groups but the last one.
		for {
			group := r.read(5)
			p.literal = p.literal<<4 | group&0xF

			if group&0x10 == 0 {
				return p
			}
		}
	}

	// The length type ID tells how the size of the sub-packets is given:
	// 0 gives the total length in bits, 1 gives the number of sub-packets.
	if r.read(1) == 0 {
		length := r.read(15)
		end := r.pos + length

		for r.pos < end {
			p.subs = append(p.subs, r.packet())
		}

		return p
	}

	count := r.read(11)
	for range count {
		p.subs = append(p.subs, r.packet())
	}

	return p
}
