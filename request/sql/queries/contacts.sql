-- name: Addcontact :one
INSERT INTO contacts(id, username, room_id ,connected_on) VALUES ($1,$2,$3,$4)
RETURNING username, connected_on;

