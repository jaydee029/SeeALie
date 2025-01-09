// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Revoked struct {
	Token     string
	RevokedAt pgtype.Timestamp
}

type Session struct {
	SessionID pgtype.UUID
	UserID    pgtype.UUID
	Jwt       string
	ExpiresAt pgtype.Timestamp
}

type User struct {
	ID        pgtype.UUID
	Email     string
	Passwd    []byte
	Username  string
	CreatedAt pgtype.Timestamp
}
