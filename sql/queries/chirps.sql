-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY chirps.created_at;
