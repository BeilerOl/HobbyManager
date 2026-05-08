#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND="$ROOT/backend"
FRONTEND="$ROOT/frontend"

cleanup() {
  echo ""
  echo "Stopping..."
  # Kill entire process groups so child processes (e.g. vite spawned by npm) are also terminated
  kill -- -"$BACKEND_PID" -"$FRONTEND_PID" 2>/dev/null || true
  wait "$BACKEND_PID" "$FRONTEND_PID" 2>/dev/null || true
}
trap cleanup INT TERM

echo "==> Starting backend (http://localhost:8080)..."
cd "$BACKEND"
setsid go run ./cmd/server/main.go &
BACKEND_PID=$!

echo "==> Installing frontend dependencies..."
cd "$FRONTEND"
npm install --silent

echo "==> Starting frontend (http://localhost:5173)..."
setsid npm run dev &
FRONTEND_PID=$!

echo ""
echo "Backend  : http://localhost:8080"
echo "Frontend : http://localhost:5173"
echo "Press Ctrl+C to stop both."
echo ""

wait "$BACKEND_PID" "$FRONTEND_PID"
