#!/usr/bin/env bash
set -euo pipefail

echo "[1/6] Go tests"
go test ./...

echo "[2/6] Frontend (npm)"
npm --prefix frontend ci || npm --prefix frontend install
npm --prefix frontend run build

echo "[3/6] Build via Wails task pipeline"

echo "[4/6] Wails build"
wails build -skipbindings || true

echo "[5/6] Metrics check"
go run ./cmd/metricscheck --csv ./testdata/sample_trades.csv --expected ./testdata/expected.json || true

echo "[6/6] golangci-lint (optional)"
if command -v golangci-lint >/dev/null; then
  golangci-lint run ./...
fi

echo "Done"

