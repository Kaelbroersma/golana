-- +goose Up

CREATE TABLE trades (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    contract TEXT NOT NULL REFERENCES contracts(contract_id),
    side TEXT NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    bought_price DECIMAL(10, 2),
    sold_price DECIMAL(10, 2),
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down

DROP TABLE trades;