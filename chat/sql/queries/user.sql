-- name: Get_username :one
SELECT username FROM users WHERE id=$1;