// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: friend.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const add_contact = `-- name: Add_contact :one
INSERT INTO friends(followed_by, followed, room_id ,connected_at) VALUES ($1,$2,$3,$4)
RETURNING connected_at
`

type Add_contactParams struct {
	FollowedBy  string
	Followed    string
	RoomID      uuid.UUID
	ConnectedAt time.Time
}

func (q *Queries) Add_contact(ctx context.Context, arg Add_contactParams) (time.Time, error) {
	row := q.db.QueryRowContext(ctx, add_contact,
		arg.FollowedBy,
		arg.Followed,
		arg.RoomID,
		arg.ConnectedAt,
	)
	var connected_at time.Time
	err := row.Scan(&connected_at)
	return connected_at, err
}

const search_user = `-- name: Search_user :one
SELECT request_by FROM connections WHERE connection_id=$1
`

func (q *Queries) Search_user(ctx context.Context, connectionID uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, search_user, connectionID)
	var request_by string
	err := row.Scan(&request_by)
	return request_by, err
}
