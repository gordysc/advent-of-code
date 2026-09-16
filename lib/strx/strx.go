// Package strx adds a few string helpers that the standard strings package
// lacks.
package strx

import "strings"

// Reverse returns the string with its runes in reverse order.
func Reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// IsDigit reports whether b is an ASCII digit.
func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// IsLetter reports whether b is an ASCII letter.
func IsLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// Chunk splits a string into pieces of n bytes. The last piece may be shorter.
func Chunk(s string, n int) []string {
	var out []string

	for i := 0; i < len(s); i += n {
		out = append(out, s[i:min(i+n, len(s))])
	}

	return out
}

// Between returns the text between the first occurrence of start and the next
// occurrence of end after it. It returns "" when either marker is missing.
func Between(s, start, end string) string {
	i := strings.Index(s, start)
	if i < 0 {
		return ""
	}
	s = s[i+len(start):]

	j := strings.Index(s, end)
	if j < 0 {
		return ""
	}

	return s[:j]
}

// Hamming returns how many positions differ between two equal-length strings.
func Hamming(a, b string) int {
	n := 0

	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			n++
		}
	}

	return n
}
