-- trades table
CREATE TABLE IF NOT EXISTS trades (
  id TEXT PRIMARY KEY,
  symbol TEXT NOT NULL,
  side TEXT NOT NULL CHECK(side IN ('long','short')),
  entry_time TIMESTAMP NOT NULL,
  exit_time TIMESTAMP NULL,
  entry_price REAL NOT NULL,
  exit_price REAL NULL,
  qty REAL NOT NULL,
  fees REAL NOT NULL DEFAULT 0,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- indexes for query performance
CREATE INDEX IF NOT EXISTS idx_trades_entry_time ON trades(entry_time);
CREATE INDEX IF NOT EXISTS idx_trades_symbol ON trades(symbol);


