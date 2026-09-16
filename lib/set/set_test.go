package set

import "testing"

// TestSetOps checks the set algebra.
func TestSetOps(t *testing.T) {
	a := Of(1, 2, 3)
	b := Of(3, 4)

	if a.Union(b).Len() != 4 {
		t.Fatal("Union wrong")
	}
	if !a.Intersect(b).Has(3) || a.Intersect(b).Len() != 1 {
		t.Fatal("Intersect wrong")
	}
	if a.Difference(b).Has(3) || a.Difference(b).Len() != 2 {
		t.Fatal("Difference wrong")
	}

	a.Remove(1)
	if a.Has(1) || a.Len() != 2 {
		t.Fatal("Remove wrong")
	}
}
