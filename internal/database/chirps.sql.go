// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chirps.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createChirps = `-- name: CreateChirps :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid (), NOw(), NOW(), $1, $2
)
RETURNING id, created_at, updated_at, body, user_id
`

type CreateChirpsParams struct {
	Body   string
	UserID uuid.UUID
}

func (q *Queries) CreateChirps(ctx context.Context, arg CreateChirpsParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirps, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const deleteChirpByID = `-- name: DeleteChirpByID :exec
DELETE FROM chirps 
WHERE id = $1
`

func (q *Queries) DeleteChirpByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteChirpByID, id)
	return err
}

const getChirp = `-- name: GetChirp :one
SELECT id, created_at, updated_at, body, user_id FROM chirps WHERE id = $1
`

func (q *Queries) GetChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const getChirps = `-- name: GetChirps :many
SELECT id, created_at, updated_at, body, user_id FROM chirps ORDER BY created_at ASC
`

func (q *Queries) GetChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
