package strx

import "testing"

// TestHelpers runs simple cases through each function.
func TestHelpers(t *testing.T) {
	if Reverse("héllo") != "olléh" {
		t.Fatal("Reverse wrong")
	}
	if Between("x=[42]", "[", "]") != "42" || Between("none", "[", "]") != "" {
		t.Fatal("Between wrong")
	}
	if Hamming("abcd", "abxd") != 1 {
		t.Fatal("Hamming wrong")
	}
	if got := Chunk("abcde", 2); len(got) != 3 || got[2] != "e" {
		t.Fatalf("Chunk = %v", got)
	}
}
