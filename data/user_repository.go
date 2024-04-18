package data

// database or mock implemenation
import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	db "github.com/yobadagne/user_registration/db/sqlc_generated"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

var Queries *db.Queries
var ctx = context.Background()

func OpenDB() (*db.Queries, error) {

	DB, err := sql.Open("postgres", "postgresql://root:yobadagne2nd@localhost:5432/users_db?sslmode=disable")
	if err != nil {
		// change error into errorx format
		util.Logger.Error("Can not open Database", zap.Error(err))
		err = model.ErrInternalServerErr.New("Can not open Database")
		return nil, err
	}
	Queries = db.New(DB)

	return Queries, nil
}

// adapters to implement the datalayer port
type Datalayer struct {
	q    *db.Queries
	auth model.AuthLayer
}

func NewDataLayer(auth model.AuthLayer) *Datalayer {
	q, _ := OpenDB()
	return &Datalayer{
		q:    q,
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
				util.Logger.Error("Username and email exists", zap.Error(err))
				err = model.ErrBadRequest.New("Username and email exists must be unique")
				model.Error_type =model.BAD_REQUEST
				return err
			}
			util.Logger.Error(pqErr.Message, zap.Error(err))
			err = model.ErrInternalServerErr.New(pqErr.Message)
			model.Error_type = model.INTERNAL_SERVER_ERROR
			return err
		}
		// change error into errorx format
		util.Logger.Error("Can not Create user", zap.Error(err))
		err = model.ErrInternalServerErr.New("Can not Create user")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return err
	}
	return nil
}

func (d *Datalayer) GetPasswordForLogin(usertolog model.User) (string, error) {
	passwordfromDB, err := d.q.GetPasswordForLogin(ctx, usertolog.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			util.Logger.Error("No user found", zap.Error(err))
			err = model.ErrBadRequest.New("No user found please register")
			model.Error_type = model.NOT_FOUND
			return "", err
		}
		util.Logger.Error("Error while reading user from DB", zap.Error(err))
		err = model.ErrBadRequest.New("Can not read from DB")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return "", err
	}
	return passwordfromDB, nil
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
				util.Logger.Error("There is a refresh token that is already saved with this username ", zap.Error(err))
				err = model.ErrBadRequest.New("Please use your refresh token to get a new refresh token")
				model.Error_type = model.BAD_REQUEST
				return err
			}
			util.Logger.Error(pqErr.Message, zap.Error(err))
			err = model.ErrInternalServerErr.New(pqErr.Message)
			model.Error_type = model.INTERNAL_SERVER_ERROR
			return err
		}
		util.Logger.Error("Can not the save refresh token", zap.Error(err))
		err = model.ErrInternalServerErr.New("Error while processing your refresh token")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return err
	}
	return nil
}
func (d *Datalayer) DeleteUsedRefreshToken(username string) error {
	err := d.q.DeleteUsedRefreshToken(ctx, username)
	if err != nil {
		util.Logger.Error("Can not delete the used refresh token", zap.Error(err))
		err = model.ErrInternalServerErr.New("Error while processing your refresh token")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return err
	}
	return nil
}
func (d *Datalayer) GetRefreshToken(username string) (string, error) {
	refresh_token, err := d.q.GetRefreshToken(ctx, username)
	if err != nil {
		util.Logger.Error("Can not fetch refresh token for comparison", zap.Error(err))
		err = model.ErrInternalServerErr.New("Error while processing your refresh token")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return "", err
	}
	return refresh_token, nil
}
func (d *Datalayer) DeleteRefreshTokenForLoginIfExists(username string) error {
	err := d.q.DeleteRefreshTokenForLoginIfExists(ctx, username)
	if err != nil {
		util.Logger.Error("Error while executing DeleteRefreshTokenForLoginIfExists", zap.Error(err))
		err = model.ErrInternalServerErr.New("Error while processing your refresh token")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return err
	}
	return nil
}

func (d *Datalayer) GetRegisteredUser(username string) (string, error) {
	username, err := d.q.GetRegisteredUser(ctx, username)
	if err != nil {
		util.Logger.Error("User not found", zap.Error(err))
		err = model.ErrBadRequest.New("User not found")
		model.Error_type = model.BAD_REQUEST
		return "", err
	}
	return username, nil
}
