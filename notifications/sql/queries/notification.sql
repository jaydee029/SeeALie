-- name: DequeNotifications :many
SELECT request_init_by, request_to FROM notifications WHERE status_sent=FALSE AND sent_attempts<3 ORDER BY created_at DESC;

-- name: NotificationSent :exec
UPDATE notifications
SET sent_attempts = sent_attempts +1 , status_sent=TRUE
WHERE request_init_by=$1 AND request_to=$2;

-- name: NotificationnotSent :one
UPDATE notifications
SET sent_attempts = sent_attempts +1 
WHERE request_init_by=$1 AND request_to=$2
RETURNING sent_attempts,status_sent;