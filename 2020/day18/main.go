// Advent of Code 2020, day 18: Operation Order.
// https://adventofcode.com/2020/day/18
package main

import (
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 18, part1, part2)
}

// part1 evaluates each line with + and * at the same precedence, from left
// to right, and adds up the results.
//
// example.txt holds the six example expressions from the puzzle, one per
// line. Their part 1 values are 71, 51, 26, 437, 12240 and 13632, so the
// sum is 26457.
func part1(in string) any {
	return sumAll(in, map[byte]int{'+': 1, '*': 1})
}

// part2 evaluates each line with + before *, and adds up the results. The
// part 2 values of the examples are 231, 51, 46, 1445, 669060 and 23340, so
// the sum is 694173.
func part2(in string) any {
	return sumAll(in, map[byte]int{'+': 2, '*': 1})
}

// sumAll evaluates every line with the given operator precedence and
// returns the sum of the results.
func sumAll(in string, prec map[byte]int) int {
	total := 0
	for _, line := range input.Lines(in) {
		p := parser{src: line, prec: prec}
		total += p.expr(1)
	}

	return total
}

// parser evaluates one expression with precedence climbing. The prec map
// gives the binding strength of each operator. A higher number binds more
// tightly. All operators are left-associative.
type parser struct {
	src  string
	pos  int
	prec map[byte]int
}

// peek skips spaces and returns the next character, or 0 at the end.
func (p *parser) peek() byte {
	for p.pos < len(p.src) && p.src[p.pos] == ' ' {
		p.pos++
	}

	if p.pos == len(p.src) {
		return 0
	}

	return p.src[p.pos]
}

// expr reads operands and operators while the next operator binds at least
// as tightly as minPrec. For the right operand it asks for a precedence one
// higher, so equal operators group from left to right, and a weaker operator
// ends the right operand and returns control to the caller.
//
// The methods use a pointer receiver (*parser) so that each call moves the
// same pos forward. A value receiver would move a copy.
func (p *parser) expr(minPrec int) int {
	left := p.operand()

	for {
		op := p.peek()
		prec, ok := p.prec[op]
		if !ok || prec < minPrec {
			return left
		}

		p.pos++
		right := p.expr(prec + 1)

		if op == '+' {
			left += right
		} else {
			left *= right
		}
	}
}

// operand reads a number or a parenthesised sub-expression. The puzzle only
// uses single-digit numbers, but the loop reads longer numbers too.
func (p *parser) operand() int {
	if p.peek() == '(' {
		p.pos++
		v := p.expr(1)
		p.pos++ // skip ')'

		return v
	}

	v := 0
	for p.pos < len(p.src) && p.src[p.pos] >= '0' && p.src[p.pos] <= '9' {
		v = v*10 + int(p.src[p.pos]-'0')
		p.pos++
	}

	return v
}
