-- name: CreateUser :one

INSERT INTO users (id, email, name, hashed_password)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetUserByID :one

SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one

SELECT * FROM users WHERE email = ?;

-- name: UpdateUserBalances :one

UPDATE users
SET buying_power = ?, exposure = ?
WHERE id = ?
RETURNING *;