-- name: Find_friends :many
SELECT
 CASE 
  WHEN followed_by=$1 THEN followed
  WHEN followed=$1 THEN followed_by
 END::varchar(12) AS friend, room_id
FROM friends;

-- name: Find_room :many
SELECT room_id FROM friends 
WHERE (followed_by=$1 AND followed=$2) 
OR 
(followed_by=$2 AND followed=$1);
