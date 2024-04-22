// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createuser = `-- name: Createuser :one
INSERT INTO users(id,email,passwd,username,created_at) VALUES($1,$2,$3,$4,$5)
RETURNING username, created_at
`

type CreateuserParams struct {
	ID        uuid.UUID
	Email     string
	Passwd    []byte
	Username  string
	CreatedAt time.Time
}

type CreateuserRow struct {
	Username  string
	CreatedAt time.Time
}

func (q *Queries) Createuser(ctx context.Context, arg CreateuserParams) (CreateuserRow, error) {
	row := q.db.QueryRowContext(ctx, createuser,
		arg.ID,
		arg.Email,
		arg.Passwd,
		arg.Username,
		arg.CreatedAt,
	)
	var i CreateuserRow
	err := row.Scan(&i.Username, &i.CreatedAt)
	return i, err
}

const find_user_email = `-- name: Find_user_email :one
SELECT id, email, passwd, username, created_at FROM users WHERE email=$1
`

func (q *Queries) Find_user_email(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, find_user_email, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Passwd,
		&i.Username,
		&i.CreatedAt,
	)
	return i, err
}

const find_user_name = `-- name: Find_user_name :one
SELECT id, email, passwd, username, created_at FROM users WHERE username=$1
`

func (q *Queries) Find_user_name(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, find_user_name, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Passwd,
		&i.Username,
		&i.CreatedAt,
	)
	return i, err
}

const if_email = `-- name: If_email :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email=$1
) AS value_exists
`

func (q *Queries) If_email(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRowContext(ctx, if_email, email)
	var value_exists bool
	err := row.Scan(&value_exists)
	return value_exists, err
}

const if_username = `-- name: If_username :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username=$1
) AS value_exists
`

func (q *Queries) If_username(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRowContext(ctx, if_username, username)
	var value_exists bool
	err := row.Scan(&value_exists)
	return value_exists, err
}
