#!/usr/bin/env sh
# Backend coverage gate: total statement coverage over the testable application
# packages must be strictly greater than 80%.
#
# cmd/* packages are thin composition roots with no unit tests of their own; the
# behaviour they wire together is covered through internal/app. They are
# therefore excluded from the gate and reported separately for transparency.
set -eu

profile="$(mktemp -t maulid-cover.XXXXXX)"
trap 'rm -f "$profile"' EXIT

go test -coverprofile="$profile" ./internal/... >/dev/null

total="$(go tool cover -func="$profile" | awk '/^total:/ {gsub(/%/, "", $3); print $3}')"

echo "backend statements: ${total}% (gate: > 80%)"

awk -v total="$total" 'BEGIN { if (total + 0 <= 80) { print "FAIL: coverage must be strictly greater than 80%"; exit 1 } }'
