// Advent of Code 2022, day 21: Monkey Math.
// https://adventofcode.com/2022/day/21
package main

import (
	"strconv"
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 21, part1, part2)
}

// The names of the two special monkeys.
const (
	root  = "root"
	human = "humn"
)

// job is what one monkey yells. A monkey with no op yells value. A monkey with
// an op yells the result of "left op right".
type job struct {
	value       int
	left, right string
	op          byte
}

// part1 finds the number that the root monkey yells.
func part1(in string) any {
	return eval(parse(in), root)
}

// part2 finds the number that we (humn) must yell so that the two numbers of
// the root monkey are equal.
//
// humn is a leaf, and in the input only one path goes from root down to humn.
// So at each monkey on that path, one side depends on humn and the other side
// is a fixed number. We know the result that the monkey must give. We undo its
// operation with the fixed side to get the result that the humn side must
// give. We repeat this down to humn.
func part2(in string) any {
	jobs := parse(in)

	// root checks for equality, so its humn side must equal its other side.
	j := jobs[root]

	var name string
	var target int

	if hasHuman(jobs, j.left) {
		target = eval(jobs, j.right)
		name = j.left
	} else {
		target = eval(jobs, j.left)
		name = j.right
	}

	for name != human {
		j := jobs[name]

		if hasHuman(jobs, j.left) {
			target = solveLeft(j.op, target, eval(jobs, j.right))
			name = j.left
		} else {
			target = solveRight(j.op, target, eval(jobs, j.left))
			name = j.right
		}
	}

	return target
}

// solveLeft finds x in "x op right = target".
func solveLeft(op byte, target, right int) int {
	switch op {
	case '+':
		return target - right
	case '-':
		return target + right
	case '*':
		return target / right
	default:
		return target * right
	}
}

// solveRight finds x in "left op x = target". Subtraction and division do
// not commute, so they differ from solveLeft.
func solveRight(op byte, target, left int) int {
	switch op {
	case '+':
		return target - left
	case '-':
		return left - target
	case '*':
		return target / left
	default:
		return left / target
	}
}

// eval finds the number that the named monkey yells.
func eval(jobs map[string]job, name string) int {
	j := jobs[name]
	if j.op == 0 {
		return j.value
	}

	a, b := eval(jobs, j.left), eval(jobs, j.right)

	switch j.op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	default:
		return a / b
	}
}

// hasHuman tells if the number of the named monkey depends on humn.
func hasHuman(jobs map[string]job, name string) bool {
	if name == human {
		return true
	}

	j := jobs[name]
	if j.op == 0 {
		return false
	}

	return hasHuman(jobs, j.left) || hasHuman(jobs, j.right)
}

// parse reads lines like "root: pppw + sjmn" or "dbpl: 5" into a map from
// monkey name to job.
func parse(in string) map[string]job {
	jobs := map[string]job{}

	for _, line := range input.Lines(in) {
		name, rest, _ := strings.Cut(line, ": ")

		// A number job has one field. An operation job has three.
		parts := strings.Fields(rest)
		if len(parts) == 1 {
			v, _ := strconv.Atoi(parts[0])
			jobs[name] = job{value: v}
			continue
		}

		jobs[name] = job{left: parts[0], op: parts[1][0], right: parts[2]}
	}

	return jobs
}
