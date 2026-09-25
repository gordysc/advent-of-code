// Advent of Code 2022, day 22: Monkey Map.
// https://adventofcode.com/2022/day/22
package main

import (
	"math"

	"aoc/lib/grid"
	"aoc/lib/input"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2022, 22, part1, part2)
}

// part1 follows the path on the flat map. When a step goes off the map, it
// wraps to the far side of the same row or column.
func part1(in string) any {
	b, path := parse(in)

	return b.walk(path, b.wrapFlat)
}

// part2 follows the path with the map folded into a cube. When a step goes
// off a face, it continues on the face that is next to it on the cube.
//
// The layout of the net is different in the example and in the real
// inputs, so the fold is not hard-coded. See newCube for how it is found.
func part2(in string) any {
	b, path := parse(in)
	c := newCube(b)

	return b.walk(path, c.wrap)
}

// board is the map. Rows have different lengths, and a space or a missing
// character is off the map.
type board struct {
	rows []string
}

// at returns the tile at p: '.' for open, '#' for wall, or ' ' when p is
// off the map.
func (b board) at(p grid.Point) byte {
	if p.Y < 0 || p.Y >= len(b.rows) || p.X < 0 || p.X >= len(b.rows[p.Y]) {
		return ' '
	}

	return b.rows[p.Y][p.X]
}

// step is one instruction of the path: move forward some tiles, then turn.
// turn is 'L', 'R', or 0 for the last move, which has no turn after it.
type step struct {
	tiles int
	turn  byte
}

// parse reads the map and the path. input.Blocks keeps the leading spaces
// of each row, which are important: they give the column of each tile.
func parse(in string) (board, []step) {
	blocks := input.Blocks(in)
	b := board{rows: blocks[0]}

	var path []step
	line := blocks[1][0]
	n := 0

	for i := range len(line) {
		c := line[i]

		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
			continue
		}

		path = append(path, step{tiles: n, turn: c})
		n = 0
	}

	path = append(path, step{tiles: n})

	return b, path
}

// wrapFunc gives the position and direction after a step from p in
// direction dir goes off the map.
type wrapFunc func(p, dir grid.Point) (grid.Point, grid.Point)

// walk follows the path from the leftmost open tile of the top row, facing
// right, and returns the final password.
func (b board) walk(path []step, wrap wrapFunc) int {
	pos := grid.P(0, 0)

	for b.at(pos) != '.' {
		pos = pos.Add(grid.Right)
	}

	dir := grid.Right

	for _, s := range path {
		for range s.tiles {
			next, nextDir := pos.Add(dir), dir
			if b.at(next) == ' ' {
				next, nextDir = wrap(pos, dir)
			}

			// A wall stops the move, and the rest of the tiles are lost.
			if b.at(next) == '#' {
				break
			}

			pos, dir = next, nextDir
		}

		switch s.turn {
		case 'L':
			dir = dir.TurnLeft()
		case 'R':
			dir = dir.TurnRight()
		}
	}

	return 1000*(pos.Y+1) + 4*(pos.X+1) + facing(dir)
}

// facing gives the score of a direction: right 0, down 1, left 2, up 3.
func facing(dir grid.Point) int {
	switch dir {
	case grid.Right:
		return 0
	case grid.Down:
		return 1
	case grid.Left:
		return 2
	default:
		return 3
	}
}

// wrapFlat goes back from p against dir to the last tile on the map. That
// is the tile on the far side of the row or column.
func (b board) wrapFlat(p, dir grid.Point) (grid.Point, grid.Point) {
	back := dir.Reverse()

	for b.at(p.Add(back)) != ' ' {
		p = p.Add(back)
	}

	return p, dir
}

// vec is a point or a direction in 3D.
type vec struct {
	X, Y, Z int
}

// add returns v + w.
func (v vec) add(w vec) vec {
	return vec{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

// scale returns v multiplied by n.
func (v vec) scale(n int) vec {
	return vec{v.X * n, v.Y * n, v.Z * n}
}

// dot returns the dot product of v and w.
func (v vec) dot(w vec) int {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

// face is one side of the cube, and where it is on the flat map.
//
// normal points out of the cube. right and down are the 3D directions that
// point to the right and down on the flat map, when the face is folded
// into its place on the cube.
type face struct {
	corner              grid.Point // top-left tile of the face on the map
	normal, right, down vec
}

// cube is the map folded into a cube.
//
// Tile centres are in 3D coordinates with the cube centre at the origin.
// The units are half tiles, so that all tile centres have integer
// coordinates. The faces are at distance side from the origin.
type cube struct {
	side   int
	byNet  map[grid.Point]*face // key: face position on the map, in faces
	byNorm map[vec]*face
}

// newCube folds the map into a cube.
//
// The side of a face is the square root of the number of tiles divided by
// six. The first face stays flat, and a breadth-first search over the net
// gives each other face its 3D direction vectors. When the search goes from
// a face to the face on its right in the net, that face folds down over the
// right edge: its normal is the old right, and its right is the old normal
// reversed. The other three directions follow the same pattern.
func newCube(b board) cube {
	tiles := 0

	for _, row := range b.rows {
		for i := range len(row) {
			if row[i] != ' ' {
				tiles++
			}
		}
	}

	c := cube{
		side:   int(math.Round(math.Sqrt(float64(tiles / 6)))),
		byNet:  map[grid.Point]*face{},
		byNorm: map[vec]*face{},
	}

	// start is the position of the first face on the net, in faces.
	var start grid.Point

	for b.at(start.Scale(c.side)) == ' ' {
		start = start.Add(grid.Right)
	}

	first := &face{
		corner: start.Scale(c.side),
		normal: vec{0, 0, 1},
		right:  vec{1, 0, 0},
		down:   vec{0, 1, 0},
	}
	c.add(start, first)

	queue := []grid.Point{start}

	for len(queue) > 0 {
		at := queue[0]
		queue = queue[1:]
		f := c.byNet[at]

		for _, dir := range grid.Dirs4 {
			next := at.Add(dir)
			corner := next.Scale(c.side)

			if c.byNet[next] != nil || b.at(corner) == ' ' {
				continue
			}

			c.add(next, fold(f, dir, corner))
			queue = append(queue, next)
		}
	}

	return c
}

// add records the face f at position at of the net.
func (c cube) add(at grid.Point, f *face) {
	c.byNet[at] = f
	c.byNorm[f.normal] = f
}

// fold returns the face that is next to f in direction dir on the net,
// folded into its place on the cube.
func fold(f *face, dir grid.Point, corner grid.Point) *face {
	n, r, d := f.normal, f.right, f.down
	minus := func(v vec) vec { return v.scale(-1) }

	switch dir {
	case grid.Right:
		return &face{corner, r, minus(n), d}
	case grid.Left:
		return &face{corner, minus(r), n, d}
	case grid.Down:
		return &face{corner, d, r, minus(n)}
	default:
		return &face{corner, minus(d), r, n}
	}
}

// wrap gives the position and direction after a step from p in direction
// dir goes off the edge of its face.
//
// On the cube, the step goes over the edge onto the face whose normal is
// the 3D direction of the move. The new 3D direction is into the cube,
// the old normal reversed. In half-tile units, the new tile centre is one
// unit further along the move and one unit less along the old normal.
func (c cube) wrap(p, dir grid.Point) (grid.Point, grid.Point) {
	from := c.byNet[grid.P(p.X/c.side, p.Y/c.side)]
	pos := c.to3D(from, p)
	move := from.right.scale(dir.X).add(from.down.scale(dir.Y))

	to := c.byNorm[move]

	pos = pos.add(move).add(from.normal.scale(-1))
	newMove := from.normal.scale(-1)

	newDir := grid.P(newMove.dot(to.right), newMove.dot(to.down))

	return c.toMap(to, pos), newDir
}

// to3D gives the 3D centre of the tile p on face f.
func (c cube) to3D(f *face, p grid.Point) vec {
	i, j := p.X-f.corner.X, p.Y-f.corner.Y

	// The tile centres on a face go from 1-side to side-1 in steps of 2.
	return f.normal.scale(c.side).
		add(f.right.scale(2*i + 1 - c.side)).
		add(f.down.scale(2*j + 1 - c.side))
}

// toMap gives the map position of the tile with 3D centre v on face f. It
// is the reverse of to3D.
func (c cube) toMap(f *face, v vec) grid.Point {
	i := (v.dot(f.right) + c.side - 1) / 2
	j := (v.dot(f.down) + c.side - 1) / 2

	return f.corner.Add(grid.P(i, j))
}
