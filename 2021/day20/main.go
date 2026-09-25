// Advent of Code 2021, day 20: Trench Map.
// https://adventofcode.com/2021/day/20
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2021, 20, part1, part2)
}

// part1 enhances the image two times and counts the light pixels.
func part1(in string) any {
	return enhanceTimes(in, 2)
}

// part2 enhances the image 50 times and counts the light pixels.
func part2(in string) any {
	return enhanceTimes(in, 50)
}

// image is the part of the infinite image that can hold light pixels, and
// the value of all the pixels outside it.
//
// The pixels outside all have the same value. In the real input, the first
// character of the algorithm is '#', so a dark background turns light, and
// the last character is '.', so it turns dark again on the next step.
type image struct {
	pixels     grid.Grid[bool]
	background bool
}

// enhanceTimes runs the enhancement the given number of times and counts the
// light pixels. The count is finite because the number of steps is even.
func enhanceTimes(in string, steps int) int {
	algorithm, img := parse(in)

	for range steps {
		img = enhance(algorithm, img)
	}

	return img.pixels.Count(func(lit bool) bool { return lit })
}

// parse reads the enhancement algorithm and the input image. The image starts
// on a dark background.
func parse(in string) ([]bool, image) {
	algoText, imageText, _ := strings.Cut(in, "\n\n")

	algorithm := make([]bool, len(algoText))
	for i, c := range algoText {
		algorithm[i] = c == '#'
	}

	pixels := grid.Map(imageText, func(c rune) bool { return c == '#' })

	return algorithm, image{pixels: pixels}
}

// enhance makes the next image.
//
// Only pixels at most one step outside the current image can see a pixel of
// that image. So the new image is one pixel larger on each side. Pixels
// further out see only background, so they all become the same new
// background value.
func enhance(algorithm []bool, img image) image {
	next := grid.New[bool](img.pixels.W+2, img.pixels.H+2)

	for p := range next.Points() {
		// The point in the old image is one up and one left.
		centre := p.Add(grid.P(-1, -1))
		next.Set(p, algorithm[index(img, centre)])
	}

	// A background pixel sees nine background pixels: index 0 or 511.
	bgIndex := 0
	if img.background {
		bgIndex = 511
	}

	return image{pixels: next, background: algorithm[bgIndex]}
}

// index reads the 3x3 square around p as a 9-bit binary number, from the top
// left to the bottom right. A light pixel is a 1.
func index(img image, p grid.Point) int {
	n := 0

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			lit, ok := img.pixels.Get(p.Add(grid.P(dx, dy)))
			if !ok {
				lit = img.background
			}

			n <<= 1
			if lit {
				n |= 1
			}
		}
	}

	return n
}
