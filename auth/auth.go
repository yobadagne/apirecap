package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"crypto/sha256"
	
)

// adapter for auth layer
type AuthLayer struct {
}

func NewAuthLayer() model.AuthLayer {
	return &AuthLayer{}
}

func (a AuthLayer) GenerateHashPassword(c *gin.Context,password string) (string, error) {
	if len(password) > 72 {
			// Hash the refresh token with SHA-256
		sha256Hash := sha256.Sum256([]byte(password))

		// Convert the SHA-256 hash to a string
		password = string(sha256Hash[:])
    }
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		// change error into errorx format
		
		err = errorx.Decorate(err, "Can not hash password")
		util.Logger.Error("Can not hash password", zap.Error(err))
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return "", err
	}
	return string(bytes), nil
}

// compare for login
func (a AuthLayer) CompareHashPassword(c *gin.Context,password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// change error into errorx format
		err = errorx.Decorate(err, "Can not compare password")
		util.Logger.Error("Error while comparing passoword", zap.Error(err))
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return err
	}
	return nil
}
