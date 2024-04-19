package data

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/lib/pq"
	db "github.com/yobadagne/user_registration/db/sqlc_generated"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

var ctx = context.Background()

// adapters to implement the datalayer port
type Datalayer struct {
	q    *db.Queries
	auth model.AuthLayer
}

func NewDataLayer(auth model.AuthLayer) model.DataLayer {
	return &Datalayer{
		q:    model.Queries,
		auth: auth,
	}
}

func (d *Datalayer) Register(newuser model.User) error {
	args := db.CreateUserParams{
		Username: newuser.Username,
		Password: newuser.Password,
		Email:    newuser.Email,
	}
	_, err := d.q.CreateUser(ctx, args)
	if err != nil {
		// check for constarint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // constraint violation error in postgress
				firstpart := strings.Split(pqErr.Detail, " ")[1]
				column := strings.Split(firstpart, "=")[0]
				if column == "(username)" {
					err = model.MyError{
						Code:    http.StatusBadRequest,
						Message: "Username exists, must be unique , if you are registerd please login",
					}
					util.Logger.Error("Username exists in the database, Should be Unique,error while excuting data.Register()", zap.Error(err))	
					
					return err
				} else {
					err = model.MyError{
						Code:    http.StatusBadRequest,
						Message: "Email exists, must be unique , if you are registerd please login",
					}
					util.Logger.Error("Email exists in the database, Should be Unique,error while excuting data.Register()", zap.Error(err))
					
					
					return err
				}
			}
			util.Logger.Error(pqErr.Message, zap.Error(err))
			err = model.MyError{
				Code:    http.StatusInternalServerError,
				Message: pqErr.Message,
			}
			return err
		}
		// change error into errorx format
		util.Logger.Error("Can not Create user, error while excuting data.Register()", zap.Error(err))
		err = model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "Can not Create user",
		}
		return err
	}
	return nil
}

func (d *Datalayer) GetUserForLogin(usertolog model.User) (db.User, error) {
	UserFromDB, err := d.q.GetUserForLogin(ctx, usertolog.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Logger.Error("No Registeed user found with the given username, error while excuting data.GetUserForLogin()", zap.Error(err))
			err = model.MyError{
				Code:    http.StatusNotFound,
				Message: "No user found please register",
			}
			return db.User{}, err
		}
		util.Logger.Error("Error while reading user from DB , error while excuting data.GetUserForLogin()", zap.Error(err))
		err = model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "Can not read from DB",
		}
		return db.User{}, err
	}
	return UserFromDB, nil
}

func (d *Datalayer) SaveNewRefreshToken(refreshtoken, username string) error {

	encrypted_refresh_token, err := d.auth.EncryptToken(refreshtoken, model.IV)
	if err != nil {
		return err
	}

	args := db.SaveNewRefreshTokenParams{
		Username:     username,
		RefreshToken: encrypted_refresh_token,
	}
	_, err = d.q.SaveNewRefreshToken(ctx, args)
	if err != nil {
		// check for constarint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // constraint violation error in postgres
				util.Logger.Error("There is a refresh token that is already saved with this username,error while excuting data.SaveNewRefreshToken() ", zap.Error(err), zap.Int("UserID", model.UserID))
				err = model.MyError{
					Code:    http.StatusBadRequest,
					Message: "Please use your refresh token to get a new refresh token",
				}
				return err
			}
			util.Logger.Error(pqErr.Message, zap.String("Message:", "error while excuting data.SaveNewRefreshToken()"), zap.Error(err), zap.Int("UserID", model.UserID))
			err = model.MyError{
				Code:    http.StatusInternalServerError,
				Message: pqErr.Message,
			}
			return err
		}
		util.Logger.Error("Can not save refresh token, Error while excuting  SaveNewRefreshToken()", zap.Error(err), zap.Int("UserID", model.UserID))
		err = model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "Error while processing your refresh token",
		}
		return err
	}
	return nil
}
func (d *Datalayer) DeleteUsedRefreshToken(username string) error {
	err := d.q.DeleteUsedRefreshToken(ctx, username)
	if err != nil {
		util.Logger.Error("Can not delete the used refresh token, error while excuting data.DeleteUsedRefreshToken()", zap.Error(err), zap.Int("UserID", model.UserID))
		
		err = model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "Error while processing your refresh token",
		}
		return err
	}
	return nil
}
func (d *Datalayer) GetRefreshToken(username string) (string, error) {
	refresh_token, err := d.q.GetRefreshToken(ctx, username)
	if err != nil {
		util.Logger.Error("Can not fetch refresh token for comparison, error while excuting data.GetRefreshToken()", zap.Error(err))
		err = model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "Error while processing your refresh token can not use one refresh token more than once, please generate a new one",
		}
		return "", err
	}
	return refresh_token, nil
}
func (d *Datalayer) DeleteRefreshTokenForLoginIfExists(username string) error {
	err := d.q.DeleteRefreshTokenForLoginIfExists(ctx, username)
	if err != nil {
		util.Logger.Error("Error while executing DeleteRefreshTokenForLoginIfExists", zap.Error(err))
		
		err = model.MyError{
			Code:    http.StatusInternalServerError,
			Message: "Error while processing your refresh token",
		}
		return err
	}
	return nil
}

func (d *Datalayer) GetRegisteredUser(username string) (string, error) {
	username, err := d.q.GetRegisteredUser(ctx, username)
	if err != nil {
		util.Logger.Error("User not found, error while excuting data.GetRegisteredUser()", zap.Error(err))
		err = model.MyError{
			Code:    http.StatusBadRequest,
			Message: "User not found",
		}
		return "", err
	}
	return username, nil
}
