-- name: CreateTrade :one
INSERT INTO trades (id, user_id, contract, open_quantity, open_price)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTrade :one
SELECT * FROM trades WHERE id = ?;

-- name: GetUserTrades :many
SELECT * FROM trades WHERE user_id = ?;

-- name: GetOpenTrades :many
SELECT * FROM trades WHERE user_id = ? AND open_quantity > 0;

-- name: GetTradesForUser :many
SELECT * FROM trades WHERE user_id = ?;

-- name: GetClosedTrades :many
SELECT * FROM trades WHERE user_id = ? AND sold_price IS NOT NULL;

-- name: UpdateTrade :one
UPDATE trades SET open_quantity = ?, closed_quantity = ?, average_close_price = ?, realized_profit = ? WHERE id = ?
RETURNING *;