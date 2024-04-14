package handler

// implement handlers
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"

	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/service"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
)

// Adapter for the handler layer
var NewServiceLayer = service.NewServiceLayer()

type HandlerLayer struct {
	servicelayer service.ServiceLayer
}

func NewHandlerLayer() model.HandlerLayer {
	return &HandlerLayer{
		servicelayer: *NewServiceLayer,
	}
}

//create the datalayer adapter to use it here

func (h HandlerLayer) Register(c *gin.Context) {

	var usertoregister model.User
	if err := c.BindJSON(&usertoregister); err != nil {
		err = errorx.Decorate(err, "Error binding user input for registration")
		util.Logger.Error("Error binding user input for registration", zap.Error(err))
		c.Error(err)
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return
	}
	err := h.servicelayer.Register(usertoregister, c)

	if err != nil {
		c.Error(err)
		fmt.Println("here")
		return
	}
	// aborting the request will be done in the erroe handler middleware
	c.JSON(http.StatusOK, "new user registered")
}

// func (h HandlerLayer) Login(c *gin.Context) {
// 	var usertolog model.User
// 	//bind user info
// 	if err := c.BindJSON(&usertolog); err != nil {
// 		err = errorx.Decorate(err, "Error binding user input for login")
// 		util.Logger.Error("Error binding user input for login", zap.Error(err))
// 		c.Error(err)
// 		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
// 		return
// 	}
// 	err := h.servicelayer.Login(usertolog, c)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	// get the generated access and refresh token from context
// 	value, exists := c.Get("access_token")
// 	if !exists {
// 		err := errorx.Decorate(errorx.InternalError.New("could not get access token from context "), "could not get access token from context ")
// 		c.Error(err)
// 		util.Logger.Error("could not get access token from context", zap.Error(err))
// 		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
// 		return
// 	}
// 	access_token := value

// 	value, exists = c.Get("refresh_token")
// 	if !exists {
// 		err := errorx.Decorate(errorx.InternalError.New("could not get refresh token from context"), "could not get refresh token from context ")
// 		c.Error(err)
// 		util.Logger.Error("could not get refresh token from context", zap.Error(err))
// 		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
// 		return
// 	}
// 	refresh_token := value

// 	c.JSON(http.StatusOK, gin.H{

// 		"access_token":  access_token,
// 		"refresh_token": refresh_token,
// 	})

// }

// func (h HandlerLayer) Refresh(c *gin.Context) {

// 	if err := h.servicelayer.Refresh(c); err != nil {
// 		c.Error(err)
// 		return
// 	}
// 	// get the generated acces and refresh token from context
// 	value, exists := c.Get("access_token")
// 	if !exists {
// 		err := errorx.Decorate(errorx.InternalError.New("could not get access token from context "), "could not get access token from context ")
// 		c.Error(err)
// 		util.Logger.Error("could not get access token from context", zap.Error(err))
// 		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
// 		return
// 	}
// 	access_token := value

// 	value, exists = c.Get("refresh_token")
// 	if !exists {
// 		err := errorx.Decorate(errorx.InternalError.New("could not get refresh token from context"), "could not get refresh token from context ")
// 		c.Error(err)
// 		util.Logger.Error("could not get refresh token from context", zap.Error(err))
// 		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
// 		return
// 	}
// 	refresh_token := value
// 	c.JSON(http.StatusOK, gin.H{
// 		"access_token":  access_token,
// 		"refresh_token": refresh_token,
// 	})

// }
