// Advent of Code 2018, day 23: Experimental Emergency Teleportation.
// https://adventofcode.com/2018/day/23
//
// The two parts use different examples in the puzzle text. example.txt holds
// the part 1 example, and part 2 gives 1 on it: the point (1, 0, 0) is in
// range of the most bots. The part 2 example gives 36.
package main

import (
	"container/heap"

	"aoc/lib/input"
	"aoc/lib/mathx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 23, part1, part2)
}

// bot is a nanobot: its position and its signal radius.
type bot struct {
	pos [3]int
	r   int
}

// box is a cube of whole-number points: it starts at lo and is size points
// long on each axis. count is how many bots reach at least one of its points,
// and dist is the distance from the origin to its nearest point.
type box struct {
	lo    [3]int
	size  int
	count int
	dist  int
}

// boxQueue is a priority queue of boxes. The box that pops first is the one
// that the most bots reach, then the one nearest the origin, then the smallest.
//
// The order has three keys, so it does not fit the single int priority of
// ds.PriorityQueue. Instead boxQueue has the five methods of heap.Interface,
// and the container/heap functions keep it in heap order. Push and Pop change
// the slice's length, so they need a pointer receiver.
type boxQueue []box

// Len returns the number of boxes in the queue.
func (q boxQueue) Len() int { return len(q) }

// Swap exchanges two boxes.
func (q boxQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

// Push adds a box at the end. heap.Push calls it and then moves the box up
// to its place.
func (q *boxQueue) Push(x any) { *q = append(*q, x.(box)) }

// Less reports whether box i must pop before box j.
func (q boxQueue) Less(i, j int) bool {
	a, b := q[i], q[j]
	if a.count != b.count {
		return a.count > b.count
	}

	if a.dist != b.dist {
		return a.dist < b.dist
	}

	return a.size < b.size
}

// Pop removes the last box. heap.Pop first swaps the best box to the end, so
// this returns the best box.
func (q *boxQueue) Pop() any {
	old := *q
	b := old[len(old)-1]
	*q = old[:len(old)-1]

	return b
}

// part1 counts the bots in range of the bot with the largest radius,
// including that bot.
func part1(in string) any {
	bots := parse(in)

	strongest := bots[0]
	for _, b := range bots {
		if b.r > strongest.r {
			strongest = b
		}
	}

	count := 0
	for _, b := range bots {
		if distance(strongest.pos, b.pos) <= strongest.r {
			count++
		}
	}

	return count
}

// part2 finds the point in range of the most bots, and returns its distance
// from the origin. When many points tie, it takes the one nearest the origin.
//
// The search starts with one cube around every point that any bot reaches,
// and splits cubes into eight smaller ones. The best point can be outside the
// box around the bots' positions, so the first cube includes each bot's range. A cube's count is the number of bots that reach any of
// its points, so no point inside can beat it. A cube's distance is that of
// its nearest point, so no point inside can be nearer. The queue always takes
// the best cube next, so the first cube of a single point to come out has the
// most bots, and it is the nearest of the points that tie.
func part2(in string) any {
	bots := parse(in)

	lo := bots[0].pos
	hi := bots[0].pos
	for _, b := range bots {
		for i := range 3 {
			lo[i] = min(lo[i], b.pos[i]-b.r)
			hi[i] = max(hi[i], b.pos[i]+b.r)
		}
	}

	// The first cube starts at the low corner of every bot's range. Its size is a
	// power of two, so it splits evenly all the way down to single points.
	size := 1
	for i := range 3 {
		for size <= hi[i]-lo[i] {
			size *= 2
		}
	}

	q := &boxQueue{measure(bots, lo, size)}

	for q.Len() > 0 {
		cur := heap.Pop(q).(box)
		if cur.size == 1 {
			return cur.dist
		}

		// Each bit of corner picks the low or the high half on one axis.
		half := cur.size / 2
		for corner := range 8 {
			child := cur.lo
			for i := range 3 {
				if corner&(1<<i) != 0 {
					child[i] += half
				}
			}

			heap.Push(q, measure(bots, child, half))
		}
	}

	return nil
}

// parse reads one bot from every line.
func parse(in string) []bot {
	var bots []bot
	for _, line := range input.Lines(in) {
		n := input.Ints(line)
		bots = append(bots, bot{pos: [3]int{n[0], n[1], n[2]}, r: n[3]})
	}

	return bots
}

// measure builds the box at lo with the given size. It counts the bots that
// reach it and finds its distance from the origin.
func measure(bots []bot, lo [3]int, size int) box {
	b := box{lo: lo, size: size}

	for _, bt := range bots {
		if gap(bt.pos, lo, size) <= bt.r {
			b.count++
		}
	}

	b.dist = gap([3]int{}, lo, size)

	return b
}

// gap returns the Manhattan distance from p to the nearest point of the cube
// at lo with the given size. On each axis the distance is zero when p is
// level with the cube, and otherwise it is the distance to the nearer face.
func gap(p, lo [3]int, size int) int {
	d := 0
	for i := range 3 {
		hi := lo[i] + size - 1
		d += max(0, lo[i]-p[i], p[i]-hi)
	}

	return d
}

// distance returns the Manhattan distance between two points.
func distance(a, b [3]int) int {
	d := 0
	for i := range 3 {
		d += mathx.Abs(a[i] - b[i])
	}

	return d
}
