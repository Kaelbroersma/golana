-- +goose Up

CREATE TABLE symbols (
    symbol TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    websocket_open BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down

DROP TABLE symbols;