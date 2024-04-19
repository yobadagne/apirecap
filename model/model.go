package model

import (
	"crypto/aes"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/yobadagne/user_registration/db/sqlc_generated"
)
var UserID int
var RequestID uuid.UUID
var Queries *db.Queries
var IV = make([]byte, aes.BlockSize)
type MyError struct{
	Code int
	Message string	
}
// MyError implement the error interface
func(m MyError) Error() string{
	return fmt.Sprintf("Error: %v,%v", m.Code,m.Message)
}
// TODO here try to handle error
var Encriptionkey = []byte("AES128Key-16Char")
type User struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Claims struct {
	Username string
	UserID   int
	jwt.StandardClaims
}
type DataLayer interface {
	Register(User) error
	GetUserForLogin(usertolog User) (db.User, error)
	SaveNewRefreshToken(refreshtoken, username string) error
	DeleteUsedRefreshToken(refreshtoken string) error
	GetRefreshToken(refreshtoken string) (string, error)
	DeleteRefreshTokenForLoginIfExists(username string) error
	GetRegisteredUser(username string) (string, error)
}

//port for handler

type HandlerLayer interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

//port for auth layer

type AuthLayer interface {
	GenerateHashPassword(password string) (string, error)
	CompareHashPassword(password, hash string) error
	EncryptToken(token string, iv []byte) (string, error)
	DecryptToken(ciphertext string) (string, error)
}

// port for validating user input
type ValidaterLayer interface {
	ValidateForRegister(user User) error
	ValidateForLogin(user User) error
}

// port for token
type TokenLayer interface {
	CreateToken(username string,userID int, duration time.Duration, key string) (string, error)
	ValidateToken(authorizationHeader, key string) (*Claims, string, error)
}

type ServiceLayer interface{
	GernerateAccessAndRefreshToken(username string, userID int) (access_token, refresh_token string, err error)
	ValidateToken(authorizationHeader string) (*Claims, error)
	Register(usertoregister User) error
	Login(usertolog User) (string, string, error)
	Refresh(authorizationHeader string) (string, string, error)
	GetRegisteredUser(username string) (string, error)
}