-- name: CreateTrade :one
INSERT INTO trades (id, user_id, contract, quantity, open_price, close_price)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTrade :one
SELECT * FROM trades WHERE id = ?;

-- name: GetUserTrades :many
SELECT * FROM trades WHERE user_id = ?;

-- name: GetOpenTrades :many
SELECT * FROM trades WHERE user_id = ? AND sold_price IS NULL;

-- name: GetClosedTrades :many
SELECT * FROM trades WHERE user_id = ? AND sold_price IS NOT NULL;

-- name: UpdateTrade :one
UPDATE trades SET quantity = ?, open_price = ?, close_price = ? WHERE id = ?
RETURNING *;
