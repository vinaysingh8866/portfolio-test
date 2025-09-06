package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"changeme/backend/services"
)

type Expected struct {
	WinRate      float64 `json:"winRate"`
	ProfitFactor float64 `json:"profitFactor"`
	MaxDD        float64 `json:"maxDD"`
	Sharpe       float64 `json:"sharpe"`
	Sortino      float64 `json:"sortino"`
	Expectancy   float64 `json:"expectancy"`
}

func main() {
	csvPath := flag.String("csv", "", "input CSV path")
	expPath := flag.String("expected", "", "expected json path")
	flag.Parse()
	if *csvPath == "" || *expPath == "" {
		fmt.Println("usage: metricscheck --csv sample.csv --expected expected.json")
		os.Exit(2)
	}
	returns := readReturns(*csvPath)
	equity := make([]float64, len(returns))
	var s float64
	for i, r := range returns {
		s += r
		equity[i] = s
	}
	got := Expected{
		WinRate:      services.WinRate(returns),
		ProfitFactor: services.ProfitFactor(returns),
		MaxDD:        services.MaxDrawdown(equity),
		Sharpe:       services.Sharpe(returns),
		Sortino:      services.Sortino(returns),
		Expectancy:   services.Expectancy(returns),
	}
	exp := Expected{}
	f, _ := os.Open(*expPath)
	defer f.Close()
	json.NewDecoder(f).Decode(&exp)
	fmt.Printf("got: %+v\nexpected: %+v\n", got, exp)
}

func readReturns(path string) []float64 {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	_, _ = r.Read()
	var returns []float64
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		side := rec[1]
		entry, _ := strconv.ParseFloat(rec[4], 64)
		if rec[5] == "" {
			continue
		}
		exit, _ := strconv.ParseFloat(rec[5], 64)
		qty, _ := strconv.ParseFloat(rec[6], 64)
		fees, _ := strconv.ParseFloat(rec[7], 64)
		pnl := 0.0
		if side == "short" {
			pnl = (entry - exit) * qty
		} else {
			pnl = (exit - entry) * qty
		}
		pnl -= fees
		returns = append(returns, pnl)
	}
	return returns
}
