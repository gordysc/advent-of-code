// Advent of Code 2019, day 23: Category Six.
// https://adventofcode.com/2019/day/23
//
// The puzzle text has no example. example.txt holds a small handmade Intcode
// program: each machine reads its address, sends one packet to 255 with
// X = address and Y = address + 100, and then reads its input forever. Part
// 1 gives 100, the packet from machine 0, which runs first. In part 2 the
// last packet to 255 comes from machine 49, the network then goes idle, and
// the NAT sends that packet to machine 0 twice, so part 2 gives 149.
package main

import (
	"aoc/lib/intcode"
	"aoc/runner"
)

// main hands the two parts to the runner, which loads the input and times them.
func main() {
	runner.Run(2019, 23, part1, part2)
}

// machines is the number of computers on the network, with addresses 0 to 49.
// natAddr is the address outside the network that the NAT listens on.
const (
	machines = 50
	natAddr  = 255
)

// packet is one message on the network: it goes to address dest and carries
// the two values x and y.
type packet struct {
	dest, x, y int
}

// network is the set of computers and the state that the round-robin loop
// keeps between rounds.
type network struct {
	nics    []*intcode.Machine
	pending []bool  // true if a packet was sent to the NIC since it last ran
	partial [][]int // output of each NIC that does not make a full packet yet
}

// part1 returns the Y value of the first packet sent to address 255. Part 1
// has no NAT, so an idle network stays idle forever, and part1 returns nil
// if the network goes idle or halts before a packet reaches 255.
func part1(in string) any {
	net := newNetwork(in)

	for {
		toNAT, idle, alive := net.round()
		if len(toNAT) > 0 {
			return toNAT[0].y
		}

		if idle || !alive {
			return nil
		}
	}
}

// part2 adds the NAT. The NAT keeps only the last packet sent to 255. When
// the whole network is idle it sends that packet to address 0 to wake it up.
// part2 returns the first Y value the NAT sends to address 0 twice in a row.
func part2(in string) any {
	net := newNetwork(in)

	var last packet
	haveLast := false

	sentY, sentBefore := 0, false

	for {
		toNAT, idle, alive := net.round()
		if len(toNAT) > 0 {
			last = toNAT[len(toNAT)-1]
			haveLast = true
		}

		if !alive {
			return nil
		}

		if !idle {
			continue
		}

		// An idle network with nothing for the NAT to send stays idle
		// forever, so there is no answer.
		if !haveLast {
			return nil
		}

		if sentBefore && last.y == sentY {
			return last.y
		}

		net.deliver(packet{dest: 0, x: last.x, y: last.y})
		sentY, sentBefore = last.y, true
	}
}

// newNetwork boots 50 copies of the program and gives each one its address
// as its first input.
func newNetwork(in string) *network {
	program := intcode.Parse(in)

	net := &network{
		nics:    make([]*intcode.Machine, machines),
		pending: make([]bool, machines),
		partial: make([][]int, machines),
	}

	for addr := range machines {
		net.nics[addr] = intcode.New(program)
		net.nics[addr].Send(addr)
		net.pending[addr] = true
	}

	return net
}

// deliver puts a packet's X and Y on the input queue of the machine it is for.
func (net *network) deliver(p packet) {
	net.nics[p.dest].Send(p.x, p.y)
	net.pending[p.dest] = true
}

// round runs every machine once, in address order, without goroutines. Run
// stops when a machine needs input it does not have, so each machine reads
// all of its queued packets and then stops. A machine with no packets gets
// -1, which the puzzle says to send when its queue is empty.
//
// Packets for the network go straight onto the target's queue, so a machine
// later in the round can read them at once. round returns the packets sent to
// 255, in the order they were sent.
//
// idle is true when no machine had a packet waiting and no machine sent
// anything. Then every machine only read -1 and went back to waiting, and no
// packet is in flight, so the next round would do the same thing again.
// alive is false when every machine has halted.
func (net *network) round() (toNAT []packet, idle, alive bool) {
	idle = true

	for addr, nic := range net.nics {
		if nic.Halted() {
			continue
		}

		alive = true

		if net.pending[addr] {
			idle = false
			net.pending[addr] = false
		} else {
			nic.Send(-1)
		}

		out := nic.Run()
		if len(out) > 0 {
			idle = false
		}

		// A packet is three outputs. Keep any leftover values until the
		// machine sends the rest of its packet in a later round.
		buf := append(net.partial[addr], out...)
		for len(buf) >= 3 {
			p := packet{dest: buf[0], x: buf[1], y: buf[2]}
			buf = buf[3:]

			if p.dest == natAddr {
				toNAT = append(toNAT, p)
			} else {
				net.deliver(p)
			}
		}

		net.partial[addr] = buf
	}

	return toNAT, idle, alive
}
