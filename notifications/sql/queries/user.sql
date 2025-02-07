-- name: GetEmail :one
SELECT email FROM users WHERE username=$1;
