-- +goose Up

CREATE TABLE trades (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    contract TEXT NOT NULL REFERENCES contracts(contract_id),
    open_quantity DECIMAL(10, 2) NOT NULL,
    closed_quantity DECIMAL(10,2),
    open_price DECIMAL(10, 2),
    average_close_price DECIMAL(10, 2),
    unrealized_profit DECIMAL(10, 2),
    realized_profit DECIMAL(10, 2),
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down

DROP TABLE trades;