package handler

// implement handlers
import (
	
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
	//ctx := c.Request.Context()
	err := h.servicelayer.Register(usertoregister)

	if err != nil {
		c.Error(err)
		return
	}
	// aborting the request will be done in the error handler middleware
	c.JSON(http.StatusOK, "new user registered")
}

func (h HandlerLayer) Login(c *gin.Context) {
	var usertolog model.User
	//bind user info
	if err := c.BindJSON(&usertolog); err != nil {
		err = errorx.Decorate(err, "Error binding user input for login")
		util.Logger.Error("Error binding user input for login", zap.Error(err))
		c.Error(err)
		c.Set(model.Error_type,model.INTERNAL_SERVER_ERROR)
		return
	}
	access_token,refresh_token,err := h.servicelayer.Login(usertolog)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"access_token":  access_token,
		"refresh_token": refresh_token,
	})

}

func (h HandlerLayer) Refresh(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	access_token,refresh_token, err := h.servicelayer.Refresh(authorization); 
	if err != nil {
		c.Error(err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})

}
