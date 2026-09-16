// Package slicesx extends the standard slices package with the generic helpers
// puzzles keep needing: sums, products, chunking, transposes, counting and
// combinatorics.
//
// The standard library already has slices.Sort, slices.Reverse, slices.Contains,
// slices.Max, slices.Min and slices.Index; use those directly.
package slicesx

import (
	"cmp"
	"iter"
)

// Number is any built-in integer or float type.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum adds up all the numbers.
func Sum[T Number](items []T) T {
	var total T

	for _, v := range items {
		total += v
	}

	return total
}

// Product multiplies all the numbers together. The product of nothing is 1.
func Product[T Number](items []T) T {
	total := T(1)

	for _, v := range items {
		total *= v
	}

	return total
}

// Map applies fn to every item and returns the results.
func Map[T, U any](items []T, fn func(T) U) []U {
	out := make([]U, len(items))

	for i, v := range items {
		out[i] = fn(v)
	}

	return out
}

// Filter returns the items for which pred is true.
func Filter[T any](items []T, pred func(T) bool) []T {
	var out []T

	for _, v := range items {
		if pred(v) {
			out = append(out, v)
		}
	}

	return out
}

// Count returns how many items satisfy pred.
func Count[T any](items []T, pred func(T) bool) int {
	n := 0

	for _, v := range items {
		if pred(v) {
			n++
		}
	}

	return n
}

// Counts tallies how many times each value appears.
func Counts[T comparable](items []T) map[T]int {
	out := make(map[T]int)

	for _, v := range items {
		out[v]++
	}

	return out
}

// MinMax returns the smallest and largest item. It panics on an empty slice.
func MinMax[T cmp.Ordered](items []T) (T, T) {
	lo, hi := items[0], items[0]

	for _, v := range items[1:] {
		lo = min(lo, v)
		hi = max(hi, v)
	}

	return lo, hi
}

// Chunk splits the slice into pieces of size n. The last piece may be shorter.
func Chunk[T any](items []T, n int) [][]T {
	var out [][]T

	for i := 0; i < len(items); i += n {
		out = append(out, items[i:min(i+n, len(items))])
	}

	return out
}

// Windows yields every run of n consecutive items. The windows share memory
// with the original slice, so copy one before you change it.
func Windows[T any](items []T, n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		for i := 0; i+n <= len(items); i++ {
			if !yield(items[i : i+n]) {
				return
			}
		}
	}
}

// Pairs yields every unordered pair of distinct indexes (i < j) as their items.
func Pairs[T any](items []T) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		for i := 0; i < len(items); i++ {
			for j := i + 1; j < len(items); j++ {
				if !yield(items[i], items[j]) {
					return
				}
			}
		}
	}
}

// Transpose turns rows into columns. All rows must have the same length.
func Transpose[T any](rows [][]T) [][]T {
	if len(rows) == 0 {
		return nil
	}

	out := make([][]T, len(rows[0]))
	for x := range out {
		out[x] = make([]T, len(rows))
		for y := range rows {
			out[x][y] = rows[y][x]
		}
	}

	return out
}

// Permutations yields every ordering of the items using Heap's algorithm.
// The yielded slice is reused between iterations, so copy it if you keep it.
func Permutations[T any](items []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		work := make([]T, len(items))
		copy(work, items)

		// Heap's algorithm generates each permutation from the previous one with
		// a single swap, which is why it is fast. c tracks the loop counter for
		// each position, replacing the recursive version's call stack.
		c := make([]int, len(work))

		if !yield(work) {
			return
		}

		for i := 0; i < len(work); {
			if c[i] < i {
				if i%2 == 0 {
					work[0], work[i] = work[i], work[0]
				} else {
					work[c[i]], work[i] = work[i], work[c[i]]
				}

				if !yield(work) {
					return
				}

				c[i]++
				i = 0
			} else {
				c[i] = 0
				i++
			}
		}
	}
}

// Combinations yields every way to choose k items, in order of the original
// slice. The yielded slice is reused between iterations, so copy it if you keep it.
func Combinations[T any](items []T, k int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if k > len(items) || k < 0 {
			return
		}

		// idx holds the chosen positions and starts as 0, 1, ..., k-1. Each step
		// advances the rightmost index that can still move and resets the ones
		// after it, like an odometer that must stay strictly increasing.
		idx := make([]int, k)
		for i := range idx {
			idx[i] = i
		}
		out := make([]T, k)

		for {
			for i, j := range idx {
				out[i] = items[j]
			}
			if !yield(out) {
				return
			}

			i := k - 1
			for i >= 0 && idx[i] == len(items)-k+i {
				i--
			}
			if i < 0 {
				return
			}

			idx[i]++
			for j := i + 1; j < k; j++ {
				idx[j] = idx[j-1] + 1
			}
		}
	}
}
