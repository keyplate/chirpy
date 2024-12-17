// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Email,
		arg.HashedPassword,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
TRUNCATE TABLE users CASCADE
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUserByToken = `-- name: GetUserByToken :one
SELECT users.id, users.created_at, users.updated_at, users.email, users.hashed_password, users.is_chirpy_red FROM users
WHERE users.id = (
    SELECT user_id FROM refresh_tokens WHERE refresh_tokens.token = $1
)
`

func (q *Queries) GetUserByToken(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByToken, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateIsChirpyRedUserTrue = `-- name: UpdateIsChirpyRedUserTrue :one
UPDATE users SET
    is_chirpy_red = true
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

func (q *Queries) UpdateIsChirpyRedUserTrue(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, updateIsChirpyRedUserTrue, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUserEmailPassword = `-- name: UpdateUserEmailPassword :one
UPDATE users SET
    email = $1,
    hashed_password = $2
WHERE id = $3
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type UpdateUserEmailPasswordParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

func (q *Queries) UpdateUserEmailPassword(ctx context.Context, arg UpdateUserEmailPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserEmailPassword, arg.Email, arg.HashedPassword, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
