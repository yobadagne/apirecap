package token

import (
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/joomcode/errorx"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

type TokenLayer struct {
}

func NewTokenLayer() model.TokenLayer {
	return &TokenLayer{}
}
func (t TokenLayer) CreateToken(username string, duration time.Duration, key string) (string, error) {
	expiretime := time.Now().Add(duration)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiretime.Unix(),
		},
	})
	token, err := jwtToken.SignedString([]byte(key))
	if err!= nil{
		err = errorx.Decorate(err, "Token creation error")
		util.Logger.Error("Token creation error", zap.Error(err))
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return "", err
	}
	return token, nil
}

func (t TokenLayer) ValidateToken(authorizationHeader,key string) (*model.Claims, string, error) {
	authorization := authorizationHeader
	if authorization == " " {
		err := errorx.IllegalState.New("Invalid authorization header")
		util.Logger.Error("Invalid authorization header", zap.Error(err))
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return nil, " ", err
	}

	fields := strings.Split(authorization, " ")
	if strings.ToLower(fields[0]) != "bearer" || len(fields) < 2 {
		err := errorx.IllegalState.New("Invalid token type")
		util.Logger.Error("Invalid token type", zap.Error(err))
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return nil, " ", err
	}
	tokenstring := fields[1]
	Claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenstring, Claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		err = errorx.Decorate(err,"Invalid token when parsing" )
		util.Logger.Error("Invalid token when parsing", zap.Error(err))
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return nil, " ", err
	}
	if !token.Valid {
		err := errorx.IllegalState.New("Unauthorized")
		util.Logger.Error("Invalid token when  validating")
		model.Error_type = model.UNAUTHORIZED
		return nil, " ", err
	}
	return Claims, fields[1], nil
}
