package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joomcode/errorx"
	"github.com/yobadagne/user_registration/model"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Next()
	}
}
func UserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := uuid.New().String()
		c.Set("UserID", userID)
		c.Next()
	}
}

// time out middleware
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel() // Make sure to call cancel to release resources when done
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
		// use map to bind error code
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Unwrap() // get the original error
			errx := errorx.Cast(err) // cast to errorx to use the Message mthod so that we can diplay the message to user
			c.AbortWithStatusJSON(model.HttpCodeGenerator[model.Error_type], gin.H{"err": errx.Message()})

		}
	}
}
