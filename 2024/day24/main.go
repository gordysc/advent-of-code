// Advent of Code 2024, day 24: Crossed Wires.
// https://adventofcode.com/2024/day/24
package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2024, 24, part1, part2)
}

// part1 runs the circuit and returns the number that the z wires make, with
// z00 as the lowest bit.
//
// The puzzle has a small and a larger example. The file example.txt holds
// the larger one, which gives 2024. The small one gives 4.
func part1(in string) any {
	c := parse(in)
	result := 0

	for _, g := range c.gates {
		if g.out[0] != 'z' {
			continue
		}

		bit, err := strconv.Atoi(g.out[1:])
		if err != nil {
			panic("bad z wire: " + g.out)
		}

		result |= c.value(g.out) << bit
	}

	return result
}

// part2 finds the eight wires whose gate outputs were swapped in pairs, and
// returns their names in order, joined by commas.
//
// The circuit should add the x bits to the y bits as a ripple-carry adder.
// Each bit i (except bit 0) of that adder uses five gates:
//
//	xi XOR yi -> sum        (half sum, no carry yet)
//	xi AND yi -> direct     (carry from xi and yi)
//	sum XOR carryIn -> zi   (the output bit)
//	sum AND carryIn -> via  (carry that goes through the sum)
//	direct OR via -> carryOut
//
// Bit 0 has no carry in, so it is only "x00 XOR y00 -> z00" and
// "x00 AND y00 -> carryOut". The last z wire is the carry out of the top bit.
//
// We do not need to find the pairs. A wire is wrong when its gate breaks one
// of these rules, which all follow from the pattern above:
//
//  1. A z wire comes from XOR. Only the last z wire comes from OR.
//  2. A XOR gate that does not read x and y must write a z wire.
//  3. A XOR gate that reads x and y (except bit 0) must feed a XOR gate.
//  4. An AND gate (except bit 0) must feed an OR gate.
//  5. Only AND gates feed an OR gate.
//
// The part 2 example uses AND gates in place of an adder, so these rules do
// not apply to it. Part 2 returns nil when the z wires do not fit an adder
// for the x and y bits, which is the case for both examples.
func part2(in string) any {
	c := parse(in)

	xBits, zBits := 0, 0
	for w := range c.values {
		if w[0] == 'x' {
			xBits++
		}
	}

	for _, g := range c.gates {
		if g.out[0] == 'z' {
			zBits++
		}
	}

	if xBits < 2 || zBits != xBits+1 {
		return nil
	}

	lastZ := fmt.Sprintf("z%02d", xBits)

	// feeds[w] lists the operations of the gates that read wire w.
	feeds := map[string][]string{}
	for _, g := range c.gates {
		feeds[g.a] = append(feeds[g.a], g.op)
		feeds[g.b] = append(feeds[g.b], g.op)
	}

	wrong := set.New[string]()

	for _, g := range c.gates {
		readsXY := isInput(g.a) && isInput(g.b)
		bit0 := g.a[1:] == "00" && readsXY
		isZ := g.out[0] == 'z'

		switch {
		case isZ && g.out != lastZ && g.op != "XOR":
			wrong.Add(g.out) // rule 1
		case g.out == lastZ && g.op != "OR":
			wrong.Add(g.out) // rule 1
		case g.op == "XOR" && !readsXY && !isZ:
			wrong.Add(g.out) // rule 2
		case g.op == "XOR" && readsXY && !bit0 && !slices.Contains(feeds[g.out], "XOR"):
			wrong.Add(g.out) // rule 3
		case g.op == "AND" && !bit0 && !slices.Contains(feeds[g.out], "OR"):
			wrong.Add(g.out) // rule 4
		case g.op != "AND" && slices.Contains(feeds[g.out], "OR"):
			wrong.Add(g.out) // rule 5
		}
	}

	names := wrong.Items()
	slices.Sort(names)

	return strings.Join(names, ",")
}

// isInput reports whether a wire is one of the x or y input bits.
func isInput(w string) bool {
	return w[0] == 'x' || w[0] == 'y'
}

// gate is one logic gate: a op b -> out.
type gate struct {
	a, op, b, out string
}

// circuit holds the wire values that are known and the gates.
type circuit struct {
	values map[string]int
	gates  []gate
	byOut  map[string]gate
}

// parse reads the initial wire values, a blank line, and then one gate per
// line, such as "x00 AND y00 -> z00".
func parse(in string) circuit {
	c := circuit{values: map[string]int{}, byOut: map[string]gate{}}
	blocks := input.Blocks(in)

	for _, line := range blocks[0] {
		name, value, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}

		c.values[name] = input.Int(value)
	}

	for _, line := range blocks[1] {
		f := strings.Fields(line)
		if len(f) != 5 {
			continue
		}

		g := gate{a: f[0], op: f[1], b: f[2], out: f[4]}
		c.gates = append(c.gates, g)
		c.byOut[g.out] = g
	}

	return c
}

// value returns the value of wire w. It works back through the gates and
// stores each result in c.values, so no gate runs twice.
//
// c.values is a map, and a map in Go is a reference to shared data. So the
// method changes the circuit's map even though c is not a pointer.
func (c circuit) value(w string) int {
	if v, ok := c.values[w]; ok {
		return v
	}

	g := c.byOut[w]
	a, b := c.value(g.a), c.value(g.b)

	var v int
	switch g.op {
	case "AND":
		v = a & b
	case "OR":
		v = a | b
	case "XOR":
		v = a ^ b
	}

	c.values[w] = v

	return v
}
