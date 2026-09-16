package input

import (
	"reflect"
	"testing"
)

// TestLines checks the trailing-newline handling.
func TestLines(t *testing.T) {
	got := Lines("a\nb\n")
	want := []string{"a", "b"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Lines = %v, want %v", got, want)
	}
}

// TestBlocks checks the blank-line grouping.
func TestBlocks(t *testing.T) {
	got := Blocks("1\n2\n\n3\n")
	want := [][]string{{"1", "2"}, {"3"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Blocks = %v, want %v", got, want)
	}
}

// TestInts checks signed and unsigned extraction.
func TestInts(t *testing.T) {
	if got := Ints("move 3 from -4 to 5"); !reflect.DeepEqual(got, []int{3, -4, 5}) {
		t.Fatalf("Ints = %v", got)
	}
	if got := UInts("3-4"); !reflect.DeepEqual(got, []int{3, 4}) {
		t.Fatalf("UInts = %v", got)
	}
	if got := Digits("a1b23"); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Digits = %v", got)
	}
}
