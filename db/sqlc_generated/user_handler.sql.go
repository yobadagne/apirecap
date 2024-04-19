// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user_handler.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  email
) VALUES (
  $1, $2, $3
) RETURNING id, username, password, email
`

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
	)
	return i, err
}

const getRegisteredUser = `-- name: GetRegisteredUser :one
SELECT username FROM users 
WHERE username = $1
`

func (q *Queries) GetRegisteredUser(ctx context.Context, username string) (string, error) {
	row := q.db.QueryRowContext(ctx, getRegisteredUser, username)
	err := row.Scan(&username)
	return username, err
}

const getUserForLogin = `-- name: GetUserForLogin :one
SELECT id, username, password, email FROM users 
WHERE 
username = $1
`

func (q *Queries) GetUserForLogin(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserForLogin, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
	)
	return i, err
}
