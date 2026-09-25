// Advent of Code 2016, day 14: One-Time Pad.
// https://adventofcode.com/2016/day/14
package main

import (
	"crypto/md5"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2016, 14, part1, part2)
}

// The rules for a key from the puzzle.
const (
	keysNeeded = 64   // the part asks for the index of the 64th key
	window     = 1000 // how many later hashes can confirm a key
	stretches  = 2016 // extra hashes of each hash in part 2
)

// hexDigits turns a nibble into its lowercase hex character.
const hexDigits = "0123456789abcdef"

// part1 finds the index that makes the 64th key with plain MD5 hashes.
func part1(in string) any {
	return findKey(strings.TrimSpace(in), 0)
}

// part2 finds the index that makes the 64th key with stretched hashes.
func part2(in string) any {
	return findKey(strings.TrimSpace(in), stretches)
}

// runs keeps the only facts about one hash that the key test needs.
type runs struct {
	triple int    // the digit of the first run of three, or -1 if there is none
	fives  uint16 // bit d is set when the hash has a run of five of digit d
}

// findKey walks the indexes in order and returns the one that makes the 64th
// key. An index makes a key when its hash has a run of three of some digit,
// and one of the next 1000 hashes has a run of five of the same digit.
//
// Each hash is needed for up to 1001 indexes, so every hash is computed once
// and kept in hashes.
func findKey(salt string, extra int) int {
	var hashes []runs
	found := 0

	for i := 0; ; i++ {
		for len(hashes) <= i+window {
			hashes = more(salt, extra, hashes)
		}

		digit := hashes[i].triple
		if digit < 0 {
			continue
		}

		for j := i + 1; j <= i+window; j++ {
			if hashes[j].fives&(1<<digit) == 0 {
				continue
			}

			found++
			if found == keysNeeded {
				return i
			}

			break
		}
	}
}

// chunkSize is how many indexes one worker hashes in each round.
const chunkSize = 1000

// more hashes the next round of indexes and appends the results to hashes.
// Every CPU hashes its own chunk at the same time. Each worker writes to its
// own part of the slice, so the workers do not need a lock, and the results
// stay in index order.
func more(salt string, extra int, hashes []runs) []runs {
	workers := runtime.NumCPU()
	from := len(hashes)
	hashes = append(hashes, make([]runs, workers*chunkSize)...)

	var wg sync.WaitGroup

	for w := range workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			start := from + w*chunkSize
			scan(salt, extra, start, hashes[start:start+chunkSize])
		}()
	}

	wg.Wait()

	return hashes
}

// scan hashes salt+n for each n from start on, one for each slot in out. Each
// hash is hashed again extra more times, and then its runs are stored.
func scan(salt string, extra, start int, out []runs) {
	// One buffer holds the salt followed by the number. strconv.AppendInt
	// writes the digits without allocating a string.
	buf := make([]byte, 0, len(salt)+20)
	buf = append(buf, salt...)

	// hex is reused for every hash, so the hot loop does not allocate.
	var hex [2 * md5.Size]byte

	for k := range out {
		sum := md5.Sum(strconv.AppendInt(buf[:len(salt)], int64(start+k), 10))
		toHex(&hex, sum)

		for range extra {
			toHex(&hex, md5.Sum(hex[:]))
		}

		out[k] = findRuns(&hex)
	}
}

// toHex writes the lowercase hex form of a digest into hex.
func toHex(hex *[2 * md5.Size]byte, sum [md5.Size]byte) {
	for i, b := range sum {
		hex[2*i] = hexDigits[b>>4]
		hex[2*i+1] = hexDigits[b&0xf]
	}
}

// findRuns finds the first run of three equal characters, and every digit
// that has a run of five.
func findRuns(hex *[2 * md5.Size]byte) runs {
	r := runs{triple: -1}

	for i := 0; i+2 < len(hex); i++ {
		c := hex[i]
		if hex[i+1] != c || hex[i+2] != c {
			continue
		}

		digit := strings.IndexByte(hexDigits, c)
		if r.triple < 0 {
			r.triple = digit
		}

		if i+4 < len(hex) && hex[i+3] == c && hex[i+4] == c {
			r.fives |= 1 << digit
		}
	}

	return r
}
