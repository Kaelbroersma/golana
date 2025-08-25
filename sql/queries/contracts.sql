-- name: CreateContract :one
INSERT INTO contracts (
    contract_id,
    name
) VALUES (
    ?, ?
) RETURNING *;

-- name: GetContractByID :one
SELECT * FROM contracts WHERE contract_id = ?;

-- name: GetContractByName :one
SELECT * FROM contracts WHERE name = ?;