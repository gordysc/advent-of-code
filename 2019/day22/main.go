// Advent of Code 2019, day 22: Slam Shuffle.
// https://adventofcode.com/2019/day/22
package main

import (
	"math/bits"
	"strings"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 22, part1, part2)
}

// smallDeck and trackedCard are the deck size and the card that part 1 asks
// about. The examples in the puzzle text use a deck of 10 cards and show the
// whole deck order instead, so part 1 on example.txt only tells where card
// 2019 of a 10007-card deck goes after the example shuffle (it gives 1219).
const (
	smallDeck   = 10007
	trackedCard = 2019
)

// bigDeck, repeats and trackedPos are the numbers part 2 asks about: the card
// at position 2020 of a 119315717514047-card deck, after the whole shuffle is
// done 101741582076661 times. Both deck sizes are prime, which part 2 needs.
// On example.txt part 2 gives 117607927195067, which the puzzle text does not
// check.
const (
	bigDeck    = 119315717514047
	repeats    = 101741582076661
	trackedPos = 2020
)

// linear is the map pos -> (a*pos + b) mod n. It tells where the card at pos
// goes. Each shuffle technique is a map like this, and so is any sequence of
// them, because putting one linear map into another gives a linear map.
type linear struct {
	a, b int
}

// part1 finds the position of card 2019 after one shuffle of a 10007-card
// deck. In a new deck card c is at position c, so its new position is f(c).
func part1(in string) any {
	return parse(in, smallDeck).apply(trackedCard, smallDeck)
}

// part2 finds the card that ends up at position 2020 of the huge deck after
// the shuffle is repeated a huge number of times.
func part2(in string) any {
	return cardAt(in, bigDeck, repeats, trackedPos)
}

// cardAt returns the card at position pos after the shuffle is repeated
// times times on a deck of n cards, where n is prime.
//
// The repeated shuffle is one linear map f(c) = a*c + b, found with
// [linear.pow]. We need the card c with f(c) = pos, so c = (pos - b) / a.
// Division mod n means multiplication by the inverse of a: the number x with
// a*x = 1 (mod n). Fermat's little theorem says a^(n-1) = 1 (mod n) when n
// is prime and a is not 0 mod n, so a^(n-2) is that inverse.
func cardAt(in string, n, times, pos int) int {
	f := parse(in, n).pow(times, n)
	inv := powMod(f.a, n-2, n)

	return mulMod(mathx.Mod(pos-f.b, n), inv, n)
}

// parse reads the shuffle for a deck of n cards and returns its linear map.
// Each technique moves the card at position p to a new position:
//
//   - deal into new stack reverses the deck: p -> n-1-p, which is -p - 1.
//   - cut k moves the top k cards to the bottom: p -> p - k. A negative k
//     works the same way.
//   - deal with increment k puts the card from position p at k*p: p -> k*p.
func parse(in string, n int) linear {
	f := linear{a: 1, b: 0}

	for _, line := range input.Lines(in) {
		var step linear

		switch {
		case line == "deal into new stack":
			step = linear{a: n - 1, b: n - 1}
		case strings.HasPrefix(line, "cut "):
			k := input.Int(strings.TrimPrefix(line, "cut "))
			step = linear{a: 1, b: mathx.Mod(-k, n)}
		default:
			k := input.Int(strings.TrimPrefix(line, "deal with increment "))
			step = linear{a: mathx.Mod(k, n), b: 0}
		}

		f = f.then(step, n)
	}

	return f
}

// apply returns f(x) mod n.
func (f linear) apply(x, n int) int {
	return (mulMod(f.a, x, n) + f.b) % n
}

// then returns the map that does f first and g after it. Putting f(x) into g
// gives g.a*(f.a*x + f.b) + g.b, so the new a is g.a*f.a and the new b is
// g.a*f.b + g.b.
func (f linear) then(g linear, n int) linear {
	return linear{
		a: mulMod(g.a, f.a, n),
		b: (mulMod(g.a, f.b, n) + g.b) % n,
	}
}

// pow returns f done k times in a row. Like [powMod] for numbers, it uses
// exponentiation by squaring: it walks the bits of k and doubles the number
// of shuffles in sq at each step, so it needs about 47 steps, not 10^14.
// result and sq are both f done some number of times, so it does not matter
// which of the two goes first.
func (f linear) pow(k, n int) linear {
	result := linear{a: 1, b: 0}
	sq := f

	for k > 0 {
		if k&1 == 1 {
			result = result.then(sq, n)
		}

		sq = sq.then(sq, n)
		k >>= 1
	}

	return result
}

// powMod returns base^exp mod n by exponentiation by squaring.
func powMod(base, exp, n int) int {
	result := 1

	for exp > 0 {
		if exp&1 == 1 {
			result = mulMod(result, base, n)
		}

		base = mulMod(base, base, n)
		exp >>= 1
	}

	return result
}

// mulMod returns a*b mod n for a and b in [0, n).
//
// The part 2 deck size is about 2^47, so a*b can be about 2^94, which does
// not fit in 64 bits: a plain a*b % n would overflow and give a wrong answer.
// bits.Mul64 returns the full 128-bit product as two 64-bit halves, hi and
// lo. bits.Rem64 then divides that 128-bit number by n and returns the
// remainder, which is less than n and so fits in an int again.
func mulMod(a, b, n int) int {
	hi, lo := bits.Mul64(uint64(a), uint64(b))

	return int(bits.Rem64(hi, lo, uint64(n)))
}
