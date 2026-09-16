package slicesx

import (
	"reflect"
	"testing"
)

// TestAggregates checks the arithmetic helpers.
func TestAggregates(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5}

	if Sum(nums) != 14 || Product(nums) != 60 {
		t.Fatal("Sum/Product wrong")
	}
	if lo, hi := MinMax(nums); lo != 1 || hi != 5 {
		t.Fatal("MinMax wrong")
	}
	if Counts(nums)[1] != 2 {
		t.Fatal("Counts wrong")
	}
	if Count(nums, func(n int) bool { return n > 2 }) != 3 {
		t.Fatal("Count wrong")
	}
}

// TestShapes checks chunking, windows and transposition.
func TestShapes(t *testing.T) {
	if got := Chunk([]int{1, 2, 3, 4, 5}, 2); !reflect.DeepEqual(got, [][]int{{1, 2}, {3, 4}, {5}}) {
		t.Fatalf("Chunk = %v", got)
	}

	var windows [][]int
	for w := range Windows([]int{1, 2, 3}, 2) {
		windows = append(windows, append([]int(nil), w...))
	}
	if !reflect.DeepEqual(windows, [][]int{{1, 2}, {2, 3}}) {
		t.Fatalf("Windows = %v", windows)
	}

	if got := Transpose([][]int{{1, 2}, {3, 4}}); !reflect.DeepEqual(got, [][]int{{1, 3}, {2, 4}}) {
		t.Fatalf("Transpose = %v", got)
	}
}

// TestCombinatorics checks the counts and contents of the generators.
func TestCombinatorics(t *testing.T) {
	seen := map[string]bool{}
	for p := range Permutations([]byte("abc")) {
		seen[string(p)] = true
	}
	if len(seen) != 6 || !seen["cab"] {
		t.Fatalf("Permutations = %v", seen)
	}

	var combos []string
	for c := range Combinations([]byte("abcd"), 2) {
		combos = append(combos, string(c))
	}
	if !reflect.DeepEqual(combos, []string{"ab", "ac", "ad", "bc", "bd", "cd"}) {
		t.Fatalf("Combinations = %v", combos)
	}

	pairs := 0
	for range Pairs([]int{1, 2, 3, 4}) {
		pairs++
	}
	if pairs != 6 {
		t.Fatalf("Pairs = %d", pairs)
	}
}
