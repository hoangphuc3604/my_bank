-- name: CreateTransfer :execresult
INSERT INTO transfers (from_account, to_account, amount)
VALUES (?, ?, ?);

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ?;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id ASC
LIMIT ? OFFSET ?;

-- name: UpdateTransfer :exec
UPDATE transfers
SET amount = ?
WHERE id = ?;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = ?;
