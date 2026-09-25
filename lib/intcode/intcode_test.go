package intcode

import (
	"slices"
	"testing"
)

// TestMemory checks the day 2 examples, which only add and multiply.
func TestMemory(t *testing.T) {
	cases := []struct {
		program string
		want    []int
	}{
		{"1,0,0,0,99", []int{2, 0, 0, 0, 99}},
		{"2,3,0,3,99", []int{2, 3, 0, 6, 99}},
		{"2,4,4,5,99,0", []int{2, 4, 4, 5, 99, 9801}},
		{"1,1,1,4,99,5,6,0,99", []int{30, 1, 1, 4, 2, 5, 6, 0, 99}},
	}

	for _, c := range cases {
		m := New(Parse(c.program))
		m.Run()

		if !m.Halted() || !slices.Equal(m.mem, c.want) {
			t.Errorf("%s: memory = %v, want %v", c.program, m.mem, c.want)
		}
	}
}

// TestCompare checks the day 5 example that outputs 999, 1000 or 1001 for
// input below, equal to or above 8. It uses jumps, compares and both modes.
func TestCompare(t *testing.T) {
	program := Parse("3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31," +
		"1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104," +
		"999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99")

	for in, want := range map[int]int{7: 999, 8: 1000, 9: 1001} {
		m := New(program)
		m.Send(in)

		if got := m.Run(); !slices.Equal(got, []int{want}) {
			t.Errorf("input %d: output %v, want [%d]", in, got, want)
		}
	}
}

// TestRelative checks the day 9 examples: a quine that uses the relative
// base and memory past the program, and a large number.
func TestRelative(t *testing.T) {
	quine := Parse("109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99")
	if got := New(quine).Run(); !slices.Equal(got, quine) {
		t.Errorf("quine output %v", got)
	}

	if got := New(Parse("104,1125899906842624,99")).Run(); got[0] != 1125899906842624 {
		t.Errorf("large output %v", got)
	}
}

// TestWaiting checks that Run stops for input and carries on after Send.
func TestWaiting(t *testing.T) {
	m := New(Parse("3,9,4,9,3,9,4,9,99,0"))

	if out := m.Run(); len(out) != 0 || !m.Waiting() {
		t.Fatalf("first run: output %v, waiting %v", out, m.Waiting())
	}

	m.Send(5)
	if out := m.Run(); !slices.Equal(out, []int{5}) || !m.Waiting() {
		t.Fatalf("second run: output %v, waiting %v", out, m.Waiting())
	}

	c := m.Clone()
	c.Send(7)
	m.Send(6)

	if out := m.Run(); !slices.Equal(out, []int{6}) || !m.Halted() {
		t.Fatalf("third run: output %v, halted %v", out, m.Halted())
	}

	if out := c.Run(); !slices.Equal(out, []int{7}) {
		t.Fatalf("clone run: output %v", out)
	}
}
