-- name: Addcontact :one
INSERT INTO contacts(id, username, connected_on) VALUES ($1,$2,$3)
RETURNING username, connected_on;

