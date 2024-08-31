// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const get_username = `-- name: Get_username :one
SELECT username FROM users WHERE id=$1
`

func (q *Queries) Get_username(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, get_username, id)
	var username string
	err := row.Scan(&username)
	return username, err
}
