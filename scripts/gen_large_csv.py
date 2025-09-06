#!/usr/bin/env python3
import argparse
import csv
import random
from datetime import datetime, timedelta, timezone


SYMBOLS = ["AAPL", "MSFT", "TSLA", "NVDA", "AMZN", "GOOG", "META"]
SIDES = ["long", "short"]


def gen_row(start_dt: datetime, idx: int) -> list[str]:
    random.seed(idx)
    symbol = random.choice(SYMBOLS)
    side = random.choice(SIDES)
    # Entry time with some jitter; include some same-second collisions
    entry_time = start_dt + timedelta(minutes=random.randint(0, 1200))
    if idx % 50 == 0:
        # Same-second timestamp edge case
        entry_time = start_dt + timedelta(minutes=idx // 2)
    # DST boundary example (US 2024-03-10): create one around that time
    if idx == 123:
        entry_time = datetime(2024, 3, 10, 6, 59, 59, tzinfo=timezone.utc)
    # Some trades left open (missing exit_price)
    open_trade = (idx % 17 == 0)

    entry_price = round(random.uniform(10, 500), 2)
    # Move price up or down depending on side
    price_move = round(random.uniform(-0.08, 0.1) * entry_price, 2)
    if side == "long":
        exit_price_val = max(0.01, entry_price + price_move)
    else:
        exit_price_val = max(0.01, entry_price - price_move)

    qty = random.choice([1, 2, 3, 5, 10])
    fees = round(random.uniform(0, 3), 2)
    notes = "auto"

    exit_time = entry_time + timedelta(minutes=random.randint(1, 360))

    entry_iso = entry_time.astimezone(timezone.utc).isoformat().replace("+00:00", "Z")
    exit_iso = exit_time.astimezone(timezone.utc).isoformat().replace("+00:00", "Z")

    if open_trade:
        return [symbol, side, entry_iso, "", f"{entry_price}", "", f"{qty}", f"{fees}", notes]
    return [symbol, side, entry_iso, exit_iso, f"{entry_price}", f"{round(exit_price_val,2)}", f"{qty}", f"{fees}", notes]


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--rows", type=int, default=3000)
    ap.add_argument("--out", type=str, default="testdata/large_trades.csv")
    args = ap.parse_args()

    start_dt = datetime(2024, 1, 1, 9, 30, tzinfo=timezone.utc)
    with open(args.out, "w", newline="") as f:
        w = csv.writer(f)
        w.writerow(["symbol", "side", "entry_time", "exit_time", "entry_price", "exit_price", "qty", "fees", "notes"])
        for i in range(args.rows):
            w.writerow(gen_row(start_dt, i))

    print(f"wrote {args.rows} rows to {args.out}")


if __name__ == "__main__":
    main()


