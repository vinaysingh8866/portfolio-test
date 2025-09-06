package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"

	"github.com/google/uuid"

	"changeme/backend/models"
)

// TradeRepo defines how trades are persisted and queried.
type TradeRepo interface {
	Upsert(ctx context.Context, t models.Trade) error
	List(ctx context.Context, q models.Query) ([]models.Trade, error)
	Delete(ctx context.Context, id string) error
}

// SQLiteRepo is a default SQLite implementation.
type SQLiteRepo struct {
	DB *sql.DB
}

// NewSQLiteRepo opens (or creates) a SQLite database at the given path.
func NewSQLiteRepo(dbPath string) (*SQLiteRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &SQLiteRepo{DB: db}, nil
}

// EnsureID sets an ID if missing.
func ensureID(id string) string {
	if id != "" {
		return id
	}
	return uuid.NewString()
}

// Upsert inserts or updates a trade based on its ID.
func (r *SQLiteRepo) Upsert(ctx context.Context, t models.Trade) error {
	if r == nil || r.DB == nil {
		return errors.New("repo not initialized")
	}
	t.ID = ensureID(t.ID)
	// Normalize to UTC
	t.EntryTime = t.EntryTime.UTC()
	if t.ExitTime != nil {
		et := t.ExitTime.UTC()
		t.ExitTime = &et
	}
	now := time.Now().UTC()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now

	_, err := r.DB.ExecContext(ctx, `
        INSERT INTO trades (
            id, symbol, side, entry_time, exit_time, entry_price, exit_price, qty, fees, notes, created_at, updated_at
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        ) ON CONFLICT(id) DO UPDATE SET
            symbol=excluded.symbol,
            side=excluded.side,
            entry_time=excluded.entry_time,
            exit_time=excluded.exit_time,
            entry_price=excluded.entry_price,
            exit_price=excluded.exit_price,
            qty=excluded.qty,
            fees=excluded.fees,
            notes=excluded.notes,
            updated_at=excluded.updated_at
    `,
		t.ID, t.Symbol, t.Side, t.EntryTime,
		t.ExitTime, t.EntryPrice, t.ExitPrice,
		t.Quantity, t.Fees, t.Notes, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

// List returns trades matching the query filters.
func (r *SQLiteRepo) List(ctx context.Context, q models.Query) ([]models.Trade, error) {
	if r == nil || r.DB == nil {
		return nil, errors.New("repo not initialized")
	}
	sqlStr := `SELECT id, symbol, side, entry_time, exit_time, entry_price, exit_price, qty, fees, notes, created_at, updated_at FROM trades WHERE 1=1`
	args := []any{}
	if q.Symbol != "" {
		sqlStr += " AND symbol = ?"
		args = append(args, q.Symbol)
	}
	if q.Side != "" {
		sqlStr += " AND side = ?"
		args = append(args, q.Side)
	}
	if q.StartTime != nil {
		sqlStr += " AND entry_time >= ?"
		args = append(args, q.StartTime.UTC())
	}
	if q.EndTime != nil {
		sqlStr += " AND entry_time < ?"
		args = append(args, q.EndTime.UTC())
	}
	sqlStr += " ORDER BY entry_time ASC"
	if q.Limit > 0 {
		sqlStr += " LIMIT ?"
		args = append(args, q.Limit)
	}
	if q.Offset > 0 {
		sqlStr += " OFFSET ?"
		args = append(args, q.Offset)
	}

	rows, err := r.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Trade
	for rows.Next() {
		var t models.Trade
		var exitTime sql.NullTime
		var exitPrice sql.NullFloat64
		if err := rows.Scan(
			&t.ID, &t.Symbol, &t.Side, &t.EntryTime, &exitTime, &t.EntryPrice, &exitPrice,
			&t.Quantity, &t.Fees, &t.Notes, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if exitTime.Valid {
			et := exitTime.Time
			t.ExitTime = &et
		}
		if exitPrice.Valid {
			ep := exitPrice.Float64
			t.ExitPrice = &ep
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// Delete removes a trade by ID.
func (r *SQLiteRepo) Delete(ctx context.Context, id string) error {
	if r == nil || r.DB == nil {
		return errors.New("repo not initialized")
	}
	_, err := r.DB.ExecContext(ctx, `DELETE FROM trades WHERE id = ?`, id)
	return err
}
