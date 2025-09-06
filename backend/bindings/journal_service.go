package bindings

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"changeme/backend/models"
	"changeme/backend/repo"
	"changeme/backend/services"
)

// JournalService exposes journal operations to the Wails frontend.
type JournalService struct {
	Repo       repo.TradeRepo
	AppVersion string
}

// Ping returns the app version to demonstrate an end-to-end call.
func (s *JournalService) Ping(ctx context.Context) (string, error) {
	if s.AppVersion == "" {
		return "dev", nil
	}
	return s.AppVersion, nil
}

// GetAnalytics calculates analytics for the given query. TODO(impl)
func (s *JournalService) GetAnalytics(ctx context.Context, q models.Query) (models.AnalyticsSummary, error) {
	// Stub implementation to show wiring end-to-end
	trades, err := s.Repo.List(ctx, q)
	if err != nil {
		return models.AnalyticsSummary{}, err
	}
	// Build a naive returns slice using exit - entry (ignoring fees, qty)
	returns := make([]float64, 0, len(trades))
	for _, t := range trades {
		if t.ExitPrice != nil {
			r := (*t.ExitPrice-t.EntryPrice)*t.Quantity - t.Fees
			if strings.ToLower(t.Side) == "short" {
				r = (t.EntryPrice-*t.ExitPrice)*t.Quantity - t.Fees
			}
			returns = append(returns, r)
		}
	}
	// Equity curve as cumulative sum of returns
	equity := make([]float64, len(returns))
	var sum float64
	for i, r := range returns {
		sum += r
		equity[i] = sum
	}
	out := models.AnalyticsSummary{
		WinRate:      services.WinRate(returns),
		ProfitFactor: services.ProfitFactor(returns),
		MaxDD:        services.MaxDrawdown(equity),
		Sharpe:       services.Sharpe(returns),
		Sortino:      services.Sortino(returns),
		Expectancy:   services.Expectancy(returns),
	}
	return out, nil
}

// ImportCSV imports trades from a CSV payload string. TODO(impl): dedupe on (symbol, entry_time, exit_time, entry_price, qty)
func (s *JournalService) ImportCSV(ctx context.Context, csvPayload string) (models.ImportReport, error) {
	r := csv.NewReader(strings.NewReader(csvPayload))
	r.FieldsPerRecord = -1
	// Expect header: symbol,side,entry_time,exit_time,entry_price,exit_price,qty,fees,notes
	// Times are ISO8601 in UTC.
	var report models.ImportReport
	_, err := r.Read()
	if err != nil {
		return report, fmt.Errorf("read header: %w", err)
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			report.Errors = append(report.Errors, err.Error())
			continue
		}
		if len(rec) < 9 {
			report.Errors = append(report.Errors, "invalid record length")
			continue
		}
		entryTime, err := time.Parse(time.RFC3339, rec[2])
		if err != nil {
			report.Errors = append(report.Errors, "bad entry_time")
			continue
		}
		var exitTime *time.Time
		if rec[3] != "" {
			et, err := time.Parse(time.RFC3339, rec[3])
			if err != nil {
				report.Errors = append(report.Errors, "bad exit_time")
				continue
			}
			exitTime = &et
		}
		entryPrice, _ := strconv.ParseFloat(rec[4], 64)
		var exitPrice *float64
		if rec[5] != "" {
			ep, _ := strconv.ParseFloat(rec[5], 64)
			exitPrice = &ep
		}
		qty, _ := strconv.ParseFloat(rec[6], 64)
		fees, _ := strconv.ParseFloat(rec[7], 64)

		t := models.Trade{
			Symbol:     rec[0],
			Side:       rec[1],
			EntryTime:  entryTime.UTC(),
			ExitTime:   exitTime,
			EntryPrice: entryPrice,
			ExitPrice:  exitPrice,
			Quantity:   qty,
			Fees:       fees,
			Notes:      rec[8],
		}
		if err := s.Repo.Upsert(ctx, t); err != nil {
			report.Errors = append(report.Errors, err.Error())
		} else {
			report.Imported++
		}
	}
	return report, nil
}

// ListTrades returns trades matching the query. Useful for populating the table in UI.
func (s *JournalService) ListTrades(ctx context.Context, q models.Query) ([]models.Trade, error) {
	if s.Repo == nil {
		return nil, fmt.Errorf("repo not configured")
	}
	return s.Repo.List(ctx, q)
}
