-- name: CreateAccount :one
INSERT INTO "accounts" ("user_id", 
                        "owner",
                      "balance",
                      "currency")
VALUES (sqlc.arg(user_id), sqlc.arg(owner), sqlc.arg(balance), sqlc.arg(currency)) RETURNING *;

-- name: GetAccount :one
SELECT *
FROM "accounts"
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT *
FROM "accounts"
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE ;

-- name: ListAccounts :many
SELECT *
FROM "accounts"
ORDER BY id DESC
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE "accounts"
SET balance = balance + sqlc.arg(balance)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM "accounts"
WHERE id = $1;