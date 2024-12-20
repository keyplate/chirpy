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
ORDER BY chirps.created_at ASC;

-- name: GetAllChirpsByUserID :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY chirps.created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE chirps.id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps 
WHERE chirps.id = $1;
