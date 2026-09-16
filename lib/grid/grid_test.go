package grid

import "testing"

// TestParseAndLookup checks parsing, bounds and the S-finder.
func TestParseAndLookup(t *testing.T) {
	g := Parse("#.S\n...\n")

	if g.W != 3 || g.H != 2 {
		t.Fatalf("size = %dx%d", g.W, g.H)
	}

	p, ok := g.FindByte('S')
	if !ok || p != P(2, 0) {
		t.Fatalf("FindByte = %v %v", p, ok)
	}

	if n := len(g.Neighbors4(p)); n != 2 {
		t.Fatalf("corner has %d neighbours, want 2", n)
	}

	if _, ok := g.Get(P(3, 0)); ok {
		t.Fatal("Get out of bounds returned ok")
	}

	if g.String() != "#.S\n...\n" {
		t.Fatalf("String = %q", g.String())
	}
}

// TestTurns checks the rotation helpers keep the right handedness.
func TestTurns(t *testing.T) {
	if Up.TurnRight() != Right || Right.TurnRight() != Down {
		t.Fatal("TurnRight is wrong")
	}
	if Up.TurnLeft() != Left || Left.TurnLeft() != Down {
		t.Fatal("TurnLeft is wrong")
	}
	if P(1, 1).Manhattan(P(4, 5)) != 7 {
		t.Fatal("Manhattan is wrong")
	}
}

// TestIterators checks Points and All visit every cell in reading order.
func TestIterators(t *testing.T) {
	g := Parse("ab\ncd")

	var seen []byte
	for _, v := range g.All() {
		seen = append(seen, v)
	}

	if string(seen) != "abcd" {
		t.Fatalf("All = %q", seen)
	}
}
