param()
$ErrorActionPreference = "Stop"

Write-Host "[1/6] Go tests"
go test ./...

Write-Host "[2/6] Frontend (npm)"
npm --prefix frontend ci
npm --prefix frontend run build

Write-Host "[4/6] Wails build"
wails build -skipbindings

Write-Host "[5/6] Metrics check"
go run ./cmd/metricscheck --csv ./testdata/sample_trades.csv --expected ./testdata/expected.json

Write-Host "[6/6] golangci-lint (optional)"
golangci-lint run ./... 2>$null

Write-Host "Done"

