package token

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
func (t TokenLayer) CreateToken(c *gin.Context, username string, duration time.Duration, key string) (string, error) {
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
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return "", err
	}
	return token, nil
}

func (t TokenLayer) ValidateToken(c *gin.Context, key string) (*model.Claims, string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		util.Logger.Error("Invalid authorization header")
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return nil, " ", errorx.IllegalState.New("Invalid authorization header")
	}

	fields := strings.Split(authorization, " ")
	if strings.ToLower(fields[0]) != "bearer" || len(fields) < 2 {
		util.Logger.Error("Invalid token type")
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return nil, " ", errorx.IllegalState.New("Invalid token type")
	}
	tokenstring := fields[1]
	Claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenstring, Claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		err = errorx.Decorate(err,"Invalid token" )
		util.Logger.Error("Invalid token", zap.Error(err))
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return nil, " ", err
	}
	if !token.Valid {
		util.Logger.Error("Invalid token")
		c.Set(model.Error_type,model.UNAUTHORIZED)
		return nil, " ", errorx.IllegalState.New("Invalid token")
	}
	return Claims, fields[1], nil
}
