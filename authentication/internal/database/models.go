// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Friend struct {
	FollowedBy  uuid.UUID
	Followed    uuid.UUID
	ConnectedAt time.Time
}

type IDName struct {
	UserID   uuid.UUID
	Username sql.NullString
}

type Revokedt struct {
	Token     string
	RevokedAt time.Time
}

type User struct {
	ID        uuid.UUID
	Email     string
	Passwd    []byte
	Username  string
	CreatedAt time.Time
}
