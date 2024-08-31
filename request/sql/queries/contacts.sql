-- name: Add_contact :one
INSERT INTO friends(followed_by, followed, room_id ,connected_on) VALUES ($1,$2,$3,$4)
RETURNING connected_on;

-- name: Add_connection_Info :one
INSERT INTO connections (request_by,connection_id) VALUES ($1,$2);

-- name: Search_user :one
SELECT request_by FROM connections WHERE connection_id=$1;

-- name: Delete_connection_Info :exec
DELETE * FROM connections WHERE connection_id=$1;

