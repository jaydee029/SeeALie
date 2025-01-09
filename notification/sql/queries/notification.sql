-- name: DequeNotifications :many
SELECT request_by, request_to, request_status FROM notifications WHERE status_sent=PENDING AND sent_attempts<3 ORDER BY created_at DESC;

-- name: NotificationSent :exec
UPDATE connections 
SET sent_attempts = sent_attempts +1 , status_sent="SENT"
WHERE connection_id=$1;

-- name: NotificationnotSent :one
UPDATE connections 
SET sent_attempts = sent_attempts +1 
WHERE connection_id=$1 RETURNING sent_attempts,status_sent;