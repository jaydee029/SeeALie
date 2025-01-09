-- name: Createuser :one
INSERT INTO users(id,email,passwd,username,created_at) VALUES($1,$2,$3,$4,$5)
RETURNING username, created_at;

-- name: Is_email :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email=$1
) AS value_exists;

-- name: Is_username :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username=$1
) AS value_exists;

-- name: Find_user_email :one
SELECT * FROM users WHERE email=$1;

-- name: Find_user_name :one
SELECT * FROM users WHERE username=$1;