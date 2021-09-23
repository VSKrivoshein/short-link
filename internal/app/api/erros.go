package api

import (
	"github.com/VSKrivoshein/short-link/internal/app/e"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrResponse struct {
	Error string `json:"error" example:"user definition of error"`
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()

		if err == nil {
			return
		}

		unwrappedErr := e.UnwrapRecursive(err.Err)

		customErr, ok := unwrappedErr.(*e.CustomErrorWithCode)
		if !ok {
			logrus.Fatalf("error: error was not found in unwrappedErr.(*e.CustomErrorWithCode)")
		}

		c.JSON(customErr.Code, ErrResponse{
			Error: customErr.ErrorForUser(),
		})
	}
}
