-- name: SaveNewRefreshToken :one
INSERT INTO sessions (
    username, 
    refresh_token 
    ) VALUES (
        $1,$2
    ) RETURNING *;
-- name: DeleteUsedRefreshToken :exec
DELETE FROM sessions
WHERE refresh_token = $1;
-- name: GetRefreshToken :one
SELECT refresh_token FROM sessions
WHERE refresh_token = $1
LIMIT 1;

