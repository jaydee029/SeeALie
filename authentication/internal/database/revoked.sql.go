// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: revoked.sql

package database

import (
	"context"
	"time"
)

const revokeToken = `-- name: RevokeToken :one
INSERT INTO revoked(token,revoked_at) VALUES($1,$2)
RETURNING token, revoked_at
`

type RevokeTokenParams struct {
	Token     string
	RevokedAt time.Time
}

func (q *Queries) RevokeToken(ctx context.Context, arg RevokeTokenParams) (Revoked, error) {
	row := q.db.QueryRowContext(ctx, revokeToken, arg.Token, arg.RevokedAt)
	var i Revoked
	err := row.Scan(&i.Token, &i.RevokedAt)
	return i, err
}

const verifyRevoke = `-- name: VerifyRevoke :one
SELECT EXISTS (SELECT 1 FROM revoked WHERE token=$1) AS value_exists
`

func (q *Queries) VerifyRevoke(ctx context.Context, token string) (bool, error) {
	row := q.db.QueryRowContext(ctx, verifyRevoke, token)
	var value_exists bool
	err := row.Scan(&value_exists)
	return value_exists, err
}
