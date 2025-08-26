-- +goose Up
CREATE TABLE users (
    id TEXT NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    buying_power DECIMAL(10, 2) NOT NULL DEFAULT 1000,
    exposure DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime('now')),    
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE users;