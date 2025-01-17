package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func GenpgtypeUUID(id string) (pgtype.UUID, error) {

	var pgid pgtype.UUID
	err := pgid.Scan(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgid, nil

}

func GenpgtypeTimestamp(t time.Time) (pgtype.Timestamp, error) {

	var pgt pgtype.Timestamp
	err := pgt.Scan(t)
	if err != nil {
		return pgtype.Timestamp{}, err
	}
	return pgt, nil

}
