-- name: CreateAccount :one
INSERT INTO "accounts" ("user_id", 
                      "balance",
                      "currency")
VALUES (sqlc.arg(user_id), sqlc.arg(balance), sqlc.arg(currency)) RETURNING *;

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
where user_id = sqlc.arg(user_id)
ORDER BY id DESC
LIMIT sqlc.arg(limits)
OFFSET sqlc.arg(offsets);

-- name: UpdateAccount :one
UPDATE "accounts"
SET balance = balance + sqlc.arg(balance)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM "accounts"
WHERE id = $1;