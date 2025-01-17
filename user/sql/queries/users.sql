-- name: Createuser :one
INSERT INTO users(id,email,passwd,username,created_at) VALUES($1,$2,$3,$4,$5)
RETURNING email,username, created_at;

-- name: IfEmail :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email=$1
) AS value_exists;

-- name: IfUsername :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username=$1
) AS value_exists;

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email=$1;

-- name: FindUserByUsername :one
SELECT * FROM users WHERE username=$1;

-- name: InsertSession :one
INSERT INTO sessions(session_id, user_id, jwt, expires_at) VALUES($1, $2, $3, $4)
RETURNING session_id, expires_at;

-- name: FindSessionByid :one

SELECT EXISTS (
    SELECT 1 FROM sessions WHERE user_id=$1 AND expires_at > NOW()
) AS value_exists;