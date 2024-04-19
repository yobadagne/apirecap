package val

import (
	"net/http"
	"regexp"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

type ValidateLayer struct {
}

func NewValidateLayer() model.ValidaterLayer {
	return &ValidateLayer{}
}
func (v *ValidateLayer) ValidateForRegister(u model.User) error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(5, 50), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$"))),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50),validation.By(VerifyPassword)),
		validation.Field(&u.Email, validation.Required, is.Email, validation.Length(5, 100)))

	if err != nil {
		util.Logger.Error("Invalid User Inputs, error while excuting val.ValidateForRegister()", zap.Error(err))
		err = model.MyError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return err
	}
	return nil
}

// validate for login
func (v *ValidateLayer) ValidateForLogin(u model.User) error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(5, 50), validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$"))),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 50)))

	if err != nil {
		util.Logger.Error("Invalid Inputs, error while excuting val.ValidateForLogin()", zap.Error(err))
		err = model.MyError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		return err
	}
	return nil
}

// password verification for its character
func VerifyPassword(value interface {}) error{
	s,ok := value.(string)
	if !ok{
		model.ErrBadRequest.New("Password must be string")
	}
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
			return model.ErrBadRequest.New("Password not supported shouldn't include # or |")
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	if !(hasNumber && hasUpperCase && hasLowercase && hasSpecial){
		util.Logger.Error("Invalid password format, error while excuting val.VerifyPassword")
		err := model.MyError{
			Code:    http.StatusBadRequest,
			Message: "Password not supported must contain capital letters and special characters",
		}
		return err	
	}
	return nil
}
