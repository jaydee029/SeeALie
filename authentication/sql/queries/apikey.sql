-- name: CreateApiKey :one
INSERT INTO apikey(id) VALUES($1)
RETURNING *;

-- name: VerifyApiKey :one
SELECT id, username
FROM users
WHERE EXISTS (SELECT 1 FROM apikey WHERE id = users.id AND api_key = $1);
