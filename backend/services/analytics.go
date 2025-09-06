package services

import "math"

// TODO(impl): Switch to a decimal library for P&L aggregation to improve precision.

// MaxDrawdown calculates maximum drawdown from an equity curve.
// NOTE: This is an intentionally naive O(n^2) implementation. TODO: optimize to O(n).
func MaxDrawdown(equity []float64) float64 {
	maxDD := 0.0
	for i := 0; i < len(equity); i++ {
		for j := i + 1; j < len(equity); j++ {
			dd := (equity[i] - equity[j])
			if dd > maxDD {
				maxDD = dd
			}
		}
	}
	return maxDD
}

// ProfitFactor = gross profits / gross losses (absolute).
func ProfitFactor(returns []float64) float64 {
	gp, gl := 0.0, 0.0
	for _, r := range returns {
		if r > 0 {
			gp += r
		} else if r < 0 {
			gl += -r
		}
	}
	if gl == 0 {
		if gp == 0 {
			return 0
		}
		return math.Inf(1)
	}
	return gp / gl
}

// WinRate = wins / total.
func WinRate(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}
	wins := 0.0
	for _, r := range returns {
		if r > 0 {
			wins += 1
		}
	}
	return wins / float64(len(returns))
}

// Expectancy = average return per trade.
func Expectancy(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range returns {
		sum += r
	}
	return sum / float64(len(returns))
}

// Sharpe ratio using population stddev, risk-free = 0 (placeholder).
func Sharpe(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}
	mean := Expectancy(returns)
	var ss float64
	for _, r := range returns {
		d := r - mean
		ss += d * d
	}
	sd := math.Sqrt(ss / float64(len(returns)))
	if sd == 0 {
		return 0
	}
	return mean / sd
}

// Sortino ratio using downside deviation only.
func Sortino(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}
	mean := Expectancy(returns)
	var downsideSS float64
	n := 0.0
	for _, r := range returns {
		if r < 0 {
			d := r
			downsideSS += d * d
			n += 1
		}
	}
	if n == 0 {
		return math.Inf(1)
	}
	dd := math.Sqrt(downsideSS / n)
	if dd == 0 {
		return math.Inf(1)
	}
	return mean / dd
}
