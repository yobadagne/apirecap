package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/yobadagne/user_registration/db/sqlc_generated"
	"github.com/yobadagne/user_registration/handler"
	"github.com/yobadagne/user_registration/middleware"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

func main() {
	util.InitializeLogger()
	err := OpenDB()
	if err != nil {
		return
	}
	fmt.Println(err)
	handler := handler.NewHandlerLayer()
	r := gin.Default()
	r.Use(middleware.RequestID())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.TimeoutMiddleware(5 * time.Second))
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/refresh", handler.Refresh)
	err = r.Run(":8080")
	if err != nil {
		util.Logger.Error("Can not start server")
		return
	}
	util.Logger.Info("Server started at 8080")

}
func OpenDB() error {

	DB, err := sql.Open("postgres", "postgresql://root:yobadagne2nd@localhost:5432/users_db?sslmode=disable")
	if err != nil {
		// change error into errorx format
		util.Logger.Error("Can not open Database", zap.Error(err))
		err = model.ErrInternalServerErr.New("Can not open Database")
		return err
	}
	model.Queries = db.New(DB)
	return nil
}
