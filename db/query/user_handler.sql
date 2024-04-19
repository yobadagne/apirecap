-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  email
) VALUES (
  $1, $2, $3
) RETURNING *;


-- name: GetUserForLogin :one
SELECT * FROM users 
WHERE 
username = $1;

-- name: GetRegisteredUser :one
SELECT username FROM users 
WHERE username = $1;




