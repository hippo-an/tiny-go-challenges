-- name: CreateUser :one
insert into users (
    username,
    hashed_password,
    full_name,
    email
) values (
    sqlc.arg(username),
    sqlc.arg(hashed_password),
    sqlc.arg(full_name),
    sqlc.arg(email)
) RETURNING *;

-- name: GetUser :one
select * from users
where id = sqlc.arg(user_id)
limit 1;


-- name: GetUserByUsername :one
select * from users
where username = sqlc.arg(username)
limit 1;