// Advent of Code 2016, day 5: How About a Nice Game of Chess?
// https://adventofcode.com/2016/day/5
package main

import (
	"crypto/md5"
	"iter"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 5, part1, part2)
}

// hexDigits turns a nibble into its lowercase hex character.
const hexDigits = "0123456789abcdef"

// part1 builds the password from the sixth hex character of each interesting
// hash, in the order the hashes are found.
func part1(in string) any {
	var password []byte

	for sum := range interesting(strings.TrimSpace(in)) {
		password = append(password, hexDigits[sum[2]&0xf])

		if len(password) == 8 {
			break
		}
	}

	return string(password)
}

// part2 uses the sixth hex character as a position and the seventh as the
// character to put there. Positions past 7, and positions already filled, are
// ignored.
func part2(in string) any {
	password := []byte("________")
	filled := 0

	for sum := range interesting(strings.TrimSpace(in)) {
		pos := sum[2] & 0xf
		if pos >= 8 || password[pos] != '_' {
			continue
		}

		password[pos] = hexDigits[sum[3]>>4]
		filled++

		if filled == 8 {
			break
		}
	}

	return string(password)
}

// chunkSize is how many indexes one worker hashes in each round.
const chunkSize = 50_000

// interesting yields, in index order, every MD5 digest of door+index whose hex
// form starts with five zeros.
//
// The work runs in rounds. In each round every CPU hashes its own chunk of
// indexes and keeps its hits. When all workers are done, the hits are yielded
// chunk by chunk, so the order is the same as a plain loop. The search stops
// after the round in which the caller stops asking.
func interesting(door string) iter.Seq[[md5.Size]byte] {
	return func(yield func([md5.Size]byte) bool) {
		workers := runtime.NumCPU()
		hits := make([][][md5.Size]byte, workers)

		for start := 0; ; start += workers * chunkSize {
			var wg sync.WaitGroup

			for w := range workers {
				wg.Add(1)

				go func() {
					defer wg.Done()

					hits[w] = scan(door, start+w*chunkSize, hits[w][:0])
				}()
			}

			wg.Wait()

			for _, chunk := range hits {
				for _, sum := range chunk {
					if !yield(sum) {
						return
					}
				}
			}
		}
	}
}

// scan hashes door+n for each n in one chunk and appends the digests that
// start with five hex zeros to out. Five hex zeros means the first two bytes
// are zero and the high nibble of the third byte is zero.
func scan(door string, from int, out [][md5.Size]byte) [][md5.Size]byte {
	// One buffer holds the door ID followed by the number. strconv.AppendInt
	// writes the digits without allocating a string.
	buf := make([]byte, 0, len(door)+20)
	buf = append(buf, door...)

	for n := from; n < from+chunkSize; n++ {
		sum := md5.Sum(strconv.AppendInt(buf[:len(door)], int64(n), 10))

		if sum[0] == 0 && sum[1] == 0 && sum[2]>>4 == 0 {
			out = append(out, sum)
		}
	}

	return out
}
