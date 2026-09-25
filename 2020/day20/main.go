// Advent of Code 2020, day 20: Jurassic Jigsaw.
// https://adventofcode.com/2020/day/20
package main

import (
	"strings"

	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2020, 20, part1, part2)
}

// part1 multiplies the IDs of the four corner tiles.
//
// We do not need to build the image for this. Each inner edge of the image is
// shared by exactly two tiles, and each outer edge belongs to one tile only.
// A corner tile is the only kind of tile with two unshared edges.
func part1(in string) any {
	tiles := parse(in)
	counts := edgeCounts(tiles)

	product := 1
	for _, t := range tiles {
		if unsharedEdges(t.rows, counts) == 2 {
			product *= t.id
		}
	}

	return product
}

// part2 builds the image, finds the sea monsters and counts the '#' cells
// that are not part of a monster.
//
// We try all 8 orientations of the image. Only one of them has monsters. In
// that one, we mark every cell that some monster covers, so cells shared by
// two monsters count only once.
func part2(in string) any {
	image := assemble(parse(in))

	for _, img := range orientations(image) {
		covered := map[[2]int]bool{}

		for y := range len(img) - monsterH + 1 {
			for x := range len(img[0]) - monsterW + 1 {
				if !monsterAt(img, x, y) {
					continue
				}

				for _, o := range monster {
					covered[[2]int{x + o[0], y + o[1]}] = true
				}
			}
		}

		if len(covered) == 0 {
			continue
		}

		total := 0
		for _, row := range img {
			total += strings.Count(row, "#")
		}

		return total - len(covered)
	}

	return nil
}

// monsterPattern is the sea monster from the puzzle text. Only its '#' cells
// matter. The spaces can be any value in the image.
var monsterPattern = []string{
	"                  # ",
	"#    ##    ##    ###",
	" #  #  #  #  #  #   ",
}

// monster, monsterW and monsterH hold the monster as (x, y) offsets of its
// '#' cells, plus its size. The init function fills them from monsterPattern.
var (
	monster            [][2]int
	monsterW, monsterH int
)

// init runs once before main. Go calls every init function in a package
// automatically, so this is a good place to build lookup tables.
func init() {
	monsterH = len(monsterPattern)
	monsterW = len(monsterPattern[0])

	for y, row := range monsterPattern {
		for x, ch := range row {
			if ch == '#' {
				monster = append(monster, [2]int{x, y})
			}
		}
	}
}

// monsterAt tells if a monster has its top-left corner at (x, y) in img.
func monsterAt(img []string, x, y int) bool {
	for _, o := range monster {
		if img[y+o[1]][x+o[0]] != '#' {
			return false
		}
	}

	return true
}

// tile is one piece of the image with its ID and its rows of '#' and '.'.
type tile struct {
	id   int
	rows []string
}

// assemble puts the tiles in a square, removes the border of each tile, and
// joins the insides into one image.
//
// In the puzzle, each edge pattern matches at most one other tile, so we can
// place the tiles greedily. We take a corner tile and turn it so that its
// two unshared edges face up and left. Then we fill the square row by row.
// For each position we look for the unused tile, in some orientation, whose
// left edge matches the right edge of the tile on its left, and whose top
// edge matches the bottom edge of the tile above it.
func assemble(tiles []tile) []string {
	counts := edgeCounts(tiles)
	side := 0
	for side*side < len(tiles) {
		side++
	}

	// Each tile has 8 orientations. Build them once for all tiles.
	variants := make([][][]string, len(tiles))
	for i, t := range tiles {
		variants[i] = orientations(t.rows)
	}

	grid := make([][]string, side*side)
	used := make([]bool, len(tiles))

	for i, t := range tiles {
		if unsharedEdges(t.rows, counts) != 2 {
			continue
		}

		for _, v := range variants[i] {
			if counts[canonical(top(v))] == 1 && counts[canonical(left(v))] == 1 {
				grid[0] = v
				used[i] = true
				break
			}
		}

		break
	}

	for pos := 1; pos < side*side; pos++ {
		r, c := pos/side, pos%side
		grid[pos] = findTile(variants, used, func(v []string) bool {
			if c > 0 && left(v) != right(grid[pos-1]) {
				return false
			}

			if r > 0 && top(v) != bottom(grid[pos-side]) {
				return false
			}

			return true
		})
	}

	return join(grid, side)
}

// findTile returns the first orientation of an unused tile that passes the
// fits check, and marks that tile as used.
func findTile(variants [][][]string, used []bool, fits func([]string) bool) []string {
	for i, vs := range variants {
		if used[i] {
			continue
		}

		for _, v := range vs {
			if fits(v) {
				used[i] = true
				return v
			}
		}
	}

	panic("no tile fits")
}

// join removes the one-cell border from each placed tile and joins the
// insides into the rows of the full image. grid holds the tiles in row-major
// order, side tiles per row.
func join(grid [][]string, side int) []string {
	inner := len(grid[0]) - 2
	var image []string

	for r := range side {
		for y := 1; y <= inner; y++ {
			var b strings.Builder
			for c := range side {
				row := grid[r*side+c][y]
				b.WriteString(row[1 : len(row)-1])
			}

			image = append(image, b.String())
		}
	}

	return image
}

// edgeCounts counts, for each edge pattern, how many tile edges have it. A
// tile can be flipped, so an edge and its reverse are the same edge. We
// store each edge under its canonical form to join the two.
func edgeCounts(tiles []tile) map[string]int {
	counts := map[string]int{}
	for _, t := range tiles {
		for _, e := range edges(t.rows) {
			counts[canonical(e)]++
		}
	}

	return counts
}

// unsharedEdges counts the edges of a tile that no other tile has.
func unsharedEdges(rows []string, counts map[string]int) int {
	n := 0
	for _, e := range edges(rows) {
		if counts[canonical(e)] == 1 {
			n++
		}
	}

	return n
}

// edges returns the four edges of a tile: top, bottom, left and right.
func edges(rows []string) []string {
	return []string{top(rows), bottom(rows), left(rows), right(rows)}
}

// top returns the first row, read from left to right.
func top(rows []string) string {
	return rows[0]
}

// bottom returns the last row, read from left to right.
func bottom(rows []string) string {
	return rows[len(rows)-1]
}

// left returns the first column, read from top to bottom.
func left(rows []string) string {
	return column(rows, 0)
}

// right returns the last column, read from top to bottom.
func right(rows []string) string {
	return column(rows, len(rows[0])-1)
}

// column returns column x, read from top to bottom.
func column(rows []string, x int) string {
	b := make([]byte, len(rows))
	for y, row := range rows {
		b[y] = row[x]
	}

	return string(b)
}

// canonical returns the smaller of an edge and its reverse. Two edges that
// can match after a flip then have the same canonical form.
func canonical(e string) string {
	b := []byte(e)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return min(e, string(b))
}

// orientations returns the 8 ways to place a square: 4 rotations of the
// square, and 4 rotations of its mirror image.
func orientations(rows []string) [][]string {
	var result [][]string

	for range 2 {
		for range 4 {
			result = append(result, rows)
			rows = rotate(rows)
		}

		rows = flip(rows)
	}

	return result
}

// rotate turns a square 90 degrees clockwise. The new row y is the old
// column y, read from bottom to top.
func rotate(rows []string) []string {
	n := len(rows)
	out := make([]string, n)

	for y := range n {
		b := make([]byte, n)
		for x := range n {
			b[x] = rows[n-1-x][y]
		}

		out[y] = string(b)
	}

	return out
}

// flip mirrors a square from left to right.
func flip(rows []string) []string {
	out := make([]string, len(rows))
	for y, row := range rows {
		b := []byte(row)
		for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
			b[i], b[j] = b[j], b[i]
		}

		out[y] = string(b)
	}

	return out
}

// parse reads the blank-line separated tiles. The first line of each block
// is "Tile N:", and the other lines are the rows.
func parse(in string) []tile {
	var tiles []tile
	for _, block := range input.Blocks(in) {
		tiles = append(tiles, tile{id: input.Ints(block[0])[0], rows: block[1:]})
	}

	return tiles
}
