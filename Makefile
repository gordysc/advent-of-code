# 🎄 Advent of Code in Go
#
# Usage:  make <target> YEAR=<year> DAY=<day>
# Examples:
#   make new YEAR=2023 DAY=1        scaffold a new day
#   make run YEAR=2023 DAY=1        run it against input.txt
#   make example YEAR=2023 DAY=1    run it against example.txt
#   make bench YEAR=2023 DAY=1      time it over several runs

.DEFAULT_GOAL := help
SHELL := /usr/bin/env bash

# 📅 Defaults: the current year, and the current day of the month.
# During December that is today's puzzle; override with YEAR= and DAY=.
YEAR ?= $(shell date +%Y)
DAY  ?= $(shell date +%-d)
RUNS ?= 10
ARGS ?=

# Zero-pad the day for the folder name: DAY=3 -> day03.
DD  := $(shell printf '%02d' $(DAY))
PKG := ./$(YEAR)/day$(DD)

# 🎨 A little colour.
GREEN  := \033[32m
RED    := \033[31m
YELLOW := \033[33m
BOLD   := \033[1m
DIM    := \033[2m
RESET  := \033[0m

.PHONY: help new run example part1 part2 bench year input test check fmt vet tidy clean today

## 🎅 Show this help
help:
	@printf '\n'
	@printf '   $(GREEN)      *      $(RESET)   $(BOLD)🎄 Advent of Code · Go 🎄$(RESET)\n'
	@printf '   $(GREEN)     ***     $(RESET)\n'
	@printf '   $(GREEN)    *****    $(RESET)   $(DIM)make <target> YEAR=$(YEAR) DAY=$(DAY)$(RESET)\n'
	@printf '   $(GREEN)   *******   $(RESET)\n'
	@printf '   $(GREEN)  *********  $(RESET)   $(BOLD)Solving$(RESET)\n'
	@printf '   $(RED)    |||||    $(RESET)   $(GREEN)new$(RESET)      🎁 scaffold YEAR/dayDD from the template\n'
	@printf '                   $(GREEN)input$(RESET)    ⬇️  download the puzzle input (needs .env)\n'
	@printf '                   $(GREEN)run$(RESET)      🛷 run both parts against input.txt\n'
	@printf '                   $(GREEN)example$(RESET)  📋 run both parts against example.txt\n'
	@printf '                   $(GREEN)part1$(RESET)    1️⃣  run only part 1\n'
	@printf '                   $(GREEN)part2$(RESET)    2️⃣  run only part 2\n'
	@printf '                   $(GREEN)bench$(RESET)    ⏱  time each part over RUNS=$(RUNS) runs\n'
	@printf '                   $(GREEN)year$(RESET)     🗓  run every day of YEAR\n'
	@printf '                   $(GREEN)today$(RESET)    ⛄ new + input for today (December only!)\n'
	@printf '\n'
	@printf '                   $(BOLD)Housekeeping$(RESET)\n'
	@printf '                   $(GREEN)test$(RESET)     🧪 test the helper library\n'
	@printf '                   $(GREEN)check$(RESET)    ✅ fmt + vet + test\n'
	@printf '                   $(GREEN)fmt$(RESET)      ✨ gofmt everything\n'
	@printf '                   $(GREEN)vet$(RESET)      🔍 go vet everything\n'
	@printf '                   $(GREEN)tidy$(RESET)     🧹 go mod tidy\n'
	@printf '                   $(GREEN)clean$(RESET)    🗑  remove build caches and binaries\n'
	@printf '\n'
	@printf '                   $(DIM)Extra flags: ARGS="-runs 5"  RUNS=20$(RESET)\n'
	@printf '\n'

## 🎁 Scaffold a new day
new:
	@scripts/new.sh $(YEAR) $(DAY)

## ⬇️  Download the input for a day
input:
	@scripts/input.sh $(YEAR) $(DAY)

## 🛷 Run both parts against input.txt
run: _exists
	@go run $(PKG) $(ARGS)

## 📋 Run both parts against example.txt
example: _exists
	@go run $(PKG) -example $(ARGS)

## 1️⃣  Run only part 1
part1: _exists
	@go run $(PKG) -part 1 $(ARGS)

## 2️⃣  Run only part 2
part2: _exists
	@go run $(PKG) -part 2 $(ARGS)

## ⏱  Time each part over RUNS runs
bench: _exists
	@go run $(PKG) -runs $(RUNS) $(ARGS)

## 🗓  Run every day of a year, in order
year:
	@if [[ ! -d "$(YEAR)" ]]; then \
		printf '$(RED)🎅 No solutions for $(YEAR) yet.$(RESET)\n'; exit 1; \
	fi
	@printf '\n$(BOLD)🗓  Advent of Code $(YEAR) · all days$(RESET)\n'
	@for dir in $(YEAR)/day*/; do \
		if [[ -f "$$dir/main.go" ]]; then go run "./$$dir" $(ARGS) || true; fi; \
	done

## ⛄ Scaffold and download today's puzzle
today:
	@if [[ "$$(date +%m)" != "12" ]]; then \
		printf '$(YELLOW)⛄ It is not December! Use: make new YEAR=<year> DAY=<day>$(RESET)\n'; exit 1; \
	fi
	@$(MAKE) --no-print-directory new YEAR=$(YEAR) DAY=$(DAY)
	@$(MAKE) --no-print-directory input YEAR=$(YEAR) DAY=$(DAY)

## 🧪 Test the helper library and the runner
test:
	@printf '🧪 Testing the elves'"'"' workshop...\n'
	@go test ./lib/... ./runner/...

## ✅ Format, vet and test everything
check: fmt vet test
	@printf '$(GREEN)✅ All good. Santa approves.$(RESET)\n'

## ✨ Format all Go code
fmt:
	@printf '✨ Polishing the code...\n'
	@gofmt -l -w .

## 🔍 Vet all Go code
vet:
	@printf '🔍 Checking the list twice...\n'
	@go vet ./...

## 🧹 Tidy the module file
tidy:
	@go mod tidy

## 🗑  Remove build caches and binaries
clean:
	@printf '🗑  Sweeping up the wrapping paper...\n'
	@go clean ./...
	@rm -rf bin/

# Internal: fail with a friendly message when the day folder is missing.
_exists:
	@if [[ ! -f "$(PKG)/main.go" ]]; then \
		printf '$(RED)🎅 No solution at $(PKG). Run: make new YEAR=$(YEAR) DAY=$(DAY)$(RESET)\n'; exit 1; \
	fi
