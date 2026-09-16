#!/usr/bin/env bash
# Download the puzzle input for one day using the session cookie in .env.
# Usage: scripts/input.sh YEAR DAY
set -euo pipefail

year="$1"
day="$((10#$2))"
root="$(cd "$(dirname "$0")/.." && pwd)"

padded="$(printf '%02d' "$day")"
dir="$root/$year/day$padded"
target="$dir/input.txt"

if [[ -f "$root/.env" ]]; then
  # shellcheck disable=SC1091
  source "$root/.env"
fi

if [[ -z "${AOC_SESSION:-}" ]]; then
  echo "🎅 No AOC_SESSION found. Copy .env.example to .env and add your session cookie."
  exit 1
fi

if [[ -s "$target" ]]; then
  echo "🎁 $year/day$padded/input.txt is already here. Delete it to download again."
  exit 0
fi

mkdir -p "$dir"

# Advent of Code asks that scripts send a User-Agent that identifies the user.
agent="${AOC_USER_AGENT:-advent-of-code go runner via curl}"

status="$(curl --silent --show-error \
  --cookie "session=$AOC_SESSION" \
  --user-agent "$agent" \
  --output "$target" \
  --write-out '%{http_code}' \
  "https://adventofcode.com/$year/day/$day/input")"

if [[ "$status" != "200" ]]; then
  echo "🎅 Download failed with HTTP $status. Is the session cookie still valid and the puzzle unlocked?"
  head -c 200 "$target" || true
  echo
  rm -f "$target"
  exit 1
fi

echo "⬇️  Saved $year/day$padded/input.txt ($(wc -l < "$target") lines)"
