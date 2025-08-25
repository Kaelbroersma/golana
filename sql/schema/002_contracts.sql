-- +goose Up

CREATE TABLE contracts (
    contract_id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    websocket_open BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down

DROP TABLE contracts;