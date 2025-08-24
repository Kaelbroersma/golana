-- +goose Up

CREATE TABLE trades (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    symbol TEXT NOT NULL REFERENCES symbols(symbol),
    side TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    bought_price DECIMAL(10, 2) NOT NULL,
    sold_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down

DROP TABLE trades;