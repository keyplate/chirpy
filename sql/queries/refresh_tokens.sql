-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE refresh_tokens.token = $1;

-- name: MarkRevoked :one
UPDATE refresh_tokens
SET updated_at = $1, 
    revoked_at =$2
WHERE token = $3
RETURNING *;
