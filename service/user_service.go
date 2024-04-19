package service

import (
	"time"

	"github.com/yobadagne/user_registration/auth"
	"github.com/yobadagne/user_registration/data"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/token"
	"github.com/yobadagne/user_registration/util"
	"github.com/yobadagne/user_registration/val"
	"go.uber.org/zap"
)

var NewAuthLayer = auth.NewAuthLayer()
var NewDataLayer = data.NewDataLayer(NewAuthLayer)
var NewValLayer = val.NewValidateLayer()
var NewTokenLayer = token.NewTokenLayer()

type ServiceLayer struct {
	datalayer     model.DataLayer
	authlayer     model.AuthLayer
	validatelayer model.ValidaterLayer
	tokenlayer    model.TokenLayer
}

func NewServiceLayer() model.ServiceLayer {
	return &ServiceLayer{
		datalayer:     NewDataLayer,
		authlayer:     NewAuthLayer,
		validatelayer: NewValLayer,
		tokenlayer:    NewTokenLayer,
	}
}

// generate access and refresh tokens
func (s *ServiceLayer) GernerateAccessAndRefreshToken(username string, userID int) (access_token, refresh_token string, err error) {
	config, err := util.LoadConfig(".")
	if err != nil {
		return "", "", err
	}
	access_token, err = s.tokenlayer.CreateToken(username, userID, 15*time.Minute, config.Access_key)
	if err != nil {
		return "", "", err
	}
	refresh_token, err = s.tokenlayer.CreateToken(username, userID, 30*24*time.Hour, config.Refersh_key)
	if err != nil {
		return "", "", err
	}
	return access_token, refresh_token, nil
}

//valiadte token

func (s *ServiceLayer) ValidateToken(authorizationHeader string) (*model.Claims, error) {
	config, err := util.LoadConfig(".")
	if err != nil {
		return nil, err
	}

	claims, refresh_token, err := s.tokenlayer.ValidateToken(authorizationHeader, config.Refersh_key)
	// set the global variable UserID
	model.UserID = claims.UserID
	if err != nil {
		return nil, err
	}

	//check for validation of refresh token session
	refresh_tokenfromDB, err := s.datalayer.GetRefreshToken(claims.Username)
	if err != nil {
		return nil, err
	}
	// refresh_tokenfromDB , err = s.authlayer.DecryptToken(c,refresh_tokenfromDB)
	// if err != nil {
	// 	return nil, err
	// }

	//encrypt the user refresh token
	refresh_token, err = s.authlayer.EncryptToken(refresh_token, model.IV)
	if err != nil {
		return nil, err
	}
	if refresh_token != refresh_tokenfromDB {
		util.Logger.Error("refresh token doesnot match, error while excuting service.ValidateToken()", zap.Error(err), zap.Int("UserID", model.UserID))
		err = model.ErrUnauthorized.NewType("").New("refresh token doesnot match, you can not use one refresh token more than once, please generate a new one")
		model.Error_type = model.UNAUTHORIZED
		return nil, err
	}

	return claims, nil
}

func (s *ServiceLayer) Register(usertoregister model.User) error {
	// validate user
	if err := s.validatelayer.ValidateForRegister(usertoregister); err != nil {
		return err
	}
	//hash password
	password := usertoregister.Password
	hashedpass, err := s.authlayer.GenerateHashPassword(password)
	if err != nil {
		return err
	}
	// assign the hashed password
	usertoregister.Password = hashedpass
	// save to DB
	if err := s.datalayer.Register(usertoregister); err != nil {
		return err
	}
	return nil
}

// for login

func (s *ServiceLayer) Login(usertolog model.User) (string, string, error) {
	// validate user
	if err := s.validatelayer.ValidateForLogin(usertolog); err != nil {
		return "", "", err
	}

	//get password from DB
	UserFromDB, err := s.datalayer.GetUserForLogin(usertolog)
	if err != nil {
		return "", "", err
	}
	// compare password

	if err := s.authlayer.CompareHashPassword(usertolog.Password, UserFromDB.Password); err != nil {
		return "", "", err
	}
	// delete its session , if it has
	if err := s.datalayer.DeleteRefreshTokenForLoginIfExists(usertolog.Username); err != nil {
		return "", "", err
	}
	// generate access and referesh token
	access_token, refresh_token, err := s.GernerateAccessAndRefreshToken(UserFromDB.Username, int(UserFromDB.ID))
	if err != nil {
		return "", "", err
	}
	// save the refresh token to session table to make sure we only use it once
	if err := s.datalayer.SaveNewRefreshToken(refresh_token, usertolog.Username); err != nil {
		return "", "", err
	}
	return access_token, refresh_token, nil
}

//for refresh

func (s *ServiceLayer) Refresh(authorizationHeader string) (string, string, error) {
	claims, err := s.ValidateToken(authorizationHeader)
	if err != nil {
		return "", "", err
	}
	//delete token from database
	if err := s.datalayer.DeleteUsedRefreshToken(claims.Username); err != nil {
		return "", "", err
	}
	// now generate new tokens
	access_token, refresh_token, err := s.GernerateAccessAndRefreshToken(claims.Username, claims.UserID)
	if err != nil {
		return "", "", err
	}
	// save the refresh token to session table to make sure we only use it once
	if err := s.datalayer.SaveNewRefreshToken(refresh_token, claims.Username); err != nil {
		return "", "", err
	}
	return access_token, refresh_token, nil
}
func (s *ServiceLayer) GetRegisteredUser(username string) (string, error) {
	username, err := s.datalayer.GetRegisteredUser(username)
	if err != nil {
		return "", err
	}
	return username, nil
}
