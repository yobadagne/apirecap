// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

type Session struct {
	ID           int32  `json:"id"`
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
