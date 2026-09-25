// Advent of Code 2016, day 4: Security Through Obscurity.
// https://adventofcode.com/2016/day/4
package main

import (
	"cmp"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/strx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 4, part1, part2)
}

// room is one line of the kiosk list: an encrypted name, a sector ID and a checksum.
type room struct {
	name     string
	sector   int
	checksum string
}

// part1 adds up the sector IDs of the real rooms.
func part1(in string) any {
	total := 0

	for _, r := range parse(in) {
		if r.real() {
			total += r.sector
		}
	}

	return total
}

// part2 finds the sector ID of the room where the North Pole objects are stored.
func part2(in string) any {
	for _, r := range parse(in) {
		if r.real() && strings.Contains(r.decrypt(), "northpole") {
			return r.sector
		}
	}

	return nil
}

// parse reads lines such as "aaaaa-bbb-z-y-x-123[abxyz]" into rooms.
func parse(in string) []room {
	var rooms []room

	for _, line := range input.Lines(in) {
		dash := strings.LastIndexByte(line, '-')

		rooms = append(rooms, room{
			name:     line[:dash],
			sector:   input.Int(line[dash+1 : strings.IndexByte(line, '[')]),
			checksum: strx.Between(line, "[", "]"),
		})
	}

	return rooms
}

// real reports whether the checksum holds the five most common letters of the
// name, most common first. Letters with the same count go in alphabetical order.
func (r room) real() bool {
	var counts [26]int

	for _, c := range r.name {
		if c != '-' {
			counts[c-'a']++
		}
	}

	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	slices.SortStableFunc(letters, func(a, b byte) int {
		return cmp.Compare(counts[b-'a'], counts[a-'a'])
	})

	return string(letters[:5]) == r.checksum
}

// decrypt rotates each letter of the name forward by the sector ID. Dashes
// become spaces.
func (r room) decrypt() string {
	var out strings.Builder

	for _, c := range []byte(r.name) {
		if c == '-' {
			out.WriteByte(' ')
			continue
		}

		out.WriteByte('a' + byte((int(c-'a')+r.sector)%26))
	}

	return out.String()
}
