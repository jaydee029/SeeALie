// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Revoked struct {
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
