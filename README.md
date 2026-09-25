# 🎄 Advent of Code in Go

Solutions for every year of [Advent of Code](https://adventofcode.com), written in Go and run from the CLI.

## Layout

```
.
├── Makefile          🎅 all the commands you need
├── runner/           ⏱  loads the input, times part 1 and part 2, prints the report
├── lib/              🧰 helper packages for puzzles
├── templates/        📄 the file that `make new` copies
├── scripts/          🔧 scaffold and download helpers
└── <year>/day<DD>/   🎁 one folder per puzzle
    ├── main.go
    ├── example.txt   (the sample from the puzzle text)
    └── input.txt     (your input, ignored by git)
```

Each day is a small `package main`. It calls `runner.Run` with its two part functions. The runner does the rest.

## Quick start

```sh
make new YEAR=2023 DAY=1     # scaffold 2023/day01/main.go
make input YEAR=2023 DAY=1   # download input.txt (see "Inputs" below)
make example YEAR=2023 DAY=1 # run against example.txt
make run YEAR=2023 DAY=1     # run against input.txt
make bench YEAR=2023 DAY=1   # best and average time over 10 runs
```

Run `make` with no target to see the full list.

`YEAR` and `DAY` default to the current date. In December, `make today` scaffolds and downloads the current puzzle.

## The report

```
🎄  Advent of Code 2023 · Day 01  (input.txt)
────────────────────────────────────────────────────
🎁  Part 1  54388                    ⏱  120.5µs
🎁  Part 2  53515                    ⏱  310.2µs
────────────────────────────────────────────────────
⭐  Total                            ⏱  430.7µs
```

The total is the sum of the two parts. It does not include the time to read the input file.

A part that returns `nil` shows as not solved. Multi-line answers, such as ASCII-art letters, print on their own lines.

## Flags

You can pass flags to a solution with `ARGS`, or run it directly with `go run`.

| Flag | Effect |
|------|--------|
| `-example` | Read `example.txt` instead of `input.txt`. |
| `-input PATH` | Read a specific file. |
| `-part 1` or `-part 2` | Run only one part. |
| `-runs N` | Run each part N times and show the best and average time. |

Example: `make run YEAR=2023 DAY=1 ARGS="-input other.txt"`.

## Inputs

Advent of Code asks that you do not publish puzzle inputs. The `.gitignore` file excludes every `input.txt`.

To download inputs, copy `.env.example` to `.env` and add your session cookie. The `.env` file is also ignored by git.

## Helper library

Import what you need from `aoc/lib/...`. Every function has a doc comment.

| Package | Contents |
|---------|----------|
| `input` | `Lines`, `Blocks`, `Ints`, `UInts`, `Int`, `IntLines`, `Digits`, `Fields` |
| `grid` | `Point` with directions and turns, generic `Grid[T]` with `Parse`, `Map`, iterators and neighbour lookups |
| `mathx` | `Abs`, `Sign`, `GCD`, `LCM`, `Mod`, `Pow`, `DivCeil`, `Triangular`, `Digits` |
| `set` | Generic `Set[T]` with union, intersection and difference |
| `ds` | `Stack`, `Queue` and `PriorityQueue` |
| `search` | `BFS`, `Flood`, `Dijkstra`, `DijkstraAll`, `BinarySearch` |
| `slicesx` | `Sum`, `Product`, `Map`, `Filter`, `Count`, `Counts`, `MinMax`, `Chunk`, `Windows`, `Pairs`, `Transpose`, `Permutations`, `Combinations` |
| `strx` | `Reverse`, `IsDigit`, `IsLetter`, `Chunk`, `Between`, `Hamming` |
| `intcode` | The 2019 Intcode computer: `Parse`, `New`, and a `Machine` that runs until it needs input or halts |

Run `make test` to test the library. Run `make check` to format, vet and test everything.
