package api

import (
	u "github.com/VSKrivoshein/short-link/internal/app/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	startErrStatuses = 500
)

// JSONLogMiddleware logs a gin HTTP requests in JSON format, with some additional custom key/values
func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := u.GetDurationInMilliseconds(start)

		entry := log.WithFields(log.Fields{
			"method":    c.Request.Method,
			"path":      c.Request.RequestURI,
			"status":    c.Writer.Status(),
			"duration":  duration,
			"user_id":   u.GetUserID(c),
			"client_ip": u.GetClientIP(c),
		})

		if c.Writer.Status() >= startErrStatuses {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("")
		}
	}
}
