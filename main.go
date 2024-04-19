package main

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/yobadagne/user_registration/handler"
	"github.com/yobadagne/user_registration/middleware"
	"github.com/yobadagne/user_registration/util"

)

func main() {
	util.InitializeLogger()
	handler := handler.NewHandlerLayer()
	r := gin.Default()
	r.Use(middleware.RequestID())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.TimeoutMiddleware(5 * time.Second))
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/refresh", handler.Refresh)
	err := r.Run(":8080")
	if err != nil {
		util.Logger.Error("Can not start server")
		return
	}
	util.Logger.Info("Server started at 8080")

}
