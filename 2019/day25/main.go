// Advent of Code 2019, day 25: Cryostasis.
// https://adventofcode.com/2019/day/25
//
// Day 25 has no second puzzle: part 2 is unlocked by collecting the other 49
// stars, so only part 1 is written here.
//
// The puzzle text has no example. example.txt holds a small handmade Intcode
// program: it prints a Hull Breach room with one door north, reads one
// command, and then prints the text of a successful step onto the pressure
// plate, with the password 2424. It does not model items or weights, so it
// tests the room parser and the password check, but not the item search.
package main

import (
	"math/bits"
	"regexp"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 25, part1, nil)
}

// dangerous lists the items that end the game or trap the droid when taken.
var dangerous = map[string]bool{
	"giant electromagnet": true,
	"infinite loop":       true,
	"molten lava":         true,
	"photons":             true,
	"escape pod":          true,
}

// opposite maps each direction to the one that walks back.
var opposite = map[string]string{
	"north": "south",
	"south": "north",
	"east":  "west",
	"west":  "east",
}

// passwordRe matches the line that gives the password for the main airlock.
var passwordRe = regexp.MustCompile(`typing (\d+) on the keypad`)

// room is what the droid learns from one room description.
type room struct {
	name  string
	doors []string
	items []string
}

// droid drives the Intcode program. items is what it has picked up, and
// checkpoint and floorDoor are the path from the start to the security
// checkpoint and the door from there onto the pressure-sensitive floor.
// password is set as soon as any reply gives it.
type droid struct {
	m          *intcode.Machine
	items      []string
	checkpoint []string
	floorDoor  string
	password   string
}

// part1 finds the password for the main airlock.
//
// The droid explores the whole ship and picks up every safe item. Then it goes
// to the security checkpoint and tries sets of items until the floor lets it
// through.
func part1(in string) any {
	d := &droid{m: intcode.New(intcode.Parse(in))}
	start := parseRoom(d.reply(d.m.RunString()))

	d.explore(start, nil, map[string]bool{})
	if d.password != "" {
		return input.Int(d.password)
	}

	if d.floorDoor == "" {
		return nil
	}

	for _, dir := range d.checkpoint {
		d.do(dir)
	}

	return d.tryItems()
}

// do sends one command and returns the program's reply.
func (d *droid) do(cmd string) string {
	d.m.SendLine(cmd)

	return d.reply(d.m.RunString())
}

// reply looks at a reply for the password before it goes back to the caller.
//
// The program halts after it gives the password. A halt at any other time
// means the droid did something fatal, and every later command would get an
// empty reply. The droid could then walk in circles, so it stops with a panic
// that shows the last reply.
func (d *droid) reply(text string) string {
	if m := passwordRe.FindStringSubmatch(text); m != nil {
		d.password = m[1]
	}

	if d.m.Halted() && d.password == "" {
		panic("day25: the program halted with no password:\n" + text)
	}

	return text
}

// explore does a depth-first walk from room r, which the path of directions
// leads to from the start. It takes the safe items in each room, walks
// through every door, and walks back the opposite way afterward, so it ends
// where it began.
//
// A step onto the pressure-sensitive floor with the wrong weight throws the
// droid back to the checkpoint with an "Alert!". That is how explore finds the
// checkpoint and the door to the floor. It does not go through that door
// again.
func (d *droid) explore(r room, path []string, seen map[string]bool) {
	seen[r.name] = true

	for _, item := range r.items {
		if !dangerous[item] {
			d.do("take " + item)
			d.items = append(d.items, item)
		}
	}

	for _, door := range r.doors {
		text := d.do(door)
		if d.password != "" {
			return
		}

		if strings.Contains(text, "Alert!") {
			d.checkpoint = slices.Clone(path)
			d.floorDoor = door
			continue
		}

		// Clip the path so the append below never writes into a slice that
		// a caller or the saved checkpoint path still uses.
		next := parseRoom(text)
		if !seen[next.name] {
			d.explore(next, append(slices.Clip(path), door), seen)
			if d.password != "" {
				return
			}
		}

		d.do(opposite[door])
	}
}

// tryItems stands at the checkpoint with every item and tries each set of
// items on the floor, until one has the right weight.
//
// Set number i is the Gray code i ^ (i >> 1): bit k says whether item k is
// held. Two Gray codes in a row differ in exactly one bit, so each new set
// needs only one take or drop. The trailing zeros of the XOR of two codes
// give the position of that bit.
func (d *droid) tryItems() any {
	for _, item := range d.items {
		d.do("drop " + item)
	}

	held := 0
	for i := range 1 << len(d.items) {
		want := i ^ i>>1

		if diff := want ^ held; diff != 0 {
			item := d.items[bits.TrailingZeros(uint(diff))]
			if want&diff != 0 {
				d.do("take " + item)
			} else {
				d.do("drop " + item)
			}
		}

		held = want
		d.do(d.floorDoor)

		if d.password != "" {
			return input.Int(d.password)
		}
	}

	return nil
}

// parseRoom reads a room description: the name between "==" marks, and the
// "- " lines under "Doors here lead:" and "Items here:".
//
// A step onto the floor with the wrong weight prints two rooms: the floor,
// then the checkpoint the droid is thrown back to. A new name line starts the
// room again, so the last room in the text wins. That is the room the droid
// is in.
func parseRoom(text string) room {
	var r room

	// list points at the slice that "- " lines go into: r.doors, r.items, or
	// nil outside a list. Through the pointer, one append case fills both.
	var list *[]string

	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, "== ") && strings.HasSuffix(line, " =="):
			r = room{name: strings.Trim(line, "= ")}
			list = nil

		case line == "Doors here lead:":
			list = &r.doors

		case line == "Items here:":
			list = &r.items

		case strings.HasPrefix(line, "- ") && list != nil:
			*list = append(*list, line[2:])

		default:
			// Any other line, such as the blank line after a list, ends it.
			list = nil
		}
	}

	return r
}
