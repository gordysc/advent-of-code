// Advent of Code 2017, day 21: Fractal Art.
// https://adventofcode.com/2017/day/21
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/lib/strx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2017, 21, part1, part2)
}

// part1Iterations and part2Iterations are how many times the art is enhanced
// in each part. The worked example in the puzzle text stops after 2
// iterations instead, and its two rules do not cover the patterns that come
// after that, so both parts have no answer for the example.
const (
	part1Iterations = 5
	part2Iterations = 18
)

// start is the pattern the program always begins with.
var start = pattern{".#.", "..#", "###"}

// pattern is a square of pixels, one string per row. A '#' is a lit pixel.
type pattern []string

// key joins the rows with slashes, the same way the rules write them.
func (p pattern) key() string {
	return strings.Join(p, "/")
}

// rules maps each input pattern, in every rotation and flip, to its output.
type rules map[string]pattern

// memoKey is a 3x3 pattern and the number of iterations still to run.
type memoKey struct {
	key   string
	iters int
}

// artist counts lit pixels and remembers the counts it has found already.
type artist struct {
	rules rules
	memo  map[memoKey]int
}

// part1 counts the lit pixels after 5 iterations.
func part1(in string) any {
	return lit(in, part1Iterations)
}

// part2 counts the lit pixels after 18 iterations.
func part2(in string) any {
	return lit(in, part2Iterations)
}

// lit reads the rules and counts the lit pixels after iters iterations. It
// returns nil when the rules do not cover a pattern that comes up.
func lit(in string, iters int) any {
	a := artist{rules: parse(in), memo: map[memoKey]int{}}

	n, ok := a.count(start, iters)
	if !ok {
		return nil
	}

	return n
}

// parse reads the rules. Each input pattern is stored in all eight of its
// rotations and flips, so a lookup needs only the square as it is.
func parse(in string) rules {
	r := rules{}

	for _, line := range input.Lines(in) {
		from, to, _ := strings.Cut(line, " => ")
		p := pattern(strings.Split(from, "/"))
		out := pattern(strings.Split(to, "/"))

		for range 4 {
			r[p.key()] = out
			r[flip(p).key()] = out
			p = rotate(p)
		}
	}

	return r
}

// rotate turns a pattern 90 degrees clockwise.
func rotate(p pattern) pattern {
	n := len(p)
	out := make(pattern, n)

	for y := range n {
		row := make([]byte, n)
		for x := range n {
			row[x] = p[n-1-x][y]
		}

		out[y] = string(row)
	}

	return out
}

// flip mirrors a pattern from left to right.
func flip(p pattern) pattern {
	out := make(pattern, len(p))
	for i, row := range p {
		out[i] = strx.Reverse(row)
	}

	return out
}

// count returns the number of lit pixels that the 3x3 pattern p grows into
// after iters iterations.
//
// A 3x3 square becomes 4x4, then 6x6, then 9x9. In each of those steps the
// squares split along the edges of the squares from the step before, so no
// rule ever sees pixels from two of the original squares at once. The 9x9
// grid then splits into nine 3x3 squares, and each of them grows on its own
// in the same way. So count follows one 3x3 square for 3 iterations, then
// adds up the counts of its nine 3x3 squares. Only a few different squares
// come up, and the memo makes each of them cost almost nothing after the
// first time.
func (a *artist) count(p pattern, iters int) (int, bool) {
	if iters < 3 {
		g, ok := a.enhanceN(p, iters)
		if !ok {
			return 0, false
		}

		return litPixels(g), true
	}

	k := memoKey{p.key(), iters}
	if n, ok := a.memo[k]; ok {
		return n, true
	}

	g, ok := a.enhanceN(p, 3)
	if !ok {
		return 0, false
	}

	total := 0
	for _, sq := range split(g, 3) {
		n, ok := a.count(sq, iters-3)
		if !ok {
			return 0, false
		}

		total += n
	}

	a.memo[k] = total

	return total, true
}

// enhanceN runs n iterations on g.
func (a *artist) enhanceN(g pattern, n int) (pattern, bool) {
	for range n {
		var ok bool
		if g, ok = a.enhance(g); !ok {
			return nil, false
		}
	}

	return g, true
}

// enhance runs one iteration: it splits g into 2x2 squares when the size is
// even and into 3x3 squares otherwise, and replaces each square with the
// output of its rule. It returns false when a square has no rule.
func (a *artist) enhance(g pattern) (pattern, bool) {
	size := 3
	if len(g)%2 == 0 {
		size = 2
	}

	blocks := len(g) / size
	outSize := size + 1

	rows := make([][]byte, blocks*outSize)
	for i := range rows {
		rows[i] = make([]byte, 0, blocks*outSize)
	}

	for i, sq := range split(g, size) {
		out, ok := a.rules[sq.key()]
		if !ok {
			return nil, false
		}

		by := i / blocks
		for y, row := range out {
			rows[by*outSize+y] = append(rows[by*outSize+y], row...)
		}
	}

	next := make(pattern, len(rows))
	for i, row := range rows {
		next[i] = string(row)
	}

	return next, true
}

// split cuts g into squares of the given size, row by row from the top left.
func split(g pattern, size int) []pattern {
	blocks := len(g) / size

	var out []pattern
	for by := range blocks {
		for bx := range blocks {
			sq := make(pattern, size)
			for y := range size {
				sq[y] = g[by*size+y][bx*size : (bx+1)*size]
			}

			out = append(out, sq)
		}
	}

	return out
}

// litPixels counts the '#' characters in g.
func litPixels(g pattern) int {
	n := 0
	for _, row := range g {
		n += strings.Count(row, "#")
	}

	return n
}
