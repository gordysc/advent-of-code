package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ANSI escape codes for colour. They are set to "" when colour is disabled, so the
// rest of the code can use them without checks.
var (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

// init turns colour off when the output is not a terminal (for example when piped
// to a file) or when the NO_COLOR convention (https://no-color.org) is set.
func init() {
	info, err := os.Stdout.Stat()
	isTerminal := err == nil && info.Mode()&os.ModeCharDevice != 0

	if _, noColor := os.LookupEnv("NO_COLOR"); noColor || !isTerminal {
		reset, bold, dim, red, green, yellow = "", "", "", "", "", ""
	}
}

// ruleWidth is the width of the horizontal lines in the report.
const ruleWidth = 52

// rule returns a horizontal line for the report.
func rule() string {
	return dim + strings.Repeat("─", ruleWidth) + reset
}

// printHeader prints the report title line.
func printHeader(year, day int, inputPath string, runs int) {
	extra := ""
	if runs > 1 {
		extra = fmt.Sprintf("  ·  %d runs, best (avg)", runs)
	}

	fmt.Printf("\n🎄  %sAdvent of Code %d · Day %02d%s  %s(%s%s)%s\n",
		bold, year, day, reset, dim, filepath.Base(inputPath), extra, reset)
	fmt.Println(rule())
}

// printResult prints one solved part with its answer and timing.
func printResult(n int, res result) {
	answer := fmt.Sprintf("%v", res.answer)

	// Multi-line answers (for example ASCII-art letters) go on their own lines.
	if strings.Contains(answer, "\n") {
		fmt.Printf("🎁  %sPart %d%s  %s\n%s\n", bold, n, reset, timing(res), answer)
		return
	}

	fmt.Printf("🎁  %sPart %d%s  %s%-24s%s %s\n", bold, n, reset, green, answer, reset, timing(res))
}

// printSkipped prints a part that produced no answer.
func printSkipped(n int, why string) {
	fmt.Printf("🎁  %sPart %d%s  %s— %s%s\n", bold, n, reset, yellow, why, reset)
}

// printFooter prints the combined time for both parts.
func printFooter(total time.Duration) {
	fmt.Println(rule())
	fmt.Printf("⭐  %sTotal%s   %-24s ⏱  %s%s%s\n\n", bold, reset, "", bold, format(total), reset)
}

// timing formats the duration column for one part.
func timing(res result) string {
	if res.runs > 1 {
		return fmt.Sprintf("⏱  %s%s%s %s(%s)%s", bold, format(res.best), reset, dim, format(res.avg), reset)
	}

	return fmt.Sprintf("⏱  %s%s%s", bold, format(res.best), reset)
}

// format renders a duration with a sensible amount of precision: whole
// nanoseconds below a microsecond, and three significant decimals above.
func format(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return d.String()
	case d < time.Millisecond:
		return d.Round(time.Nanosecond * 10).String()
	case d < time.Second:
		return d.Round(time.Microsecond * 10).String()
	default:
		return d.Round(time.Millisecond).String()
	}
}
