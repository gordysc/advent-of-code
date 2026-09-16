// Package grid provides 2D points, directions and a generic grid type.
//
// The coordinate system matches how puzzle inputs are laid out as text:
// X grows to the right and Y grows downward, so (0, 0) is the top-left
// character and Up is (0, -1).
package grid

// Point is a position or a direction on the grid.
type Point struct {
	X, Y int
}

// Compass directions. Y grows downward, so Up has a negative Y.
var (
	Up    = Point{0, -1}
	Down  = Point{0, 1}
	Left  = Point{-1, 0}
	Right = Point{1, 0}

	UpLeft    = Point{-1, -1}
	UpRight   = Point{1, -1}
	DownLeft  = Point{-1, 1}
	DownRight = Point{1, 1}
)

// Dirs4 lists the four orthogonal directions in clockwise order starting at Up.
var Dirs4 = []Point{Up, Right, Down, Left}

// Dirs8 lists all eight directions in clockwise order starting at Up.
var Dirs8 = []Point{Up, UpRight, Right, DownRight, Down, DownLeft, Left, UpLeft}

// DirFromRune maps the characters puzzles commonly use for directions.
var DirFromRune = map[rune]Point{
	'^': Up, 'v': Down, '<': Left, '>': Right,
	'U': Up, 'D': Down, 'L': Left, 'R': Right,
	'N': Up, 'S': Down, 'W': Left, 'E': Right,
}

// P builds a Point. It is shorter than Point{x, y} in dense solution code.
func P(x, y int) Point {
	return Point{x, y}
}

// Add returns p + q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns p - q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Scale returns p multiplied by n. Useful to step n times in a direction.
func (p Point) Scale(n int) Point {
	return Point{p.X * n, p.Y * n}
}

// Manhattan returns the taxicab distance |dx| + |dy| between p and q.
func (p Point) Manhattan(q Point) int {
	return abs(p.X-q.X) + abs(p.Y-q.Y)
}

// Chebyshev returns max(|dx|, |dy|), the number of king moves between p and q.
func (p Point) Chebyshev(q Point) int {
	return max(abs(p.X-q.X), abs(p.Y-q.Y))
}

// Neighbors4 returns the four orthogonal neighbours of p.
func (p Point) Neighbors4() []Point {
	return p.neighbors(Dirs4)
}

// Neighbors8 returns the eight surrounding points of p.
func (p Point) Neighbors8() []Point {
	return p.neighbors(Dirs8)
}

// neighbors returns p offset by each direction.
func (p Point) neighbors(dirs []Point) []Point {
	out := make([]Point, len(dirs))

	for i, d := range dirs {
		out[i] = p.Add(d)
	}

	return out
}

// TurnRight rotates a direction 90° clockwise (Up becomes Right).
func (p Point) TurnRight() Point {
	return Point{-p.Y, p.X}
}

// TurnLeft rotates a direction 90° counter-clockwise (Up becomes Left).
func (p Point) TurnLeft() Point {
	return Point{p.Y, -p.X}
}

// Reverse returns the opposite direction.
func (p Point) Reverse() Point {
	return Point{-p.X, -p.Y}
}

// abs returns the absolute value. It lives here so the package has no imports.
func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
