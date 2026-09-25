// Advent of Code 2018, day 24: Immune System Simulator 20XX.
// https://adventofcode.com/2018/day/24
package main

import (
	"cmp"
	"slices"
	"strings"

	"aoc/lib/input"
	"aoc/lib/set"
	"aoc/lib/strx"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2018, 24, part1, part2)
}

// The two armies.
const (
	immune = iota
	infection
)

// stalemate is the winner that fight reports when a round kills no units.
// The same round then repeats forever, so no army can win.
const stalemate = -1

// group is a group of identical units in one army.
type group struct {
	army       int
	units      int
	hp         int
	damage     int
	attack     string
	initiative int
	weak       set.Set[string]
	immune     set.Set[string]
}

// power returns the group's effective power: its units times their damage.
func (g *group) power() int {
	return g.units * g.damage
}

// damageTo returns the damage the group would deal to the defender. Immunity
// makes it zero, and weakness doubles it.
func (g *group) damageTo(d *group) int {
	switch {
	case d.immune.Has(g.attack):
		return 0
	case d.weak.Has(g.attack):
		return 2 * g.power()
	default:
		return g.power()
	}
}

// part1 lets the armies fight and counts the units the winner has left.
func part1(in string) any {
	_, units := fight(parse(in), 0)

	return units
}

// part2 finds the smallest boost to the immune system's damage that lets it
// win, and counts the units it has left. A bigger boost does not always help,
// because some boosts end in a stalemate, so the boosts are tried in order.
func part2(in string) any {
	groups := parse(in)

	for boost := 1; ; boost++ {
		winner, units := fight(groups, boost)
		if winner == immune {
			return units
		}
	}
}

// parse reads the two armies. Each block starts with the army's name, and
// every line after it is one group.
func parse(in string) []group {
	var groups []group

	for _, block := range input.Blocks(in) {
		army := immune
		if strings.HasPrefix(block[0], "Infection") {
			army = infection
		}

		for _, line := range block[1:] {
			groups = append(groups, parseGroup(army, line))
		}
	}

	return groups
}

// parseGroup reads one line such as:
//
//	17 units each with 5390 hit points (weak to radiation, bludgeoning) with
//	an attack that does 4507 fire damage at initiative 2
//
// The part in brackets can list weaknesses and immunities in either order, or
// be missing.
func parseGroup(army int, line string) group {
	n := input.Ints(line)
	g := group{
		army:       army,
		units:      n[0],
		hp:         n[1],
		damage:     n[2],
		initiative: n[3],
		weak:       set.New[string](),
		immune:     set.New[string](),
	}

	// The attack type is the word after the damage number.
	g.attack = strings.Fields(strx.Between(line, "does ", " damage"))[1]

	// The bracket part looks like "immune to fire; weak to bludgeoning,
	// slashing". Each piece before " to " says which set the list goes in.
	for _, part := range strings.Split(strx.Between(line, "(", ")"), "; ") {
		kind, list, ok := strings.Cut(part, " to ")
		if !ok {
			continue
		}

		target := g.weak
		if kind == "immune" {
			target = g.immune
		}

		target.Add(strings.Split(list, ", ")...)
	}

	return g
}

// fight runs the battle with the immune system's damage raised by boost. It
// returns the winning army and the units it has left, or stalemate.
func fight(start []group, boost int) (int, int) {
	// Each fight works on copies, so part 2 can reuse the parsed groups. The
	// weak and immune sets are maps, which the copies share, but no one
	// changes them.
	groups := make([]*group, len(start))
	for i := range start {
		g := start[i]
		if g.army == immune {
			g.damage += boost
		}

		groups[i] = &g
	}

	for {
		alive := [2]int{}
		for _, g := range groups {
			alive[g.army] += g.units
		}

		if alive[immune] == 0 {
			return infection, alive[infection]
		}

		if alive[infection] == 0 {
			return immune, alive[immune]
		}

		if !round(groups) {
			return stalemate, 0
		}

		// Drop the groups that died, so later rounds skip them.
		groups = slices.DeleteFunc(groups, func(g *group) bool { return g.units <= 0 })
	}
}

// round plays one round: every group picks a target, and then every group
// attacks. It reports whether any unit died.
func round(groups []*group) bool {
	// Groups choose in order of effective power, then initiative.
	slices.SortFunc(groups, func(a, b *group) int {
		return cmp.Or(cmp.Compare(b.power(), a.power()), cmp.Compare(b.initiative, a.initiative))
	})

	targets := map[*group]*group{}
	taken := map[*group]bool{}

	for _, g := range groups {
		var best *group
		bestDamage := 0

		// The best target takes the most damage. Ties go to the target with
		// more effective power, then to the one with higher initiative.
		for _, d := range groups {
			if d.army == g.army || taken[d] {
				continue
			}

			dmg := g.damageTo(d)
			if dmg == 0 {
				continue
			}

			if best == nil || dmg > bestDamage ||
				dmg == bestDamage && (d.power() > best.power() ||
					d.power() == best.power() && d.initiative > best.initiative) {
				best, bestDamage = d, dmg
			}
		}

		if best != nil {
			targets[g] = best
			taken[best] = true
		}
	}

	// Groups attack in order of initiative. A group that died earlier in the
	// round does not attack, and a group that lost units deals less damage.
	slices.SortFunc(groups, func(a, b *group) int { return cmp.Compare(b.initiative, a.initiative) })

	killed := false
	for _, g := range groups {
		d, ok := targets[g]
		if !ok || g.units <= 0 {
			continue
		}

		dead := min(g.damageTo(d)/d.hp, d.units)
		d.units -= dead
		if dead > 0 {
			killed = true
		}
	}

	return killed
}
