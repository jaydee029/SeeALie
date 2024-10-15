-- name: RevokeToken :one
INSERT INTO revoked(token,revoked_at) VALUES($1,$2)
RETURNING *;

-- name: VerifyRevoke :one
SELECT EXISTS (SELECT 1 FROM revoked WHERE token=$1) AS value_exists;
