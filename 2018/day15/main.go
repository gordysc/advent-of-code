// Advent of Code 2018, day 15: Beverage Bandits.
// https://adventofcode.com/2018/day/15
package main

import (
	"slices"

	"aoc/lib/grid"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 15, part1, part2)
}

// startHP and basePower are the hit points and attack power every unit starts
// with. Only the elves' power changes, and only in part 2.
const (
	startHP   = 200
	basePower = 3
)

// readingOrder lists the four neighbours of a square in reading order: up,
// left, right, down. Every tie in the puzzle breaks this way.
var readingOrder = []grid.Point{grid.Up, grid.Left, grid.Right, grid.Down}

// unit is one elf ('E') or goblin ('G').
type unit struct {
	kind  byte
	pos   grid.Point
	hp    int
	power int
}

// battle is the cave during one fight. walls holds only the walls and open
// floor. at holds the living unit on each square, or nil, so a unit can find
// its neighbours without a search.
type battle struct {
	walls grid.Grid[byte]
	at    grid.Grid[*unit]
	units []*unit
}

// part1 returns the outcome of the fight when elves hit as hard as goblins.
func part1(in string) any {
	outcome, _ := fight(in, basePower, false)

	return outcome
}

// part2 finds the smallest elf attack power that lets the elves win without
// a single loss, and returns the outcome of that fight. A fight stops as soon
// as an elf dies, because that power is too low.
func part2(in string) any {
	for power := basePower + 1; ; power++ {
		if outcome, ok := fight(in, power, true); ok {
			return outcome
		}
	}
}

// fight runs a whole battle with the given elf attack power. It returns the
// outcome: the number of full rounds times the hit points that are left.
// With stopOnElfDeath it stops at the end of the round in which the first
// elf dies and returns false. In all other cases it returns true.
func fight(in string, elfPower int, stopOnElfDeath bool) (int, bool) {
	b := parse(in, elfPower)
	elves := b.count('E')

	for rounds := 0; ; rounds++ {
		done := !b.round()

		// Check for a dead elf before the end of the fight, because the
		// last elf can die in the round that ends it.
		if stopOnElfDeath && b.count('E') < elves {
			return 0, false
		}

		if done {
			hp := 0
			for _, u := range b.units {
				hp += u.hp
			}

			return rounds * hp, true
		}
	}
}

// parse reads the cave and places a unit on each E and G. Elves get elfPower.
// The units move off the wall grid, so their squares become open floor.
func parse(in string, elfPower int) *battle {
	walls := grid.Parse(in)
	b := &battle{walls: walls, at: grid.New[*unit](walls.W, walls.H)}

	for p, c := range walls.All() {
		if c != 'E' && c != 'G' {
			continue
		}

		u := &unit{kind: c, pos: p, hp: startHP, power: basePower}
		if c == 'E' {
			u.power = elfPower
		}

		b.units = append(b.units, u)
		b.at.Set(p, u)
		walls.Set(p, '.')
	}

	return b
}

// count returns how many living units of one kind there are.
func (b *battle) count(kind byte) int {
	n := 0

	for _, u := range b.units {
		if u.kind == kind {
			n++
		}
	}

	return n
}

// round gives every unit one turn, in reading order of where the units stand
// when the round starts. It returns false when a unit finds no enemy left,
// which ends the fight before the round is complete.
func (b *battle) round() bool {
	slices.SortFunc(b.units, func(u, v *unit) int { return compare(u.pos, v.pos) })

	for _, u := range slices.Clone(b.units) {
		// A unit that died earlier in this round gets no turn.
		if u.hp <= 0 {
			continue
		}

		if b.count(enemy(u.kind)) == 0 {
			return false
		}

		if b.target(u) == nil {
			b.move(u)
		}

		if t := b.target(u); t != nil {
			b.attack(u, t)
		}
	}

	return true
}

// move takes one step toward the nearest square that is next to an enemy.
// First a search from the unit finds the distance to every open square it can
// reach. The chosen square is the closest one next to an enemy, the first in
// reading order on a tie. A second search, back from the chosen square, gives
// the step: the first neighbour in reading order that is one step closer.
func (b *battle) move(u *unit) {
	fromUnit := b.distances(u.pos)

	goal, best := grid.Point{}, -1
	for _, v := range b.units {
		if v.kind == u.kind {
			continue
		}

		for _, d := range readingOrder {
			p := v.pos.Add(d)
			dist := fromUnit.At(p)

			if dist < 0 {
				continue
			}

			if best < 0 || dist < best || dist == best && compare(p, goal) < 0 {
				goal, best = p, dist
			}
		}
	}

	if best < 0 {
		return
	}

	fromGoal := b.distances(goal)

	for _, d := range readingOrder {
		p := u.pos.Add(d)

		if fromGoal.At(p) == best-1 {
			b.at.Set(u.pos, nil)
			b.at.Set(p, u)
			u.pos = p

			return
		}
	}
}

// target returns the enemy next to u with the fewest hit points, the first
// in reading order on a tie, or nil when no enemy is next to u.
func (b *battle) target(u *unit) *unit {
	var best *unit

	for _, d := range readingOrder {
		v := b.at.At(u.pos.Add(d))

		if v != nil && v.kind != u.kind && (best == nil || v.hp < best.hp) {
			best = v
		}
	}

	return best
}

// attack deals u's damage to t. A unit with no hit points left dies and
// leaves the battle, which frees its square.
func (b *battle) attack(u, t *unit) {
	t.hp -= u.power
	if t.hp > 0 {
		return
	}

	b.at.Set(t.pos, nil)
	b.units = slices.DeleteFunc(b.units, func(v *unit) bool { return v == t })
}

// distances returns the number of steps from start to every open square that
// is free of units, or -1 where no path exists. It is a plain breadth-first
// search on a flat grid, which is much faster than a map-based search here
// because it runs twice per unit per round. The cave has a wall all around,
// so the search never leaves the grid.
func (b *battle) distances(start grid.Point) grid.Grid[int] {
	dist := grid.New[int](b.walls.W, b.walls.H)
	for i := range dist.Cells {
		dist.Cells[i] = -1
	}

	dist.Set(start, 0)
	queue := []grid.Point{start}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		for _, d := range readingOrder {
			q := p.Add(d)

			if b.walls.At(q) != '.' || b.at.At(q) != nil || dist.At(q) >= 0 {
				continue
			}

			dist.Set(q, dist.At(p)+1)
			queue = append(queue, q)
		}
	}

	return dist
}

// enemy returns the kind of unit that the given kind fights.
func enemy(kind byte) byte {
	if kind == 'E' {
		return 'G'
	}

	return 'E'
}

// compare orders two points in reading order: top to bottom, then left to
// right. It returns a negative number when a comes first.
func compare(a, b grid.Point) int {
	if a.Y != b.Y {
		return a.Y - b.Y
	}

	return a.X - b.X
}
