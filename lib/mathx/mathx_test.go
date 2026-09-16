package mathx

import (
	"reflect"
	"testing"
)

// TestBasics runs a table of simple cases through the helpers.
func TestBasics(t *testing.T) {
	if GCD(12, 18) != 6 || LCM(4, 6, 10) != 60 {
		t.Fatal("GCD/LCM wrong")
	}
	if Mod(-1, 5) != 4 || Mod(7, 5) != 2 {
		t.Fatal("Mod wrong")
	}
	if Pow(2, 10) != 1024 || Pow(7, 0) != 1 {
		t.Fatal("Pow wrong")
	}
	if DivCeil(7, 2) != 4 || Triangular(4) != 10 {
		t.Fatal("DivCeil/Triangular wrong")
	}
	if !reflect.DeepEqual(Digits(1203), []int{1, 2, 0, 3}) || NumDigits(1203) != 4 {
		t.Fatal("Digits wrong")
	}
}
