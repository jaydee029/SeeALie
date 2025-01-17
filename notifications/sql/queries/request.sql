-- name: DequeRequests :many
SELECT request_by, request_to, connection_id FROM connections WHERE status_sent=PENDING AND sent_attempts<3 ORDER BY created_at DESC;

-- name: MailSent :exec
UPDATE connections 
SET sent_attempts = sent_attempts +1 , status_sent="SENT"
WHERE connection_id=$1;

-- name: MailnotSent :one
UPDATE connections 
SET sent_attempts = sent_attempts +1 
WHERE connection_id=$1 RETURNING sent_attempts,status_sent;