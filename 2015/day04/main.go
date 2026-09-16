// Advent of Code 2015, day 4: The Ideal Stocking Stuffer.
// https://adventofcode.com/2015/day/4
package main

import (
	"crypto/md5"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2015, 4, part1, part2)
}

// part1 finds the lowest number whose MD5 hash of secret+number starts with
// five hex zeros.
func part1(in string) any {
	return mine(strings.TrimSpace(in), fiveZeros)
}

// part2 does the same for six hex zeros.
func part2(in string) any {
	return mine(strings.TrimSpace(in), sixZeros)
}

// fiveZeros checks the raw 16-byte digest instead of formatting it as hex.
// Five hex zeros means the first two bytes are zero and the high nibble of
// the third byte is zero.
func fiveZeros(sum [md5.Size]byte) bool {
	return sum[0] == 0 && sum[1] == 0 && sum[2]>>4 == 0
}

// sixZeros means the first three bytes are all zero.
func sixZeros(sum [md5.Size]byte) bool {
	return sum[0] == 0 && sum[1] == 0 && sum[2] == 0
}

// chunkSize is how many numbers a worker claims at a time. Large enough that
// the atomic counter is not a bottleneck, small enough that the search stops
// soon after the answer is found.
const chunkSize = 10_000

// mine searches for the lowest number that makes the hash pass the check. The
// work is split across one goroutine per CPU. Each worker claims a chunk of
// numbers from a shared atomic counter, so no two workers hash the same value.
//
// Chunks are claimed in increasing order, but a worker on a later chunk can find
// a hit before a worker on an earlier chunk finishes. To guarantee the lowest
// answer, workers keep going until the chunk they would claim next starts
// beyond the best hit seen so far. Every number below the best hit has then
// been checked.
func mine(secret string, ok func([md5.Size]byte) bool) int {
	var (
		next atomic.Int64 // start of the next unclaimed chunk
		best atomic.Int64 // lowest hit so far; 0 means none yet
		wg   sync.WaitGroup
	)

	// The answer is never 0, so 0 works as "not found" and the search starts at 1.
	next.Store(1)

	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Each worker reuses one buffer: the secret followed by the number.
			// strconv.AppendInt writes the digits without allocating a string.
			buf := make([]byte, 0, len(secret)+20)
			buf = append(buf, secret...)

			for {
				start := next.Add(chunkSize) - chunkSize

				// Stop once every number below the best hit has been covered.
				if found := best.Load(); found != 0 && start > found {
					return
				}

				for n := start; n < start+chunkSize; n++ {
					sum := md5.Sum(strconv.AppendInt(buf[:len(secret)], n, 10))
					if !ok(sum) {
						continue
					}

					// Record n if it beats the current best. The loop retries when
					// another worker changed best in between, which is the standard
					// compare-and-swap pattern for a lock-free minimum.
					for {
						cur := best.Load()
						if cur != 0 && cur <= n {
							break
						}
						if best.CompareAndSwap(cur, n) {
							break
						}
					}

					// Nothing higher in this chunk can beat n.
					break
				}
			}
		}()
	}

	// Wait blocks until every worker has called Done.
	wg.Wait()

	return int(best.Load())
}
