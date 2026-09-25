// Advent of Code 2019, day 8: Space Image Format.
// https://adventofcode.com/2019/day/8
package main

import (
	"strings"

	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 8, part1, part2)
}

// width and height are the size of the image in the puzzle. The worked
// examples in the puzzle text use smaller images: 3x2 for part 1 and 2x2 for
// part 2. example.txt is the part 1 example. Its 12 digits do not fill even
// one 25x6 layer, so neither part has an answer for it with these values.
const (
	width  = 25
	height = 6
)

// Pixel colors in the image data.
const (
	black       = '0'
	white       = '1'
	transparent = '2'
)

// on and off are the characters that part 2 draws for white and black pixels.
const (
	on  = '#'
	off = ' '
)

// part1 finds the layer with the fewest 0 digits, and multiplies its number
// of 1 digits by its number of 2 digits.
func part1(in string) any {
	layers := split(in)
	if layers == nil {
		return nil
	}

	best := layers[0]
	for _, layer := range layers[1:] {
		if strings.Count(layer, "0") < strings.Count(best, "0") {
			best = layer
		}
	}

	return strings.Count(best, "1") * strings.Count(best, "2")
}

// part2 stacks the layers and draws the image. The first layer is on top, so
// each pixel takes its color from the first layer that is not transparent
// there. The white pixels spell out the message as capital letters.
// Grid.String ends every row with a newline, so the last one is trimmed.
func part2(in string) any {
	layers := split(in)
	if layers == nil {
		return nil
	}

	image := grid.New[byte](width, height)

	for i := range image.Cells {
		image.Cells[i] = off

		for _, layer := range layers {
			if layer[i] == transparent {
				continue
			}

			if layer[i] == white {
				image.Cells[i] = on
			}

			break
		}
	}

	return strings.TrimRight(image.String(), "\n")
}

// split cuts the digits into layers of width*height pixels. Slicing a string
// does not copy it, so each layer shares memory with the input. It returns
// nil when the digits do not fill a whole number of layers.
func split(in string) []string {
	size := width * height
	if len(in) == 0 || len(in)%size != 0 {
		return nil
	}

	var layers []string
	for i := 0; i < len(in); i += size {
		layers = append(layers, in[i:i+size])
	}

	return layers
}
