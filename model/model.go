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
	Register(User, *gin.Context) error
	GetPasswordForLogin(User, *gin.Context) (string, error)
	SaveNewRefreshToken(c *gin.Context, refreshtoken, username string) error
	DeleteUsedRefreshToken(c *gin.Context, refreshtoken string) error
	GetRefreshToken(c *gin.Context, refreshtoken string) (string, error)
}

//port for handler

type HandlerLayer interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

//port for auth layer

type AuthLayer interface {
	GenerateHashPassword(c *gin.Context, password string) (string, error)
	CompareHashPassword(c *gin.Context, password, hash string) error
	EncryptToken(c *gin.Context, token string, iv []byte) (string, error)
	DecryptToken(c *gin.Context, ciphertext string) (string, error)
}

// port for validating user input
type ValidaterLayer interface {
	ValidateForRegister(c *gin.Context, user User) error
	ValidateForLogin(c *gin.Context, user User) error
	ValidateEmail(c *gin.Context, email string) error
	VerifyPassword(c *gin.Context, s string) error
}

// port for token
type TokenLayer interface {
	CreateToken(c *gin.Context, username string, duration time.Duration, key string) (string, error)
	ValidateToken(c *gin.Context, key string) (*Claims, string, error)
}
