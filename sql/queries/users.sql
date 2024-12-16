-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users CASCADE;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByToken :one
SELECT users.* FROM users
WHERE users.id = (
    SELECT user_id FROM refresh_tokens WHERE refersh_tokens.token = $1
);
