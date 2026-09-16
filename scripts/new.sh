#!/usr/bin/env bash
# Scaffold a new puzzle day from the template.
# Usage: scripts/new.sh YEAR DAY
set -euo pipefail

year="$1"
day="$2"
root="$(cd "$(dirname "$0")/.." && pwd)"

padded="$(printf '%02d' "$day")"
dir="$root/$year/day$padded"

if [[ -e "$dir/main.go" ]]; then
  echo "🎅 Ho ho no! $year/day$padded already exists. Nothing was changed."
  exit 1
fi

mkdir -p "$dir"

sed -e "s/{{YEAR}}/$year/g" -e "s/{{DAY}}/$((10#$day))/g" \
  "$root/templates/day.go.tmpl" > "$dir/main.go"

touch "$dir/example.txt"

echo "🎁 Unwrapped $year/day$padded/"
echo "   ├── main.go       ✏️  write your solution here"
echo "   ├── example.txt   📋 paste the example from the puzzle text"
echo "   └── input.txt     ⬇️  run: make input YEAR=$year DAY=$day"
