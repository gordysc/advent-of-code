// Package runner runs one day's puzzle solution and reports how long each part takes.
//
// Every day is its own `package main` that calls [Run] with its two part functions.
// The runner handles everything else: it finds and reads the input file, parses the
// command-line flags, times each part, and prints a festive report.
//
// Flags (all optional):
//
//	-input PATH   read this file instead of input.txt next to the solution
//	-example      read example.txt next to the solution instead of input.txt
//	-part N       run only part 1 or part 2
//	-runs N       run each part N times and report the best and average time
package runner

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Part is the signature of a puzzle part. It receives the whole input file as a
// string, with the single trailing newline removed, and returns the answer. The
// answer can be any value; it is printed with %v. Return nil to mark the part as
// not solved yet.
type Part func(input string) any

// result holds what we learned from timing one part.
type result struct {
	answer any
	best   time.Duration // fastest single run
	avg    time.Duration // mean over all runs
	runs   int
}

// options holds the parsed command-line flags.
type options struct {
	inputPath string
	example   bool
	part      int
	runs      int
}

// Run is the entry point every day's main function calls.
//
// year and day identify the puzzle for the report header. part1 and part2 are
// the solutions; pass nil for a part you have not written yet.
func Run(year, day int, part1, part2 Part) {
	opts := parseFlags()

	// Find the directory of the file that called Run. This lets `go run ./2015/day01`
	// find its own input.txt no matter what the working directory is.
	// runtime.Caller(1) means "one stack frame up from here", which is the caller's main.
	_, callerFile, _, ok := runtime.Caller(1)
	if !ok {
		fail("could not find the solution directory")
	}
	dir := filepath.Dir(callerFile)

	inputPath, input, err := loadInput(dir, opts)
	if err != nil {
		fail(err.Error())
	}

	printHeader(year, day, inputPath, opts.runs)

	var total time.Duration
	total += runPart(1, part1, input, opts)
	total += runPart(2, part2, input, opts)

	printFooter(total)
}

// parseFlags reads the command-line flags into an options struct.
func parseFlags() options {
	var opts options

	flag.StringVar(&opts.inputPath, "input", "", "path to the puzzle input file")
	flag.BoolVar(&opts.example, "example", false, "use example.txt instead of input.txt")
	flag.IntVar(&opts.part, "part", 0, "run only this part (1 or 2); 0 runs both")
	flag.IntVar(&opts.runs, "runs", 1, "number of timed runs per part")
	flag.Parse()

	if opts.part < 0 || opts.part > 2 {
		fail("-part must be 1 or 2")
	}
	if opts.runs < 1 {
		opts.runs = 1
	}

	return opts
}

// loadInput decides which file to read and reads it. It returns the path it used
// so the header can show it, and the file contents with one trailing newline removed.
func loadInput(dir string, opts options) (string, string, error) {
	path := opts.inputPath
	if path == "" {
		name := "input.txt"
		if opts.example {
			name = "example.txt"
		}
		path = filepath.Join(dir, name)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// errors.Is matches wrapped errors too, so this works no matter how the
		// os package built the error.
		if errors.Is(err, os.ErrNotExist) {
			return "", "", fmt.Errorf("no input file at %s\n   💡 run `make input YEAR=<year> DAY=<day>` to download it, or create the file by hand", path)
		}
		return "", "", fmt.Errorf("could not read %s: %w", path, err)
	}

	if len(data) == 0 {
		return "", "", fmt.Errorf("the input file %s is empty", path)
	}

	// Most puzzles end with a newline that solutions never want to see.
	// Only one is removed so inputs that end in a blank line stay intact.
	return path, strings.TrimSuffix(string(data), "\n"), nil
}

// runPart times one part and prints its line of the report. It returns the best
// time so the caller can build the combined total. A skipped part contributes zero.
func runPart(n int, part Part, input string, opts options) time.Duration {
	if opts.part != 0 && opts.part != n {
		return 0
	}

	if part == nil {
		printSkipped(n, "not written yet")
		return 0
	}

	res := timePart(part, input, opts.runs)

	if res.answer == nil {
		printSkipped(n, "returned nil")
		return 0
	}

	printResult(n, res)

	return res.best
}

// timePart calls part on the input `runs` times and collects the timings.
// The answer from the last run is kept.
func timePart(part Part, input string, runs int) result {
	res := result{runs: runs}

	var sum time.Duration
	for i := 0; i < runs; i++ {
		start := time.Now()
		res.answer = part(input)
		elapsed := time.Since(start)

		sum += elapsed
		if i == 0 || elapsed < res.best {
			res.best = elapsed
		}
	}

	res.avg = sum / time.Duration(runs)

	return res
}

// fail prints an error in the festive style and exits with a non-zero status.
func fail(msg string) {
	fmt.Fprintf(os.Stderr, "%s❌ %s%s\n", red, msg, reset)
	os.Exit(1)
}
