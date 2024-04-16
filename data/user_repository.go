package data

// database or mock implemenation
import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/lib/pq"
	db "github.com/yobadagne/user_registration/db/sqlc_generated"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

var Queries *db.Queries

func OpenDB() (*db.Queries, error) {

	DB, err := sql.Open("postgres", "postgresql://root:yobadagne2nd@localhost:5432/users_db?sslmode=disable")
	if err != nil {
		// change error into errorx format
		err = errorx.Decorate(err, "Can not open Database")
		util.Logger.Error("Can not open Database", zap.Error(err))
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
	q, err := OpenDB()
	if err != nil {
		util.Logger.Error("Can not open DB")
	}
	return &Datalayer{
		q:    q,
		auth: auth,
	}
}

func (d *Datalayer) Register(newuser model.User, c *gin.Context) error {
	args := db.CreateUserParams{
		Username: newuser.Username,
		Password: newuser.Password,
		Email:    newuser.Email,
	}
	_, err := d.q.CreateUser(c, args)
	if err != nil {
		// check for constarint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // constraint violation error in postgress
				err = errorx.Decorate(pqErr, "Username or email exists")
				util.Logger.Error("Username or email exists", zap.Error(err))
				c.Set(model.Error_type, model.BAD_REQUEST)
				return err
			}
			err = errorx.Decorate(pqErr, pqErr.Message)
			util.Logger.Error(pqErr.Message, zap.Error(err))
			c.Set(model.Error_type, model.BAD_REQUEST)
			return err
		}
		// change error into errorx format
		err = errorx.Decorate(err, "Can not Create user")
		util.Logger.Error("Can not Create user", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return err
	}
	return nil
}

func (d *Datalayer) GetPasswordForLogin(usertolog model.User, c *gin.Context) (string, error) {
	passwordfromDB, err := d.q.GetPasswordForLogin(c, usertolog.Username)
	if err != nil {
		if err == sql.ErrNoRows {

			err = errorx.Decorate(err, "No user found")
			util.Logger.Error("No user found", zap.Error(err))
			c.Set(model.Error_type, model.NOT_FOUND)
			return "", err
		}
		err = errorx.Decorate(err, "Error while reading user from DB")
		util.Logger.Error("Error while reading user from DB", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return "", err
	}
	return passwordfromDB, nil
}

func (d *Datalayer) SaveNewRefreshToken(c *gin.Context, refreshtoken, username string) error {

	encrypted_refresh_token, err := d.auth.EncryptToken(c, refreshtoken, model.IV)
	if err != nil {
		return err
	}

	args := db.SaveNewRefreshTokenParams{
		Username:     username,
		RefreshToken: encrypted_refresh_token,
	}
	_, err = d.q.SaveNewRefreshToken(c, args)
	if err != nil {
		// check for constarint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // constraint violation error in postgres
				err = errorx.Decorate(pqErr, "Username exists")
				util.Logger.Error("Username exists", zap.Error(err))
				c.Set(model.Error_type, model.BAD_REQUEST)
				return err
			}
			err = errorx.Decorate(pqErr, pqErr.Message)
			util.Logger.Error(pqErr.Message, zap.Error(err))
			c.Set(model.Error_type, model.BAD_REQUEST)
			return err
		}
		err = errorx.Decorate(err, "Can not save refresh token")
		util.Logger.Error("Can not save refresh token", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return err
	}
	return nil
}
func (d *Datalayer) DeleteUsedRefreshToken(c *gin.Context, username string) error {
	err := d.q.DeleteUsedRefreshToken(c, username)
	if err != nil {
		err = errorx.Decorate(err, "Can not delete refresh token")
		util.Logger.Error("Can not delete refresh token", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return err
	}
	return nil
}
func (d *Datalayer) GetRefreshToken(c *gin.Context, username string) (string, error) {
	refresh_token, err := d.q.GetRefreshToken(c, username)
	if err != nil {
		err = errorx.Decorate(err, "Can not fetch refresh token")
		util.Logger.Error("Can not fetch refresh token", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return "", err
	}
	return refresh_token, nil
}
func (d *Datalayer) DeleteRefreshTokenForLoginIfExists(c *gin.Context, username string) error {
	err := d.q.DeleteRefreshTokenForLoginIfExists(c, username)
	if err != nil {
		err = errorx.Decorate(err, "Error while executing DeleteRefreshTokenForLoginIfExists")
		util.Logger.Error("Error while executing DeleteRefreshTokenForLoginIfExists", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return err
	}
	return nil
}
