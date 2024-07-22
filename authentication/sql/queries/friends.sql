-- name: Find_friends :many
SELECT
 CASE 
  WHEN f.followed_by=$1 THEN f.followed
  WHEN f.followed=$1 THEN f.followed_by
 END AS friend,
u.username
FROM friends f INNER JOIN users u 
ON u.user_id= CASE
   WHEN f.followed_by=$1 THEN f.followed
   WHEN f.followed=$1 THEN f.followed_by
END;


