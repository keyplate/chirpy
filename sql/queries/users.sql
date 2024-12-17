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
    SELECT user_id FROM refresh_tokens WHERE refresh_tokens.token = $1
);

-- name: UpdateUserEmailPassword :one
UPDATE users SET
    email = $1,
    hashed_password = $2
WHERE id = $3
RETURNING *;

-- name: UpdateIsChirpyRedUserTrue :one
UPDATE users SET
    is_chirpy_red = true
WHERE id = $1
RETURNING *;
