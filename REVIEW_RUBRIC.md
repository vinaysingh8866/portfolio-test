## Review Rubric (100 pts)

### Correctness (40)
- MaxDrawdown percent‑based O(n) (15)
- ProfitFactor, Expectancy, Sharpe, Sortino (10)
- CSV import: UTC storage + dedupe behavior (10)
- Passing unit tests and CI (5)

### Code Quality (25)
- Clear, readable Go code and naming (10)
- Error handling and edge cases (5)
- Reasonable precision handling with decimal (5)
- Simple, testable functions (5)

### Performance (15)
- MaxDD O(n) and no quadratic hot paths (10)
- Handles large CSV reasonably (5)

### UX & Wiring (10)
- Filters refresh analytics & table (5)
- Equity chart renders minimal line (5)

### Testing (10)
- Adds/maintains unit tests for analytics and timezone (10)

### Timebox & Submission
- Expected effort: 6–8 hours (cap at 10 hours).
- Submission window: 72 hours from assignment (flexible to 5 days on request).
- Stretch items (decimal precision, perf on large CSV, extra tests/polish) are optional. Passing bar: CI green, analytics correct (incl. percent MaxDD), CSV dedupe + UTC correct, filters + table + equity chart minimally working.


