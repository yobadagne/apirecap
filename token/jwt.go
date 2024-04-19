package token

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

type TokenLayer struct {
}

func NewTokenLayer() model.TokenLayer {
	return &TokenLayer{}
}
func (t *TokenLayer) CreateToken(username string,userID int, duration time.Duration, key string) (string, error) {
	expiretime := time.Now().Add(duration)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		Username: username,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiretime.Unix(),
		},
	})
	token, err := jwtToken.SignedString([]byte(key))
	if err!= nil{
		util.Logger.Error("Token creation error,error while excuting token.CreateToken()", zap.Error(err))
		err = model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "can not create token",
		}
		return "", err
	}
	return token, nil
}

func (t *TokenLayer) ValidateToken(authorizationHeader,key string) (*model.Claims, string, error) {
	authorization := authorizationHeader
	if authorization == " " {
		err := model.MyError{
			Code:    http.StatusBadRequest,
			Message: "Invalid authorization header",
		}
		
		util.Logger.Error("Invalid authorization header,error while excuting token.ValidateToken()", zap.Error(err))
		
		return nil, " ", err
	}

	fields := strings.Split(authorization, " ")
	if strings.ToLower(fields[0]) != "bearer" || len(fields) < 2 {
		err := model.MyError{
			Code:    http.StatusBadRequest,
			Message: "Invalid token type",
		}
		util.Logger.Error("Invalid token type,error while excuting token.ValidateToken()", zap.Error(err))	
		return nil, " ", err
	}
	tokenstring := fields[1]
	Claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenstring, Claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		util.Logger.Error("Invalid token when parsing, error while excuting token.ValidateToken()", zap.Error(err))
		err := model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid token when parsing",
		}
		return nil, " ", err
	}
	if !token.Valid {
		err := model.MyError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid token when validating",
		}
		
		util.Logger.Error("Invalid token when validating,error while excuting token.ValidateToken()",  zap.Error(err))
		
		return nil, " ", err
	}
	return Claims, fields[1], nil
}
