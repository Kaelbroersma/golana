-- name: CreateSymbol :one

-- name: GetSymbol :one

-- name: GetAllSymbols :many

-- name: UpdateSymbol :one

-- name: DeleteSymbol :exec

-- name OpenSocket :one
UPDATE contracts SET websocket_open = true WHERE contract_id = ?;

-- name: CloseSocket :one
UPDATE contracts SET websocket_open = false WHERE contract_id = ?;