#!/bin/sh
set -eu

cleanup() {
  kill "$WEB_PID" "$API_PID" 2>/dev/null || true
}

trap cleanup EXIT INT TERM

npm run start:web &
WEB_PID=$!

npm run start:api &
API_PID=$!

wait
