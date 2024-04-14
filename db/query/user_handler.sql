-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  email
) VALUES (
  $1, $2, $3
) RETURNING *;


-- name: GetPasswordForLogin :one
SELECT password FROM users 
WHERE 
username = $1;


