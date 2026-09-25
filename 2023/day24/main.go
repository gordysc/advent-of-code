// Advent of Code 2023, day 24: Never Tell Me The Odds.
// https://adventofcode.com/2023/day/24
package main

import (
	"math/big"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2023, 24, part1, part2)
}

// hailstone is a position and a velocity, each as X, Y and Z.
type hailstone struct {
	pos [3]int
	vel [3]int
}

// part1 counts the pairs of hailstones whose paths cross inside the test area
// in the future, when we look only at X and Y.
func part1(in string) any {
	stones := parse(in)
	lo, hi := testArea(stones)

	count := 0

	for i := range stones {
		for j := i + 1; j < len(stones); j++ {
			if crossInside(stones[i], stones[j], lo, hi) {
				count++
			}
		}
	}

	return count
}

// part2 finds the rock throw that hits every hailstone, and returns the sum of
// the X, Y and Z of the rock's start position.
//
// Let the rock start at P with velocity V. It hits hailstone i at some time t,
// so P - p_i = t (v_i - V). The two sides are parallel, so their cross product
// is zero: (P - p_i) x (V - v_i) = 0. Expand this, and the only term that is
// not linear is P x V. That term is the same for every hailstone. Thus when we
// subtract the equation of hailstone i from the equation of hailstone j, we
// get three linear equations in the six unknowns:
//
//	P x (v_j - v_i) + (p_j - p_i) x V = p_j x v_j - p_i x v_i
//
// Two such pairs give six equations. We solve them with exact fractions,
// because the coordinates are near 4e14 and the products overflow int64.
func part2(in string) any {
	stones := parse(in)

	// Most choices of three hailstones give a system with one solution. If a
	// choice does not (for example, two velocities are parallel), try the next.
	for j := 1; j < len(stones); j++ {
		for k := j + 1; k < len(stones); k++ {
			rows := append(equations(stones[0], stones[j]), equations(stones[0], stones[k])...)

			sol, ok := solve(rows)
			if !ok {
				continue
			}

			sum := new(big.Rat).Add(sol[0], sol[1])
			sum.Add(sum, sol[2])
			if !sum.IsInt() {
				return nil
			}

			return int(sum.Num().Int64())
		}
	}

	return nil
}

// parse reads one hailstone from each line.
func parse(in string) []hailstone {
	var stones []hailstone

	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		stones = append(stones, hailstone{
			pos: [3]int{n[0], n[1], n[2]},
			vel: [3]int{n[3], n[4], n[5]},
		})
	}

	return stones
}

// testArea returns the limits of the test area for X and Y.
//
// The example uses 7 to 27 and the real input uses 200000000000000 to
// 400000000000000. The puzzle gives no area in the input, so we choose by the
// size of the coordinates: the example has only small numbers.
func testArea(stones []hailstone) (float64, float64) {
	for _, s := range stones {
		for _, c := range s.pos {
			if mathx.Abs(c) >= 1000 {
				return 200000000000000, 400000000000000
			}
		}
	}

	return 7, 27
}

// crossInside tells if the X-Y paths of a and b cross in the future of both
// hailstones, at a point inside the test area.
//
// The paths are a.pos + t*a.vel and b.pos + s*b.vel. Take the 2D cross product
// of both sides with b.vel, and then with a.vel, to get t and s as fractions
// with the same denominator d. The signs of t and s come from exact integer
// math. Only the position check uses floats, and there an error of a few
// units cannot change the result.
func crossInside(a, b hailstone, lo, hi float64) bool {
	d := a.vel[0]*b.vel[1] - a.vel[1]*b.vel[0]
	if d == 0 {
		// The paths are parallel, so they never cross.
		return false
	}

	dx := b.pos[0] - a.pos[0]
	dy := b.pos[1] - a.pos[1]
	tNum := dx*b.vel[1] - dy*b.vel[0]
	sNum := dx*a.vel[1] - dy*a.vel[0]

	// A negative time means the paths crossed in the past.
	if mathx.Sign(tNum)*mathx.Sign(d) < 0 || mathx.Sign(sNum)*mathx.Sign(d) < 0 {
		return false
	}

	t := float64(tNum) / float64(d)
	x := float64(a.pos[0]) + t*float64(a.vel[0])
	y := float64(a.pos[1]) + t*float64(a.vel[1])

	return lo <= x && x <= hi && lo <= y && y <= hi
}

// equations returns the three linear equations that hailstones a and b give
// for the rock. Each row holds the factors of Px, Py, Pz, Vx, Vy, Vz and then
// the right side.
func equations(a, b hailstone) [][]*big.Rat {
	var u, w [3]int
	for i := range 3 {
		u[i] = b.pos[i] - a.pos[i]
		w[i] = b.vel[i] - a.vel[i]
	}

	// rhs = b.pos x b.vel - a.pos x a.vel.
	ca, cb := cross(a.pos, a.vel), cross(b.pos, b.vel)
	var rhs [3]*big.Rat
	for i := range 3 {
		rhs[i] = new(big.Rat).SetInt(new(big.Int).Sub(cb[i], ca[i]))
	}

	// These rows are the X, Y and Z parts of P x w + u x V.
	coeffs := [3][6]int{
		{0, w[2], -w[1], 0, -u[2], u[1]},
		{-w[2], 0, w[0], u[2], 0, -u[0]},
		{w[1], -w[0], 0, -u[1], u[0], 0},
	}

	rows := make([][]*big.Rat, 3)
	for r, cs := range coeffs {
		for _, c := range cs {
			rows[r] = append(rows[r], big.NewRat(int64(c), 1))
		}

		rows[r] = append(rows[r], rhs[r])
	}

	return rows
}

// cross returns the cross product a x b. It uses big integers because each
// product of a coordinate and a velocity can be near the int64 limit.
func cross(a, b [3]int) [3]*big.Int {
	mul := func(x, y int) *big.Int {
		return new(big.Int).Mul(big.NewInt(int64(x)), big.NewInt(int64(y)))
	}

	return [3]*big.Int{
		new(big.Int).Sub(mul(a[1], b[2]), mul(a[2], b[1])),
		new(big.Int).Sub(mul(a[2], b[0]), mul(a[0], b[2])),
		new(big.Int).Sub(mul(a[0], b[1]), mul(a[1], b[0])),
	}
}

// solve solves a square linear system by Gauss-Jordan elimination. Each row
// holds the factors and then the right side. It returns false when the system
// has no single solution.
func solve(rows [][]*big.Rat) ([]*big.Rat, bool) {
	n := len(rows)

	for col := range n {
		pivot := -1
		for r := col; r < n; r++ {
			if rows[r][col].Sign() != 0 {
				pivot = r

				break
			}
		}

		if pivot < 0 {
			return nil, false
		}

		rows[col], rows[pivot] = rows[pivot], rows[col]

		// Remove this column from every other row.
		for r := range n {
			if r == col || rows[r][col].Sign() == 0 {
				continue
			}

			f := new(big.Rat).Quo(rows[r][col], rows[col][col])
			for c := col; c <= n; c++ {
				step := new(big.Rat).Mul(f, rows[col][c])
				rows[r][c] = new(big.Rat).Sub(rows[r][c], step)
			}
		}
	}

	sol := make([]*big.Rat, n)
	for i := range n {
		sol[i] = new(big.Rat).Quo(rows[i][n], rows[i][i])
	}

	return sol, true
}
