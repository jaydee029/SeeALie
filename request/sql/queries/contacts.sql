-- name: Add_connection_Info :one
INSERT INTO connections (request_by,request_to,connection_id,created_at) VALUES ($1,$2,$3,$4)
RETURNING created_at;

-- name: Delete_connection_Info :exec
DELETE FROM connections WHERE connection_id=$1;

-- name: User_exists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username=$1
) AS value_exists;

-- name: Get_user_email :one
SELECT email FROM users WHERE username=$1;

-- name: GetRequestInfo :one
SELECT email FROM users WHERE username=$1;