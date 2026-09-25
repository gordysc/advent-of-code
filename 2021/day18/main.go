// Advent of Code 2021, day 18: Snailfish.
// https://adventofcode.com/2021/day/18
package main

import (
	"slices"

	"aoc/lib/input"
	"aoc/lib/slicesx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 18, part1, part2)
}

// part1 adds all the snailfish numbers in order and returns the magnitude of
// the final sum.
//
// A snailfish number is kept as a flat list of its regular numbers, from left
// to right, each with its nesting depth. Explode and split then only change
// items next to each other in that list, so no tree is necessary.
func part1(in string) any {
	numbers := slicesx.Map(input.Lines(in), parse)

	sum := numbers[0]
	for _, n := range numbers[1:] {
		sum = add(sum, n)
	}

	return magnitude(sum)
}

// part2 finds the largest magnitude from the sum of two different numbers.
//
// Snailfish addition is not commutative, so each pair is tried in both
// orders.
func part2(in string) any {
	numbers := slicesx.Map(input.Lines(in), parse)
	best := 0

	for i, a := range numbers {
		for j, b := range numbers {
			if i == j {
				continue
			}

			best = max(best, magnitude(add(a, b)))
		}
	}

	return best
}

// leaf is one regular number in a snailfish number, with the number of pairs
// that enclose it.
type leaf struct {
	value int
	depth int
}

// number is a snailfish number as a flat list of leaves, from left to right.
type number []leaf

// parse reads one snailfish number. All regular numbers in the input have
// one digit.
func parse(line string) number {
	var n number
	depth := 0

	for _, c := range line {
		switch {
		case c == '[':
			depth++
		case c == ']':
			depth--
		case c >= '0' && c <= '9':
			n = append(n, leaf{value: int(c - '0'), depth: depth})
		}
	}

	return n
}

// add makes a new pair from a and b and then reduces it. The inputs are not
// changed, so part2 can use them again.
func add(a, b number) number {
	sum := make(number, 0, len(a)+len(b))

	// The new pair encloses every leaf of both numbers one more time.
	for _, l := range slices.Concat(a, b) {
		sum = append(sum, leaf{value: l.value, depth: l.depth + 1})
	}

	return reduce(sum)
}

// reduce does explode and split actions until no more action applies.
//
// Explode always has priority. The two inputs are already reduced, so only
// pairs at depth 5 can explode. Such a pair always has two regular numbers,
// so they are two items next to each other in the list with depth 5.
func reduce(n number) number {
	for {
		if i := slices.IndexFunc(n, func(l leaf) bool { return l.depth > 4 }); i >= 0 {
			n = explode(n, i)
			continue
		}

		if i := slices.IndexFunc(n, func(l leaf) bool { return l.value >= 10 }); i >= 0 {
			n = split(n, i)
			continue
		}

		return n
	}
}

// explode removes the pair whose left number is at index i. Its numbers go to
// the nearest regular numbers on each side, and the pair becomes a 0.
func explode(n number, i int) number {
	left, right := n[i], n[i+1]

	if i > 0 {
		n[i-1].value += left.value
	}

	if i+2 < len(n) {
		n[i+2].value += right.value
	}

	n[i] = leaf{value: 0, depth: left.depth - 1}

	return slices.Delete(n, i+1, i+2)
}

// split replaces the regular number at index i with a pair. The left half is
// rounded down and the right half is rounded up.
func split(n number, i int) number {
	l := n[i]
	half := leaf{value: l.value / 2, depth: l.depth + 1}
	other := leaf{value: l.value - half.value, depth: l.depth + 1}

	n[i] = half

	return slices.Insert(n, i+1, other)
}

// magnitude calculates 3 times the left magnitude plus 2 times the right
// magnitude, from the bottom up.
//
// The leaves go on a stack from left to right. When the two top items have
// the same depth, they are the left and right halves of one pair, because
// everything to the left of the top item is already fully combined. So they
// are replaced by the magnitude of that pair, one level up.
func magnitude(n number) int {
	var stack []leaf

	for _, l := range n {
		stack = append(stack, l)

		for len(stack) >= 2 {
			a, b := stack[len(stack)-2], stack[len(stack)-1]
			if a.depth != b.depth {
				break
			}

			stack = stack[:len(stack)-2]
			stack = append(stack, leaf{value: 3*a.value + 2*b.value, depth: a.depth - 1})
		}
	}

	return stack[0].value
}
