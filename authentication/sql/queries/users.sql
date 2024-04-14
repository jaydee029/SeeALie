-- name: Createuser :one
INSERT INTO users(id,email,passwd,username,created_at) VALUES($1,$2,$3,$4,$5)
RETURNING *;

-- name: If_email :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email=$1
) AS value_exists;

-- name: If_username :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username=$1
) AS value_exists;