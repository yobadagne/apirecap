package val

import (
	"net/mail"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/joomcode/errorx"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

type ValidateLayer struct {
}

func NewValidateLayer() model.ValidaterLayer {
	return &ValidateLayer{}
}
func (v ValidateLayer) ValidateForRegister(c *gin.Context, u model.User) error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(5, 50), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$"))),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&u.Email, validation.Required, validation.Length(5, 100)))

	if err != nil {
		err  = errorx.Decorate(err,"Invalid User Inputs") 
		util.Logger.Error("Invalid User Inputs", zap.Error(err))
		c.Set(model.Error_type, model.BAD_REQUEST)
		return err
	}
	return nil
}

// validate for login
func (v ValidateLayer) ValidateForLogin(c *gin.Context,u model.User) error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(5, 50), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$"))),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)))

	if err != nil {
		err := errorx.Decorate(err, "Invalid Inputs") 
		util.Logger.Error("Invalid Inputs", zap.Error(err))
		c.Set(model.Error_type, model.BAD_REQUEST)
		return err
	}
	return nil
}
func (v ValidateLayer) ValidateEmail(c *gin.Context,email string) error {
	_, err := mail.ParseAddress(email)
	if err!= nil{
		err = errorx.Decorate(err, "Invalid Email")
		util.Logger.Error("Invalid Email", zap.Error(err))
		c.Set(model.Error_type, model.BAD_REQUEST)
		return err
	}
	 return nil
}

// password verification for its character
func (v ValidateLayer) VerifyPassword(c *gin.Context,s string) error{
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, ch := range s {
		switch {
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsUpper(ch):
			hasUpperCase = true
		case unicode.IsLower(ch):
			hasLowercase = true
		case ch == '#' || ch == '|':
			return errorx.Decorate(errorx.IllegalFormat.New("Invalid password format"), "Password not supported")
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	if !(hasNumber && hasUpperCase && hasLowercase && hasSpecial){
		util.Logger.Error("Invalid password format")
		c.Set(model.Error_type, model.BAD_REQUEST)
		return errorx.Decorate(errorx.IllegalFormat.New("Invalid password format"), "Password not supported")
	}
	return nil
}
