// Advent of Code 2019, day 17: Set and Forget.
// https://adventofcode.com/2019/day/17
//
// The puzzle text has no example program, only pictures of camera views.
// example.txt holds a small handmade Intcode program that prints the part 1
// example view as ASCII and halts, so part 1 gives 76 for it. Its first
// instruction adds two zeros, like the real program's first instruction, so
// the wake-up write of 2 to address 0 makes it a harmless multiply. The
// program then has no movement prompt, so part 2 has no answer for it.
package main

import (
	"slices"
	"strconv"
	"strings"

	"aoc/lib/grid"
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 17, part1, part2)
}

// maxLen is the most characters the robot's memory holds for the main
// routine and for each movement function, not counting the newline.
const maxLen = 20

// maxFuncs is how many movement functions the robot has: A, B and C.
const maxFuncs = 3

// part1 reads the camera view and adds up the alignment parameters of the
// scaffold intersections.
func part1(in string) any {
	view := intcode.New(intcode.Parse(in)).RunString()

	return alignment(view)
}

// part2 wakes the robot up, sends it a movement routine that visits every
// part of the scaffold, and returns the amount of dust it collects.
func part2(in string) any {
	program := intcode.Parse(in)

	view := intcode.New(program).RunString()
	routine, funcs, ok := compress(path(view))
	if !ok {
		return nil
	}

	// Writing 2 to address 0 wakes the robot up. The program then prints the
	// view again and asks for the main routine, so it must stop and wait for
	// input. A program that halts instead has no robot to drive.
	robot := intcode.New(program)
	robot.Poke(0, 2)
	robot.Run()
	if !robot.Waiting() {
		return nil
	}

	robot.SendLine(routine)
	for _, f := range funcs {
		robot.SendLine(f)
	}
	robot.SendLine("n") // no continuous video feed

	// The robot prints its prompts and the view as ASCII text first. The dust
	// amount comes last, and it is the only value too large to be a
	// character. A last value that is a character means the robot failed.
	out := robot.Run()
	if len(out) == 0 || out[len(out)-1] < 128 {
		return nil
	}

	return out[len(out)-1]
}

// alignment finds every scaffold cell with scaffold on all four sides and
// adds up x*y for each one.
func alignment(view string) int {
	rows := strings.Fields(view)

	sum := 0
	for y, row := range rows {
		for x := range row {
			p := grid.P(x, y)
			if !isScaffold(rows, p) {
				continue
			}

			crossing := true
			for _, n := range p.Neighbors4() {
				crossing = crossing && isScaffold(rows, n)
			}

			if crossing {
				sum += x * y
			}
		}
	}

	return sum
}

// path follows the scaffold from the robot to the far end and returns the
// moves. A move is a turn and a number of steps, such as "R,8". The robot
// turns only at corners, and it goes straight on across every intersection,
// which is the path that visits all of the scaffold on these maps.
func path(view string) []string {
	rows := strings.Fields(view)

	pos, dir, found := robot(rows)
	if !found {
		return nil
	}

	// The robot can face along the scaffold at the start. Then the first
	// move is only a number of steps, with no turn.
	turn := ""
	if !isScaffold(rows, pos.Add(dir)) {
		turn, dir = nextTurn(rows, pos, dir)
		if turn == "" {
			return nil
		}
	}

	var moves []string
	for {
		steps := 0
		for isScaffold(rows, pos.Add(dir)) {
			pos = pos.Add(dir)
			steps++
		}

		move := strconv.Itoa(steps)
		if turn != "" {
			move = turn + "," + move
		}
		moves = append(moves, move)

		turn, dir = nextTurn(rows, pos, dir)
		if turn == "" {
			break
		}
	}

	return moves
}

// nextTurn looks left and right of the robot for more scaffold. It returns
// the turn ("L" or "R") and the new direction, or an empty turn at the end of
// the scaffold.
func nextTurn(rows []string, pos, dir grid.Point) (string, grid.Point) {
	if isScaffold(rows, pos.Add(dir.TurnLeft())) {
		return "L", dir.TurnLeft()
	}

	if isScaffold(rows, pos.Add(dir.TurnRight())) {
		return "R", dir.TurnRight()
	}

	return "", dir
}

// robot finds the robot in the view and the direction it faces.
func robot(rows []string) (grid.Point, grid.Point, bool) {
	for y, row := range rows {
		for x, r := range row {
			if dir, ok := grid.DirFromRune[r]; ok && strings.ContainsRune("^v<>", r) {
				return grid.P(x, y), dir, true
			}
		}
	}

	return grid.Point{}, grid.Point{}, false
}

// isScaffold reports whether p is on the scaffold. The robot always stands on
// the scaffold, so its cell counts too. Cells outside the view do not.
func isScaffold(rows []string, p grid.Point) bool {
	if p.Y < 0 || p.Y >= len(rows) || p.X < 0 || p.X >= len(rows[p.Y]) {
		return false
	}

	return strings.IndexByte("#^v<>", rows[p.Y][p.X]) >= 0
}

// compress splits the moves into a main routine of calls to A, B and C, and
// the three movement functions, each at most maxLen characters. It returns
// false if no split fits.
//
// It is a depth-first search. At each point in the moves it first tries every
// function that is already defined and matches there. If fewer than three are
// defined, it then tries every length for a new function that starts there.
// A dead end undoes the last choice and tries the next one. The search works
// on whole moves, so a function never splits a turn from its steps.
func compress(moves []string) (string, []string, bool) {
	var funcs [][]string
	var calls []string

	// solve is a closure, so it can read and change funcs and calls from
	// compress. It has to be declared before it is set, because it calls
	// itself.
	var solve func(rest []string) bool
	solve = func(rest []string) bool {
		if len(strings.Join(calls, ",")) > maxLen {
			return false
		}

		if len(rest) == 0 {
			return true
		}

		for i, f := range funcs {
			if len(f) > len(rest) || !slices.Equal(f, rest[:len(f)]) {
				continue
			}

			calls = append(calls, string(rune('A'+i)))
			if solve(rest[len(f):]) {
				return true
			}
			calls = calls[:len(calls)-1]
		}

		if len(funcs) == maxFuncs {
			return false
		}

		// rest[:k] shares its array with moves. That is safe here because
		// nothing ever writes to the moves.
		for k := 1; k <= len(rest) && len(strings.Join(rest[:k], ",")) <= maxLen; k++ {
			funcs = append(funcs, rest[:k])
			calls = append(calls, string(rune('A'+len(funcs)-1)))
			if solve(rest[k:]) {
				return true
			}
			funcs = funcs[:len(funcs)-1]
			calls = calls[:len(calls)-1]
		}

		return false
	}

	if len(moves) == 0 || !solve(moves) {
		return "", nil, false
	}

	// The robot asks for all three functions, even when fewer are needed.
	// The main routine never calls the spare ones, so a copy of A will do.
	for len(funcs) < maxFuncs {
		funcs = append(funcs, funcs[0])
	}

	texts := make([]string, len(funcs))
	for i, f := range funcs {
		texts[i] = strings.Join(f, ",")
	}

	return strings.Join(calls, ","), texts, true
}
