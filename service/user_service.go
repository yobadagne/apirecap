package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yobadagne/user_registration/auth"
	"github.com/yobadagne/user_registration/data"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/token"
	"github.com/yobadagne/user_registration/util"
	"github.com/yobadagne/user_registration/val"
)

var NewAuthLayer = auth.NewAuthLayer()
var NewDataLayer= data.NewDataLayer(NewAuthLayer)
var NewValLayer = val.NewValidateLayer()
var NewTokenLayer = token.NewTokenLayer()

type ServiceLayer struct {
	datalayer     model.DataLayer
	authlayer     model.AuthLayer
	validatelayer model.ValidaterLayer
	tokenlayer    model.TokenLayer
}

func NewServiceLayer() *ServiceLayer {
	return &ServiceLayer{
		datalayer:     NewDataLayer,
		authlayer:     NewAuthLayer,
		validatelayer: NewValLayer,
		tokenlayer:    NewTokenLayer,
	}
}

// generate access and refresh tokens
func (s ServiceLayer) GernerateAccessAndRefreshToken(c *gin.Context, username string) (access_token, refresh_token string, err error) {
	config, err := util.LoadConfig(c,".")
	if err != nil {
		return "", "", err
	}
	access_token, err = s.tokenlayer.CreateToken(c,username, 15*time.Minute, config.Access_key)
	if err != nil {
		return "", "", err
	}
	refresh_token, err = s.tokenlayer.CreateToken(c, username, 30*24*time.Hour, config.Refersh_key)
	if err != nil {
		return "", "", err
	}
	return access_token, refresh_token, nil
}

//valiadte token

func (s ServiceLayer) ValidateToken(c *gin.Context) (*model.Claims, string, error) {
	config, err := util.LoadConfig(c,".")
	if err != nil {
		return nil, " ", err
	}

	claims, refresh_token, err := s.tokenlayer.ValidateToken(c, config.Refersh_key)
	if err != nil {
		return nil, "", err
	}

	//check for validation of refresh token session
	refresh_token,err = s.datalayer.GetRefreshToken(c,refresh_token)
	if err!= nil{
		return nil, "", err
	}
	return claims, refresh_token, nil
}

// port for database communication

func (s ServiceLayer) Register(usertoregister model.User, c *gin.Context) error {
	// validate user
	if err := s.validatelayer.ValidateForRegister(c, usertoregister); err != nil {
		return err
	}
	if err := s.validatelayer.VerifyPassword(c,usertoregister.Password); err != nil {
		return err
	}
	if err := s.validatelayer.ValidateEmail(c,usertoregister.Email); err != nil {
		return err
	}

	//hash password
	password := usertoregister.Password
	hashedpass, err := s.authlayer.GenerateHashPassword(c,password)
	if err != nil {
		return err
	}
	// assign the hashed password
	usertoregister.Password = hashedpass
	// save to DB
	if err := s.datalayer.Register(usertoregister, c); err != nil {
		return err
	}
	return nil
}

// for login

func (s ServiceLayer) Login(usertolog model.User, c *gin.Context) error {
	// validate user
	if err := s.validatelayer.ValidateForLogin(c,usertolog); err != nil {
		return err
	}

	//get password from DB
	passwordfromDB, err := s.datalayer.GetPasswordForLogin(usertolog, c)
	if err != nil {
		return err
	}
	// compare password

	if err := s.authlayer.CompareHashPassword(c, usertolog.Password, passwordfromDB); err != nil {
		return err
	}
	// generate access and referesh token
	access_token, refresh_token, err := s.GernerateAccessAndRefreshToken(c,usertolog.Username)
	if err != nil {
		return err
	}
	// save the refresh token to session table to make sure we only use it once
	if err := s.datalayer.SaveNewRefreshToken(c, refresh_token, usertolog.Username); err != nil {
		return err
	}
	c.Set("access_token", access_token)
	c.Set("refresh_token", refresh_token)
	return nil
}
//for refresh

func (s ServiceLayer) Refresh (c *gin.Context) error{
	claims, refresh_token, err := s.ValidateToken(c)
	if err!= nil{
		return err
	}
	// delete token from database
	// if err := s.datalayer.DeleteUsedRefreshToken(c,refresh_token); err!= nil{
	// 	return err
	// }
	// now generate new tokens
	access_token, refresh_token, err := s.GernerateAccessAndRefreshToken(c,claims.Username)
	if err != nil {
		return err
	}
	// save the refresh token to session table to make sure we only use it once
	if err := s.datalayer.SaveNewRefreshToken(c, refresh_token, claims.Username); err != nil {
		return err
	}
	c.Set("access_token", access_token)
	c.Set("refresh_token", refresh_token)
	return nil
}