-- name: Search_user :one
SELECT request_by FROM connections WHERE connection_id=$1;

-- name: Add_contact :one
INSERT INTO friends(followed_by, followed, room_id ,connected_at) VALUES ($1,$2,$3,$4)
RETURNING connected_at;

