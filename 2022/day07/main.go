// Advent of Code 2022, day 7: No Space Left On Device.
// https://adventofcode.com/2022/day/7
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 7, part1, part2)
}

// part1 gives the sum of the total sizes of all directories that have a
// total size of at most 100000. A file inside nested small directories
// counts one time for each of these directories.
func part1(in string) any {
	total := 0

	for _, size := range dirSizes(in) {
		if size <= 100000 {
			total += size
		}
	}

	return total
}

// part2 finds the smallest directory that frees enough space when we delete
// it. The disk has 70000000 units and the update needs 30000000 free units.
func part2(in string) any {
	sizes := dirSizes(in)
	free := 70000000 - sizes["/"]
	need := 30000000 - free
	best := sizes["/"]

	for _, size := range sizes {
		if size >= need && size < best {
			best = size
		}
	}

	return best
}

// dirSizes replays the terminal output and gives the total size of each
// directory, keyed by its full path.
//
// A stack of paths holds the current directory and all of its parents. When
// the output lists a file, its size goes to every directory on that stack,
// because each of them contains the file. The "ls" and "dir" lines give no
// size, so the loop skips them.
func dirSizes(in string) map[string]int {
	sizes := map[string]int{}
	path := []string{"/"}

	for _, line := range input.Lines(in) {
		fields := strings.Fields(line)

		switch {
		case line == "$ cd /":
			path = path[:1]

		case line == "$ cd ..":
			path = path[:len(path)-1]

		case fields[0] == "$" && fields[1] == "cd":
			// Build the full path so that two directories with the same
			// name in different places do not share a total.
			path = append(path, path[len(path)-1]+fields[2]+"/")

		case fields[0] == "$" || fields[0] == "dir":
			continue

		default:
			size := input.Int(fields[0])

			for _, dir := range path {
				sizes[dir] += size
			}
		}
	}

	return sizes
}
