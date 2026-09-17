// Advent of Code 2015, day 7: Some Assembly Required.
// https://adventofcode.com/2015/day/7
package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 7, part1, part2)
}

// gate is one parsed instruction: the operator and the operands that feed it.
// An operand is kept as the raw text from the input, which is either a wire
// name or a literal number; the evaluator tells the two apart when it needs
// the value. A plain assignment such as "123 -> x" or "lx -> a" has an empty
// operator and a single operand.
type gate struct {
	op   string
	args []string
}

// circuit maps every wire to the gate that drives it, plus a cache of the
// signals already worked out.
type circuit struct {
	gates map[string]gate
	memo  map[string]uint16
}

// part1 returns the signal that ends up on wire a.
func part1(in string) any {
	return parse(in).answer()
}

// part2 feeds the part 1 signal from wire a back into wire b, overriding the
// gate that drove b, and reports the new signal on a.
func part2(in string) any {
	c := parse(in)
	if _, ok := c.gates["a"]; !ok {
		return nil
	}

	a := c.signal("a")

	// Seeding a fresh cache with b's value is enough to override it: signal
	// returns a cached value before it ever looks at the gate behind the wire.
	// The cache must be fresh because every other value depended on the old b.
	c.memo = map[string]uint16{"b": a}

	return c.signal("a")
}

// answer returns the signal on wire a. The puzzle's example has no wire a, so
// when it is missing every wire is listed instead, one per line, to match the
// table shown in the puzzle text.
func (c *circuit) answer() any {
	if _, ok := c.gates["a"]; ok {
		return c.signal("a")
	}

	var lines []string
	for _, wire := range slices.Sorted(maps.Keys(c.gates)) {
		lines = append(lines, fmt.Sprintf("%s: %d", wire, c.signal(wire)))
	}

	return strings.Join(lines, "\n")
}

// signal returns the value on a wire, first working out every wire it depends
// on. The result is cached: a wire that feeds many gates is evaluated once,
// and without the cache the recursion would revisit shared sub-circuits over
// and over, which is exponential in the depth of the circuit.
func (c *circuit) signal(wire string) uint16 {
	if v, ok := c.memo[wire]; ok {
		return v
	}

	g, ok := c.gates[wire]
	if !ok {
		panic("no gate drives wire " + wire)
	}

	// Every operation is on 16-bit values, so uint16 arithmetic gives the
	// wrap-around the puzzle asks for with no masking. ^ on an unsigned
	// integer is bitwise NOT in Go, so NOT 123 is 65412 as the puzzle shows.
	var v uint16
	switch g.op {
	case "":
		v = c.operand(g.args[0])
	case "NOT":
		v = ^c.operand(g.args[0])
	case "AND":
		v = c.operand(g.args[0]) & c.operand(g.args[1])
	case "OR":
		v = c.operand(g.args[0]) | c.operand(g.args[1])
	case "LSHIFT":
		v = c.operand(g.args[0]) << c.operand(g.args[1])
	case "RSHIFT":
		v = c.operand(g.args[0]) >> c.operand(g.args[1])
	default:
		panic("unknown operator " + g.op)
	}

	c.memo[wire] = v

	return v
}

// operand resolves one operand: a literal number gives its value directly,
// and anything else is a wire name to evaluate.
func (c *circuit) operand(s string) uint16 {
	if s[0] >= '0' && s[0] <= '9' {
		return uint16(input.Int(s))
	}

	return c.signal(s)
}

// parse builds the circuit from the input. Each line is "<expression> -> <wire>",
// where the expression is one operand, "NOT" and one operand, or two operands
// around a binary operator.
func parse(in string) *circuit {
	c := &circuit{
		gates: make(map[string]gate),
		memo:  make(map[string]uint16),
	}

	for _, line := range input.Lines(in) {
		expr, wire, ok := strings.Cut(line, " -> ")
		if !ok {
			panic("malformed line: " + line)
		}

		f := strings.Fields(expr)
		switch len(f) {
		case 1:
			c.gates[wire] = gate{args: f}
		case 2:
			c.gates[wire] = gate{op: f[0], args: f[1:]}
		default:
			c.gates[wire] = gate{op: f[1], args: []string{f[0], f[2]}}
		}
	}

	return c
}
