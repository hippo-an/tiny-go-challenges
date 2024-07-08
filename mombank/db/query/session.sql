-- name: CreateSession :one
insert into sessions (
    id,
    user_id,
    refersh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at
) values (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(refersh_token),
    sqlc.arg(user_agent),
    sqlc.arg(client_ip),
    sqlc.arg(is_blocked),
    sqlc.arg(expires_at)
) RETURNING *;

-- name: GetSession :one
select * from sessions
where id = sqlc.arg(id)
limit 1;