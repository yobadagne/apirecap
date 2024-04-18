package model

import (
	"crypto/aes"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var IV = make([]byte, aes.BlockSize)
//  TODO here try to handle error 
var Encriptionkey = []byte("AES128Key-16Char")
var Error_type = "error_type"
var (
	BAD_REQUEST           = "bad request"
	INTERNAL_SERVER_ERROR = "internal server error"
	UNAUTHORIZED          = "unauthorized"
	NOT_FOUND             = "not Found"
)
var HttpCodeGenerator = map[string]int{
	BAD_REQUEST:           http.StatusBadRequest,
	INTERNAL_SERVER_ERROR: http.StatusInternalServerError,
	UNAUTHORIZED:          http.StatusUnauthorized,
	NOT_FOUND:             http.StatusNotFound,
}

type Mystring string
type User struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Claims struct {
	Username string
	jwt.StandardClaims
}
type DataLayer interface {
	Register(User) error
	GetPasswordForLogin(User) (string, error)
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
	CreateToken(username string, duration time.Duration, key string) (string, error)
	ValidateToken(authorizationHeader,key string) (*Claims, string, error)
}
