// Package mathx holds the integer maths that puzzles ask for again and again.
//
// Go's math package works on float64, so integer helpers such as GCD and
// integer power have to be written by hand. Note that Go 1.21+ has built-in
// min and max, so those are not repeated here.
package mathx

// Abs returns the absolute value of n.
func Abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

// Sign returns -1, 0 or 1 depending on the sign of n.
func Sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	default:
		return 0
	}
}

// GCD returns the greatest common divisor using Euclid's algorithm.
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return Abs(a)
}

// LCM returns the least common multiple of all the given numbers. Cycle
// puzzles ("when do all the ghosts arrive at once?") almost always need this.
func LCM(nums ...int) int {
	result := 1

	for _, n := range nums {
		// Divide first so the intermediate value stays small.
		result = result / GCD(result, n) * n
	}

	return result
}

// Mod returns a modulo m with a result in [0, m), unlike Go's % operator which
// keeps the sign of a. Use it for wrapping indexes and coordinates.
func Mod(a, m int) int {
	r := a % m
	if r < 0 {
		r += m
	}

	return r
}

// Pow returns base raised to exp for integers, using exponentiation by squaring.
func Pow(base, exp int) int {
	result := 1

	for exp > 0 {
		if exp&1 == 1 {
			result *= base
		}
		base *= base
		exp >>= 1
	}

	return result
}

// DivCeil returns a / b rounded up. Both values must be positive.
func DivCeil(a, b int) int {
	return (a + b - 1) / b
}

// Triangular returns 1 + 2 + ... + n.
func Triangular(n int) int {
	return n * (n + 1) / 2
}

// Digits returns the decimal digits of n from most to least significant.
func Digits(n int) []int {
	if n == 0 {
		return []int{0}
	}

	n = Abs(n)

	var out []int
	for n > 0 {
		out = append([]int{n % 10}, out...)
		n /= 10
	}

	return out
}

// NumDigits returns how many decimal digits n has.
func NumDigits(n int) int {
	n = Abs(n)
	count := 1

	for n >= 10 {
		n /= 10
		count++
	}

	return count
}
