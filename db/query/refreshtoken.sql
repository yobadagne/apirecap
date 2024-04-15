-- name: SaveNewRefreshToken :one
INSERT INTO sessions (
    username, 
    refresh_token 
    ) VALUES (
        $1,$2
    ) RETURNING *;
-- name: DeleteUsedRefreshToken :exec
DELETE FROM sessions
WHERE username = $1;
-- name: GetRefreshToken :one
SELECT refresh_token FROM sessions
WHERE username = $1
LIMIT 1;

