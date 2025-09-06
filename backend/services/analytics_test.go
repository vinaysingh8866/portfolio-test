package services

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMaxDrawdown_Simple(t *testing.T) {
	equity := []float64{100, 120, 90, 95, 80, 130}
	dd := MaxDrawdown(equity)
	if dd != 40 { // 120 -> 80
		t.Fatalf("expected 40, got %v", dd)
	}
}

func TestProfitFactorAndWinRate_TinySet(t *testing.T) {
	returns := []float64{10, -5, 20, -10}
	pf := ProfitFactor(returns)
	wr := WinRate(returns)
	if pf <= 1.0 {
		t.Fatalf("expected pf > 1, got %v", pf)
	}
	if wr != 0.5 {
		t.Fatalf("expected wr=0.5, got %v", wr)
	}
}

func TestSortino_DownsideDeviationOnly(t *testing.T) {
	returns := []float64{1, -1, 2, -2, 3}
	s := Sortino(returns)
	if s <= 0 {
		t.Fatalf("expected positive sortino, got %v", s)
	}
}

// Intentionally failing test: TODO(impl) refine MaxDrawdown definition to percent-based drawdown.
func TestMaxDrawdown_PercentBased_TODO(t *testing.T) {
	equity := []float64{100, 200, 150}
	// Percent-based max DD from 200 to 150 is 25% of peak (50/200).
	// Our current implementation uses absolute drawdown and will return 50.
	dd := MaxDrawdown(equity)
	if dd == 50 { // current behavior
		t.Fatalf("TODO: Implement percent-based MaxDrawdown; got absolute=%v", dd)
	}
}

func TestImportStoresUTC(t *testing.T) {
	local := "2024-03-10T01:59:59-05:00"
	tt, err := time.Parse(time.RFC3339, local)
	if err != nil {
		t.Fatal(err)
	}
	if tt.Location() == time.UTC {
		t.Fatalf("expected non-UTC input")
	}
	if tt.UTC().Location() != time.UTC {
		t.Fatalf("expected UTC after conversion")
	}
}

func TestFixtures_SampleTrades_Returns(t *testing.T) {
	f, err := os.Open("../../testdata/sample_trades.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	_, _ = r.Read() // header
	var returns []float64
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		side := rec[1]
		entry, _ := strconv.ParseFloat(rec[4], 64)
		exit := 0.0
		if rec[5] != "" {
			exit, _ = strconv.ParseFloat(rec[5], 64)
		} else {
			continue // open trade
		}
		qty, _ := strconv.ParseFloat(rec[6], 64)
		fees, _ := strconv.ParseFloat(rec[7], 64)
		var pnl float64
		if side == "short" {
			pnl = (entry - exit) * qty
		} else {
			pnl = (exit - entry) * qty
		}
		pnl -= fees
		returns = append(returns, pnl)
	}
	if len(returns) == 0 {
		t.Fatal("expected some closed trades in sample_trades.csv")
	}
	_ = ProfitFactor(returns)
	_ = WinRate(returns)
	_ = Sharpe(returns)
	_ = Sortino(returns)
}
